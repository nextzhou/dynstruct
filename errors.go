package dynstruct

import (
	"fmt"
	"reflect"
)

type unmatchedTypeError struct {
	t        string
	field    string
	expected string
	got      string
}

func makeUnmatchedTypeError(t *DynStruct, field string, expected, got reflect.Type) error {
	return unmatchedTypeError{
		t:        t.String(),
		field:    field,
		expected: expected.String(),
		got:      got.String(),
	}
}

func (e unmatchedTypeError) Error() string {
	return fmt.Sprintf("field %#v of type %#v unmatched type: expected %#v, got %#v", e.field, e.t, e.expected, e.got)
}

type missingFieldError struct {
	t     string
	field string
}

func makeMissingFieldError(t *DynStruct, field string) error {
	return missingFieldError{t: t.String(), field: field}
}

func (e missingFieldError) Error() string {
	return fmt.Sprintf("type %#v missing field %#v", e.t, e.field)
}

type unknownTypeError struct{}

func makeUnknownTypeError() error {
	return unknownTypeError{}
}

func (e unknownTypeError) Error() string {
	return "unknown type"
}

type invalidNameError struct {
	t    string
	name string
}

func makeInvalidNameError(t, name string) error {
	return invalidNameError{t: t, name: name}
}

func (e invalidNameError) Error() string {
	return fmt.Sprintf("invalid %s name: %#v", e.t, e.name)
}

type recallError struct {
	f string
}

func makeRecallError(f string) error {
	return recallError{f: f}
}

func (e recallError) Error() string {
	return fmt.Sprintf("call function %#v more than once", e.f)
}

type repeatedNameError struct {
	t    string
	name string
}

func makeRepeatedNameError(t, name string) error {
	return repeatedNameError{t: t, name: name}
}

func (e repeatedNameError) Error() string {
	return fmt.Sprintf("repeated %s name: %#v", e.t, e.name)
}

type nilFieldTypeError struct {
	field string
}

func makeNilTypeError(field string) error {
	return nilFieldTypeError{field: field}
}

func (e nilFieldTypeError) Error() string {
	return fmt.Sprintf("type of field %#v is nil", e.field)
}
