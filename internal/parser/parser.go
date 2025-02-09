package parser

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func getPrefixFieldMap(typ reflect.Type) map[string]reflect.StructField {
	result := make(map[string]reflect.StructField)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		prefix := field.Tag.Get("mpd_prefix")
		if prefix != "" {
			result[prefix] = field
		}
	}
	return result
}

func getNewElementPrefixesSlice(typ reflect.Type) []string {
	var result []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		prefix := field.Tag.Get("mpd_prefix")
		if str := field.Tag.Get("is_new_element_prefix"); str == "true" {
			result = append(result, prefix+":")
		}
	}
	return result
}

func parseLineAndSetFieldValue(fields map[string]reflect.StructField, targetElement reflect.Value, line string) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return errors.New("mpd_prefix does not contain an element prefix")
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if field, ok := fields[key]; ok {
		fieldVal := targetElement.FieldByName(field.Name)
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(value)
		case reflect.Int:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			fieldVal.SetInt(int64(intVal))
		case reflect.Uint16:
			uint16Val, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldVal.SetUint(uint16Val)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldVal.SetBool(boolVal)
		case reflect.Struct:
			if fieldVal.Type() == reflect.TypeOf(time.Time{}) {
				parsedTime, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return err
				}
				fieldVal.Set(reflect.ValueOf(parsedTime))
			}
		default:
			return errors.New("unknown field type")
		}
	}
	return nil
}

func ParseSingleValue[T any](mpdAnswer []string) (T, error) {
	val := reflect.ValueOf(new(T))
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return *new(T), errors.New("T must be a struct")
	}
	typ := val.Elem().Type()
	fields := getPrefixFieldMap(typ)
	result := reflect.New(typ).Elem()
	for _, line := range mpdAnswer {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if err := parseLineAndSetFieldValue(fields, result, line); err != nil {
			return *new(T), err
		}
	}
	return result.Interface().(T), nil
}

func ParseMultiValue[T any](mpdAnswer []string) ([]T, error) {
	val := reflect.ValueOf(new(T))
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, errors.New("T must be a struct")
	}
	typ := val.Elem().Type()
	var results []T
	fields := getPrefixFieldMap(typ)
	var newElementPrefix = getNewElementPrefixesSlice(typ)
	if len(newElementPrefix) == 0 {
		return nil, errors.New("T must have at least one field marked as \"new element prefix\"")
	}
	var current reflect.Value
	for _, line := range mpdAnswer {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		for _, prefix := range newElementPrefix {
			if strings.HasPrefix(line, prefix) {
				if current.IsValid() {
					results = append(results, current.Interface().(T))
				}
				current = reflect.New(typ).Elem()
			}
		}
		if err := parseLineAndSetFieldValue(fields, current, line); err != nil {
			return nil, err
		}
	}
	if current.IsValid() {
		results = append(results, current.Interface().(T))
	}
	return results, nil
}
