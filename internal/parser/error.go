package parser

import (
	"fmt"
	"reflect"
)

var (
	NoFieldMarkedAsNewElementError       = fmt.Errorf("T must have at least one field marked as \"new element prefix\"")
	TargetTypeMustBeStructError          = fmt.Errorf("T must be a struct")
	UnsupportedFieldType                 = fmt.Errorf("unsupported field type")
	FieldParsingError              error = fmt.Errorf("Field parsing error")
)

func NewFieldParsingError(fieldName, value string, fieldVal reflect.Value, err error) error {
	return fmt.Errorf("error parsing field %s: cannot convert %q to %v: %w", fieldName, value, fieldVal.Type(), fmt.Errorf("%w: %v", FieldParsingError, err))
}
