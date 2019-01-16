package dynstruct

import (
	"reflect"
)

type Value struct {
	t     *DynStruct
	value map[string]interface{}
}

func (v Value) Copy() Value {
	newVal := Value{
		t:     v.t,
		value: make(map[string]interface{}, len(v.t.fields)),
	}
	for k, v := range v.value {
		newVal.value[k] = v
	}
	return newVal
}

func (v Value) Set(field string, val interface{}) {
	t, ok := v.t.fieldIndex[field]
	if !ok {
		panic(makeMissingFieldError(v.t, field))
	}
	if !isMatchedType(t, reflect.TypeOf(val)) {
		panic(makeUnmatchedTypeError(v.t, field, t, reflect.TypeOf(val)))
	}
	if val != nil {
		v.value[field] = val
	} else {
		v.value[field] = reflect.New(t).Elem().Interface()
	}
}

func (v Value) UncheckSet(field string, value interface{}) {
	v.value[field] = value
}

func (v Value) Scan(field string, val interface{}) {
	// TODO more exact panic message
	fv, ok := v.value[field]
	if !ok {
		panic(makeMissingFieldError(v.t, field))
	}
	reflect.ValueOf(val).Elem().Set(reflect.ValueOf(fv))
}

func (v Value) UncheckScan(field string, val interface{}) {
	reflect.ValueOf(val).Elem().Set(reflect.ValueOf(v.value[field]))
}

func (v Value) Get(field string) interface{} {
	fv, ok := v.value[field]
	if !ok {
		panic(makeMissingFieldError(v.t, field))
	}
	return fv
}

func (v Value) UncheckGet(field string) interface{} {
	return v.value[field]
}

func isMatchedType(t, vt reflect.Type) bool {
	if t == vt {
		return true
	}
	if t.Kind() == reflect.Ptr && vt == nil {
		return true
	}
	if t.Kind() == reflect.Interface && vt.Implements(t) {
		return true
	}
	return false
}
