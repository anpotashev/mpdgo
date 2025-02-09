package parser

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ParsedType struct {
	IntField    int       `mpd_prefix:"int_field" is_new_element_prefix:"true"`
	Uint16Field uint16    `mpd_prefix:"uint16_field"`
	StringField string    `mpd_prefix:"string_field"`
	BoolValue   bool      `mpd_prefix:"bool_field"`
	DateField   time.Time `mpd_prefix:"date_field"`
}

func TestParseSingleValue(t *testing.T) {
	t.Run("happy pass", func(t *testing.T) {
		var expectedValue ParsedType
		faker.FakeData(&expectedValue)
		expectedValue.DateField = expectedValue.DateField.Round(time.Second)
		list := []string{
			fmt.Sprintf("int_field: %d", expectedValue.IntField),
			fmt.Sprintf("uint16_field: %d", expectedValue.Uint16Field),
			fmt.Sprintf("string_field: %s", expectedValue.StringField),
			fmt.Sprintf("bool_field: %d", b2int(expectedValue.BoolValue)),
			fmt.Sprintf("date_field: %s", expectedValue.DateField.Format(time.RFC3339)),
		}
		result, err := ParseSingleValue[ParsedType](list)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
	})
}

func TestParseMultiValue(t *testing.T) {
	t.Run("happy pass", func(t *testing.T) {
		expectedSlice := make([]ParsedType, 3)
		var parsedList []string
		for i := 0; i < 3; i++ {
			faker.FakeData(&expectedSlice[i])
			expectedSlice[i].DateField = expectedSlice[i].DateField.Round(time.Second)
			parsedList = append(parsedList, fmt.Sprintf("int_field: %d", expectedSlice[i].IntField))
			parsedList = append(parsedList, fmt.Sprintf("uint16_field: %d", expectedSlice[i].Uint16Field))
			parsedList = append(parsedList, fmt.Sprintf("string_field: %s", expectedSlice[i].StringField))
			parsedList = append(parsedList, fmt.Sprintf("bool_field: %d", b2int(expectedSlice[i].BoolValue)))
			parsedList = append(parsedList, fmt.Sprintf("date_field: %s", expectedSlice[i].DateField.Format(time.RFC3339)))
		}
		result, err := ParseMultiValue[ParsedType](parsedList)
		assert.NoError(t, err)
		assert.Equal(t, expectedSlice, result)
	})
}
func b2int(b bool) int8 {
	if b {
		return 1
	}
	return 0
}
