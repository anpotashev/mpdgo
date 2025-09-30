package parser

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/assert"
)

type ParsedType struct {
	IntField               int        `mpd_prefix:"int_field" is_new_element_prefix:"true"`
	Uint16Field            uint16     `mpd_prefix:"uint16_field"`
	StringField            string     `mpd_prefix:"string_field"`
	BoolValue              bool       `mpd_prefix:"bool_field"`
	DateField              time.Time  `mpd_prefix:"date_field"`
	IntPtrField            *int       `mpd_prefix:"int_ptr_field"`
	Uint16PtrField         *uint16    `mpd_prefix:"uint16_ptr_field"`
	StringPtrField         *string    `mpd_prefix:"string_ptr_field"`
	BoolPtrValue           *bool      `mpd_prefix:"bool_ptr_field"`
	DatePtrField           *time.Time `mpd_prefix:"date_ptr_field"`
	IntPtrFieldNilValue    *int       `mpd_prefix:"int_ptr_nil_field"`
	Uint16PtrFieldNilValue *uint16    `mpd_prefix:"uint16_ptr_nil_field"`
	StringPtrFieldNilValue *string    `mpd_prefix:"string_ptr_nil_field"`
	BoolPtrValueNilValue   *bool      `mpd_prefix:"bool_ptr_nil_field"`
	DatePtrFieldNilValue   *time.Time `mpd_prefix:"date_ptr_nil_field"`
}

func TestParseSingleValue(t *testing.T) {
	t.Run("successful parsing to struct with all supported field types", func(t *testing.T) {
		var expectedValue ParsedType
		faker.FakeData(&expectedValue)
		list := toMpdResponse(&expectedValue)
		result, err := ParseSingleValue[ParsedType](list)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
	})
	t.Run("parsing field from input line with two or more colons", func(t *testing.T) {
		type parsedType struct {
			Field string `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf:fdsa"}
		actual, err := ParseSingleValue[parsedType](lines)
		assert.NoError(t, err)
		assert.Equal(t, "asdf:fdsa", actual.Field)
	})
	t.Run("parsing field from input with empty lines", func(t *testing.T) {
		type parsedType struct {
			Field string `mpd_prefix:"field"`
		}
		lines := []string{"", "field: asdf", ""}
		actual, err := ParseSingleValue[parsedType](lines)
		assert.NoError(t, err)
		assert.Equal(t, "asdf", actual.Field)
	})
	t.Run("error parsing value for int field", func(t *testing.T) {
		type parsedType struct {
			Field int `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrParsingField))
	})
	t.Run("error parsing value for *int field", func(t *testing.T) {
		type parsedType struct {
			Field *int `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing value for uint16 field", func(t *testing.T) {
		type parsedType struct {
			Field uint16 `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing value for *uint16 field", func(t *testing.T) {
		type parsedType struct {
			Field *uint16 `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing value for bool field", func(t *testing.T) {
		type parsedType struct {
			Field bool `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing value for *bool field", func(t *testing.T) {
		type parsedType struct {
			Field *bool `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing value for time.Time field", func(t *testing.T) {
		type parsedType struct {
			Field time.Time `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing value for *time.Time  field", func(t *testing.T) {
		type parsedType struct {
			Field *time.Time `mpd_prefix:"field"`
		}
		lines := []string{"field: asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing from input line with zero colons", func(t *testing.T) {
		type parsedType struct {
			Field int `mpd_prefix:"field"`
		}
		lines := []string{"field asdf"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
	})
	t.Run("error parsing with unsupported field type", func(t *testing.T) {
		type parsedType struct {
			Field uint64 `mpd_prefix:"field"`
		}
		lines := []string{"field: 111"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnsupportedFieldType))
	})
	t.Run("error parsing with unsupported field type (ptr)", func(t *testing.T) {
		type parsedType struct {
			Field *uint64 `mpd_prefix:"field"`
		}
		lines := []string{"field: 111"}
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnsupportedFieldType))
	})
	t.Run("wrong target type", func(t *testing.T) {
		type parsedType interface{}
		var lines []string
		_, err := ParseSingleValue[parsedType](lines)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrTargetTypeMustBeStruct))
	})
}

