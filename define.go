package dynstruct

import (
	"reflect"
	"runtime"
	"strings"
)

type definer struct {
	result DynStruct
	err    error
}

func (d *definer) AddField(name string, template interface{}) *definer {
	if d.err != nil {
		return d
	}

	if !isValidIdent(name) {
		d.err = makeInvalidNameError("field", name)
		return d
	}

	if _, ok := d.result.fieldIndex[name]; ok {
		d.err = makeRepeatedNameError("field", name)
		return d
	}

	if template == nil {
		d.err = makeNilTypeError(name)
		return d
	}

	var t reflect.Type
	if tt, ok := template.(reflect.Type); ok {
		t = tt
	} else {
		t = reflect.TypeOf(template)
	}

	field := makeField(name, t)
	d.result.fieldIndex[name] = t
	d.result.fields = append(d.result.fields, field)
	return d
}

func (d *definer) Finish() (DynStruct, error) {
	if d.err != nil {
		return d.result, d.err
	}
	d.result.zeroValue = d.result.newWithoutInit()
	for _, field := range d.result.fields {
		d.result.zeroValue.value[field.name] = reflect.New(field.t).Elem().Interface()
	}
	d.err = makeRecallError("definer.Finish()")
	return d.result, nil
}

func Define(name string) *definer {
	if !isValidIdent(name) {
		return &definer{err: makeInvalidNameError("type", name)}
	}
	pkg := getPkgName()
	return &definer{
		result: DynStruct{
			pkg:        pkg,
			name:       name,
			fullName:   pkg + "." + name,
			fieldIndex: make(map[string]reflect.Type),
		},
	}
}

func getPkgName() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	name := f.Name()
	name = name[strings.LastIndexByte(name, '/')+1:]
	return name[:strings.IndexByte(name, '.')]
}

func isValidIdent(ident string) bool {
	if len(ident) == 0 {
		return false
	}
	c := ident[0]
	if isNumber(rune(c)) {
		return false
	}
	for _, c := range ident {
		if !isValidIdentChar(c) {
			return false
		}
	}
	return true
}

func isValidIdentChar(c rune) bool {
	return isNumber(c) || c == '_' || isChar(c)
}

func isNumber(c rune) bool {
	return '0' <= c && c <= '9'
}

func isChar(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
