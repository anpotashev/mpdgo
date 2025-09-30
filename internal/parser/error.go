package parser

import (
	"fmt"
	"reflect"
)

var (
	//lint:ignore ST1005 ignore
	ErrNoFieldMarkedAsNewElement = fmt.Errorf("T must have at least one field marked as \"new element prefix\"")
	//lint:ignore ST1005 ignore
	ErrTargetTypeMustBeStruct       = fmt.Errorf("T must be a struct")
	ErrUnsupportedFieldType         = fmt.Errorf("unsupported field type")
	ErrParsingField           error = fmt.Errorf("field parsing error")
)

func NewFieldParsingError(fieldName, value string, fieldVal reflect.Value, err error) error {
	return fmt.Errorf("error parsing field %s: cannot convert %q to %v: %w", fieldName, value, fieldVal.Type(), fmt.Errorf("%w: %v", ErrParsingField, err))
}