func TestParseMultiValue(t *testing.T) {
	t.Run("successful parsing to struct with all supported field types", func(t *testing.T) {
		expectedSlice := make([]ParsedType, 3)
		var parsedList []string
		for i := 0; i < 3; i++ {
			faker.FakeData(&expectedSlice[i])
			list := toMpdResponse(&expectedSlice[i])
			parsedList = append(parsedList, list...)
		}
		result, err := ParseMultiValue[ParsedType](parsedList)
		assert.NoError(t, err)
		assert.Equal(t, expectedSlice, result)
	})
	t.Run("parsing field from input with empty lines", func(t *testing.T) {
		type parsedType struct {
			Field string `mpd_prefix:"field" is_new_element_prefix:"true"`
		}
		expected := []parsedType{{Field: "asdf"}, {Field: "fdsa"}}
		lines := []string{"field: asdf", "", "field: fdsa", ""}
		actual, err := ParseMultiValue[parsedType](lines)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("parsing with more than one field marked with is_new_element_prefix:\"true\"", func(t *testing.T) {
		type parsedType struct {
			Field1 string `mpd_prefix:"field1" is_new_element_prefix:"true"`
			Field2 string `mpd_prefix:"field2" is_new_element_prefix:"true"`
			Field3 string `mpd_prefix:"field3"`
		}
		expected := []parsedType{{Field1: "field1", Field3: "field3_1"}, {Field2: "field2", Field3: "field3_2"}}
		lines := []string{"field1: field1", "field3: field3_1", "field2: field2", "field3: field3_2"}
		actual, err := ParseMultiValue[parsedType](lines)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("No field in the target struct is marked with is_new_element_prefix:\"true\"", func(t *testing.T) {
		type targetStruct struct {
			//lint:ignore U1000 ignore
			fieldOne string `mpd_prefix:"string_field1"`
			//lint:ignore U1000 ignore
			fieldTwo string `mpd_prefix:"string_field2"`
		}
		list := []string{
			"string_field1: aaa",
			"string_field2: bbb",
			"string_field1: ccc",
			"string_field2: dddd",
		}
		_, err := ParseMultiValue[targetStruct](list)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoFieldMarkedAsNewElement))
	})
	t.Run("wrong target type", func(t *testing.T) {
		type parsedType interface{}
		var lines []string
		_, err := ParseMultiValue[parsedType](lines)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrTargetTypeMustBeStruct))
	})
	t.Run("error parsing with unsupported field type", func(t *testing.T) {
		type parsedType struct {
			Field uint64 `mpd_prefix:"field" is_new_element_prefix:"true"`
		}
		lines := []string{"field: 100"}
		_, err := ParseMultiValue[parsedType](lines)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnsupportedFieldType))
	})

}

func toMpdResponse(value *ParsedType) []string {
	value.DateField = value.DateField.Round(time.Second)
	datePtrValue := (*value.DatePtrField).Round(time.Second)
	value.DatePtrField = &datePtrValue
	value.IntPtrFieldNilValue = nil
	value.Uint16PtrFieldNilValue = nil
	value.StringPtrFieldNilValue = nil
	value.BoolPtrValueNilValue = nil
	value.DatePtrFieldNilValue = nil
	return []string{
		fmt.Sprintf("int_field: %d", value.IntField),
		fmt.Sprintf("uint16_field: %d", value.Uint16Field),
		fmt.Sprintf("string_field: %s", value.StringField),
		fmt.Sprintf("bool_field: %d", b2int(value.BoolValue)),
		fmt.Sprintf("date_field: %s", value.DateField.Format(time.RFC3339)),
		fmt.Sprintf("int_ptr_field: %d", *value.IntPtrField),
		fmt.Sprintf("uint16_ptr_field: %d", *value.Uint16PtrField),
		fmt.Sprintf("string_ptr_field: %s", *value.StringPtrField),
		fmt.Sprintf("bool_ptr_field: %d", b2int(*value.BoolPtrValue)),
		fmt.Sprintf("date_ptr_field: %s", value.DatePtrField.Format(time.RFC3339)),
	}
}

func b2int(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
