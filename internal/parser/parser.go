package parser

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// getPrefixFieldMap parses a reflect.Type and returns a map of fields tagged with `mpd_prefix`.
// The map keys are tag values, and values are corresponding StructField definitions.
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

// getNewElementPrefixesSlice parses a reflect.Type and returns a slice of strings.
// It collects the values of the `mpd_prefix` struct tag for fields that also have the
// flag `is_new_element_prefix=true`.
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

// parseLineAndSetFieldValue parses a line, extracts a value from it, and sets the corresponding field
// on the targetElement using the provided map of field definitions.
//
// The 'fields' map contains struct fields indexed by expected prefixes.
// The 'targetElement' must be a reflect.Value pointing to a struct.
// Returns an error if parsing or assignment fails.
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
		case reflect.Ptr:
			switch fieldVal.Type() {
			case reflect.TypeOf((*string)(nil)):
				fieldVal.Set(reflect.ValueOf(&value))
			case reflect.TypeOf((*bool)(nil)):
				v, err := strconv.ParseBool(value)
				if err != nil {
					err = NewFieldParsingError(field.Name, value, fieldVal, err)
					return err
				}
				fieldVal.Set(reflect.ValueOf(&v))
			case reflect.TypeOf((*int)(nil)):
				v, err := strconv.Atoi(value)
				if err != nil {
					err = NewFieldParsingError(field.Name, value, fieldVal, err)
					return err
				}
				fieldVal.Set(reflect.ValueOf(&v))
			case reflect.TypeOf((*uint16)(nil)):
				v, err := strconv.Atoi(value)
				if err != nil {
					err = NewFieldParsingError(field.Name, value, fieldVal, err)
					return err
				}
				uintVal := uint16(v)
				fieldVal.Set(reflect.ValueOf(&uintVal))
			case reflect.TypeOf((*time.Time)(nil)):
				parsedTime, err := time.Parse(time.RFC3339, value)
				if err != nil {
					err = NewFieldParsingError(field.Name, value, fieldVal, err)
					return err
				}
				fieldVal.Set(reflect.ValueOf(&parsedTime))
			default:
				return ErrUnsupportedFieldType
			}
		case reflect.String:
			fieldVal.SetString(value)
		case reflect.Int:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				err = NewFieldParsingError(field.Name, value, fieldVal, err)
				return err
			}
			fieldVal.SetInt(int64(intVal))
		case reflect.Uint16:
			uint16Val, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				err = NewFieldParsingError(field.Name, value, fieldVal, err)
				return err
			}
			fieldVal.SetUint(uint16Val)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				err = NewFieldParsingError(field.Name, value, fieldVal, err)
				return err
			}
			fieldVal.SetBool(boolVal)
		case reflect.Struct:
			if fieldVal.Type() == reflect.TypeOf(time.Time{}) {
				parsedTime, err := time.Parse(time.RFC3339, value)
				if err != nil {
					err = NewFieldParsingError(field.Name, value, fieldVal, err)
					return err
				}
				fieldVal.Set(reflect.ValueOf(parsedTime))
			}
		default:
			return ErrUnsupportedFieldType
		}
	}
	return nil
}

// ParseSingleValue parses the provided MPD response lines into a single value of type T.
//
// Fields in the struct T must be tagged with `mpd_prefix` to allow correct mapping.
// Returns an error if parsing fails or the response is invalid.
func ParseSingleValue[T any](mpdAnswer []string) (T, error) {
	val := reflect.ValueOf(new(T))
	if val.Elem().Kind() != reflect.Struct {
		return *new(T), ErrTargetTypeMustBeStruct
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

// ParseMultiValue parses the provided MPD response lines into a slice of values of type T.
//
// Each group of lines is mapped to a new instance of T.
// Fields in the struct T must be tagged with `mpd_prefix` to allow correct mapping.
// At least one of these tags must also have the flag `is_new_element_prefix=true`
// to indicate the beginning of a new element in the response.
// Returns an error if parsing fails or the input format is invalid.
func ParseMultiValue[T any](mpdAnswer []string) ([]T, error) {
	val := reflect.ValueOf(new(T))
	if val.Elem().Kind() != reflect.Struct {
		return nil, ErrTargetTypeMustBeStruct
	}
	typ := val.Elem().Type()
	var results []T
	fields := getPrefixFieldMap(typ)
	var newElementPrefix = getNewElementPrefixesSlice(typ)
	if len(newElementPrefix) == 0 {
		return nil, ErrNoFieldMarkedAsNewElement
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
