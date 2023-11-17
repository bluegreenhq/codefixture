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
	InvalidTypeError struct {
		Target any
	}
	UnexpectedTypeError struct {
		ExpectedType reflect.Type
		ActualType   reflect.Type
	}
	WriterNotFoundError struct {
		Type reflect.Type
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
func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type: %T", e.Target)
}
func (e *UnexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type: expected %v, actual %v", e.ExpectedType, e.ActualType)
}
func (e *WriterNotFoundError) Error() string {
	return fmt.Sprintf("writer not found for type %v", e.Type)
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
	return &InvalidTypeError{Target: target}
}
func NewUnexpectedTypeError(expectedType reflect.Type, actualType reflect.Type) error {
	return &UnexpectedTypeError{ExpectedType: expectedType, ActualType: actualType}
}
func NewWriterNotFoundError(t reflect.Type) error {
	return &WriterNotFoundError{Type: t}
}
