package dynstruct

import (
	"reflect"
)

type DynStruct struct {
	pkg, name  string
	fullName   string
	fieldIndex map[string]reflect.Type
	zeroValue  Value
	fields     []field
}

type field struct {
	name string
	t    reflect.Type
}

func makeField(name string, t reflect.Type) field {
	return field{name: name, t: t}
}

func (ds *DynStruct) New() Value {
	return ds.zeroValue.Copy()
}

func (ds *DynStruct) newWithoutInit() Value {
	return Value{
		t:     ds,
		value: make(map[string]interface{}, len(ds.fields)),
	}
}

func (ds *DynStruct) NewFromMapStrictly(m map[string]interface{}) (Value, error) {
	value := ds.New()
	for field, fv := range m {
		if t, ok := ds.fieldIndex[field]; ok {
			if reflect.TypeOf(fv) != t {
				return value, makeUnmatchedTypeError(ds, field, t, reflect.TypeOf(fv))
			}
			value.value[field] = fv
		} else {
			return value, makeMissingFieldError(ds, field)
		}
	}
	return value, nil
}

func (ds *DynStruct) NewFromMap(m map[string]interface{}) (Value, error) {
	value := ds.New()
	for field, fv := range m {
		if t, ok := ds.fieldIndex[field]; ok {
			if reflect.TypeOf(fv) != t {
				return value, makeUnmatchedTypeError(ds, field, t, reflect.TypeOf(fv))
			}
			value.value[field] = fv
		}
	}
	return value, nil
}

func (ds *DynStruct) NewFromMapUnsafely(m map[string]interface{}) Value {
	return Value{
		t:     ds,
		value: m,
	}
}

func (ds DynStruct) String() string {
	return ds.fullName
}
