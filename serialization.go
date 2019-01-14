package dynstruct

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/json-iterator/go"
	"github.com/nextzhou/dynstruct/internal/jsonscan"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (v Value) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')
	for i, field := range v.t.fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('"')
		buf.WriteString(field.name)
		buf.WriteString(`":`)
		d, err := json.Marshal(v.value[field.name])
		if err != nil {
			return nil, err
		}
		buf.Write(d)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (v *Value) UnmarshalJSON(data []byte) error {
	if v.t == nil {
		return makeUnknownTypeError()
	}
	kvs, err := jsonscan.Scan(data)
	if err != nil {
		return err
	}
	for _, field := range v.t.fields {
		var d []byte
		// TODO optimize
		for _, kv := range kvs {
			if kv.Key == field.name {
				d = kv.Value
				break
			}
		}
		if len(d) == 0 {
			v.value[field.name] = v.t.zeroValue.value[field.name]
			continue
		}
		fv, err := unmarshal(field.t, d)
		if err != nil {
			return err
		}
		v.value[field.name] = fv
	}
	return nil
}

func unmarshal(t reflect.Type, data []byte) (interface{}, error) {
	switch t.Kind() {
	case reflect.Int:
		return strconv.Atoi(string(data))
	case reflect.Int8:
		n, err := strconv.ParseInt(string(data), 10, 8)
		return int8(n), err
	case reflect.Int16:
		n, err := strconv.ParseInt(string(data), 10, 16)
		return int16(n), err
	case reflect.Int32:
		n, err := strconv.ParseInt(string(data), 10, 32)
		return int32(n), err
	case reflect.Int64:
		n, err := strconv.ParseInt(string(data), 10, 64)
		return int64(n), err
	case reflect.Uint:
		return strconv.Atoi(string(data))
	case reflect.Uint8:
		n, err := strconv.ParseUint(string(data), 10, 8)
		return uint8(n), err
	case reflect.Uint16:
		n, err := strconv.ParseUint(string(data), 10, 16)
		return uint16(n), err
	case reflect.Uint32:
		n, err := strconv.ParseUint(string(data), 10, 32)
		return uint32(n), err
	case reflect.Uint64:
		n, err := strconv.ParseUint(string(data), 10, 64)
		return uint64(n), err
	case reflect.String:
		if bytes.IndexByte(data, byte('\\')) == -1 {
			l := len(data)
			return string(data[1 : l-1]), nil
		}
		var s string
		err := json.Unmarshal(data, &s)
		return s, err
	case reflect.Float32:
		n, err := strconv.ParseFloat(string(data), 32)
		return float32(n), err
	case reflect.Float64:
		return strconv.ParseFloat(string(data), 64)
	}
	v := reflect.New(t)
	err := json.Unmarshal(data, v.Interface())
	if err != nil {
		return nil, err
	}
	return v.Elem().Interface(), nil
}
