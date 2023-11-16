package codefixture

import (
	"fmt"
	"reflect"
)

type (
	NotPointerError struct {
		Type reflect.Type
	}
	NotStructError struct {
		Type reflect.Type
	}
	ModelRefNotFoundError struct {
		Ref ModelRef
	}
	InvalidType struct {
		Target any
	}
)

func (e *NotPointerError) Error() string {
	return fmt.Sprintf("type %v is not a pointer", e.Type)
}
func (e *NotStructError) Error() string {
	return fmt.Sprintf("type %v is not a struct", e.Type)
}
func (e *ModelRefNotFoundError) Error() string {
	return fmt.Sprintf("model ref %v not found", e.Ref)
}
func (e *InvalidType) Error() string {
	return fmt.Sprintf("invalid type: %T", e.Target)
}

func NewNotPointerError(t reflect.Type) error {
	return &NotPointerError{Type: t}
}
func NewNotStructError(t reflect.Type) error {
	return &NotStructError{Type: t}
}
func NewModelRefNotFoundError(ref ModelRef) error {
	return &ModelRefNotFoundError{Ref: ref}
}
func NewInvalidTypeError(target any) error {
	return &InvalidType{Target: target}
}
