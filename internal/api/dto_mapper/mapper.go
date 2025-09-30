package dto_mapper

import (
	"fmt"
	"github.com/anpotashev/mpdgo/internal/api/dto"
	"github.com/anpotashev/mpdgo/pkg/mpdapi"
	"reflect"
)

func MapDirectory(in mpdapi.DirectoryItem) dto.DirectoryItem {
	return dto.DirectoryItem{}
}
func Map(from, to interface{}) error {
	valFrom := reflect.ValueOf(from)
	valTo := reflect.ValueOf(to)
	if valFrom.Kind() == reflect.Array || valFrom.Kind() == reflect.Slice {
		if valTo.Kind() != reflect.Slice || valTo.Kind() != reflect.Array {
			return fmt.Errorf("expected slice or array, got %s", valFrom.Kind())
		}
	}
	if valFrom.Kind() == reflect.Struct {
		ptr := reflect.New(valFrom.Type())
		ptr.Elem().Set(valFrom)
		valFrom = ptr
	}
	if valFrom.Kind() != reflect.Ptr || valFrom.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("first argument must be a struct or a pointer to a struct1")
	}
	if valTo.Kind() != reflect.Ptr || valTo.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("second argument must be a pointer to a struct")
	}
	for i := 0; i < valFrom.Elem().NumField(); i++ {
		fieldFrom := valFrom.Elem().Field(i)
		fieldFromName := valFrom.Elem().Type().Field(i).Name
		fieldTo := valTo.Elem().FieldByName(fieldFromName)
		if fieldTo.CanSet() {
			if fieldFrom.Type().ConvertibleTo(fieldTo.Type()) {
				converted := fieldFrom.Convert(fieldTo.Type())
				fieldTo.Set(converted)
			}
		}
	}
	return nil
}
