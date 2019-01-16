package dynstruct

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/nextzhou/dynstruct/internal/jsonscan"

	"github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
)

var data = []byte(`{"Field1":"abcdefg","Field2":1234,"Field3":1234.5678}`)

func TestDefine(t *testing.T) {
	Convey("define", t, func() {
		Convey("invalid type name", func() {
			_, err := Define("").Finish()
			So(err, ShouldBeError, `invalid type name: ""`)

			_, err = Define(" ").Finish()
			So(err, ShouldBeError, `invalid type name: " "`)

			_, err = Define(" abc").Finish()
			So(err, ShouldBeError, `invalid type name: " abc"`)

			_, err = Define("abc def").Finish()
			So(err, ShouldBeError, `invalid type name: "abc def"`)

			_, err = Define("abc ").Finish()
			So(err, ShouldBeError, `invalid type name: "abc "`)

			_, err = Define("123").Finish()
			So(err, ShouldBeError, `invalid type name: "123"`)

			_, err = Define("123abc").Finish()
			So(err, ShouldBeError, `invalid type name: "123abc"`)

			_, err = Define("类型").Finish()
			So(err, ShouldBeError, `invalid type name: "类型"`)

			_, err = Define("😈").Finish()
			So(err, ShouldBeError, `invalid type name: "😈"`)
		})

		Convey("valid type name", func() {
			typ, err := Define("abc").Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct.abc")

			typ, err = Define("Abc").Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct.Abc")

			typ, err = Define("ABC").Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct.ABC")

			typ, err = Define("Abc1").Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct.Abc1")

			typ, err = Define("abc_def").Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct.abc_def")

			typ, err = Define("_abc").Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct._abc")

			typ, err = Define("abc_").Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct.abc_")
		})

		Convey("invalid field name", func() {
			st := reflect.TypeOf("")
			_, err := Define("Abc").AddField("", st).Finish()
			So(err, ShouldBeError, `invalid field name: ""`)

			_, err = Define("Abc").AddField(" ", st).Finish()
			So(err, ShouldBeError, `invalid field name: " "`)

			_, err = Define("Abc").AddField(" abc", st).Finish()
			So(err, ShouldBeError, `invalid field name: " abc"`)

			_, err = Define("Abc").AddField("abc def", st).Finish()
			So(err, ShouldBeError, `invalid field name: "abc def"`)

			_, err = Define("Abc").AddField("abc ", st).Finish()
			So(err, ShouldBeError, `invalid field name: "abc "`)

			_, err = Define("Abc").AddField("123", st).Finish()
			So(err, ShouldBeError, `invalid field name: "123"`)

			_, err = Define("Abc").AddField("123abc", st).Finish()
			So(err, ShouldBeError, `invalid field name: "123abc"`)

			_, err = Define("Abc").AddField("类型", st).Finish()
			So(err, ShouldBeError, `invalid field name: "类型"`)

			_, err = Define("Abc").AddField("😈", st).Finish()
			So(err, ShouldBeError, `invalid field name: "😈"`)

			_, err = Define("Abc").AddField("field", st).Finish()
			So(err, ShouldBeNil)
			_, err = Define("Abc").AddField("field", st).AddField("", st).Finish()
			So(err, ShouldBeError, `invalid field name: ""`)
			_, err = Define("Abc").AddField("", st).AddField("field", st).Finish()
			So(err, ShouldBeError, `invalid field name: ""`)
		})

		Convey("reuse definer", func() {
			definer := Define("Abc").AddField("Field1", "")
			_, err := definer.Finish()
			So(err, ShouldBeNil)

			definer = Define("Abc").AddField("Field1", reflect.TypeOf(""))
			_, err = definer.Finish()
			So(err, ShouldBeNil)

			_, err = definer.AddField("Field2", "").Finish()
			So(err, ShouldBeError, `call function "definer.Finish()" more than once`)
		})

		Convey("repeated field", func() {
			_, err := Define("Abc").
				AddField("field1", "").
				AddField("field1", "").
				Finish()
			So(err, ShouldBeError, `repeated field name: "field1"`)

			_, err = Define("Abc").
				AddField("field1", "").
				AddField("field1", 0).
				Finish()
			So(err, ShouldBeError, `repeated field name: "field1"`)
		})

		Convey("nil field type", func() {
			_, err := Define("Abc").AddField("field", nil).Finish()
			So(err, ShouldBeError, `type of field "field" is nil`)

			_, err = Define("Abc").AddField("field", reflect.TypeOf(nil)).Finish()
			So(err, ShouldBeError, `type of field "field" is nil`)
		})

		Convey("valid define", func() {
			typ, err := Define("Abc").
				AddField("StrField", reflect.TypeOf("")).
				AddField("IntField", reflect.TypeOf(int(0))).
				AddField("FloatField", reflect.TypeOf(float64(0))).
				AddField("StructField", reflect.TypeOf(time.Now())).
				Finish()
			So(err, ShouldBeNil)
			So(typ.String(), ShouldEqual, "dynstruct.Abc")
		})
	})
}

func TestFormat(t *testing.T) {
	type Sub struct {
		Abc string
		Def []uint
	}
	Convey("format", t, func() {
		typ, err := Define("Abc").
			AddField("StrField", reflect.TypeOf("")).
			AddField("IntField", reflect.TypeOf(int(0))).
			AddField("FloatField", reflect.TypeOf(float64(0))).
			AddField("StructField", reflect.TypeOf(Sub{})).
			Finish()
		So(err, ShouldBeNil)
		So(fmt.Sprint(typ), ShouldEqual, "dynstruct.Abc")

		val := typ.New()
		val.Set("StrField", "abc")
		val.Set("IntField", 123)
		val.Set("FloatField", 123.456)
		val.Set("StructField", Sub{Abc: "sub abc", Def: []uint{1, 2, 3}})

		So(fmt.Sprint(val), ShouldEqual, "{abc 123 123.456 {sub abc [1 2 3]}}")
		So(fmt.Sprintf("%v", val), ShouldEqual, "{abc 123 123.456 {sub abc [1 2 3]}}")
		So(fmt.Sprintf("%+v", val), ShouldEqual,
			"{StrField:abc IntField:123 FloatField:123.456 StructField:{Abc:sub abc Def:[1 2 3]}}")
		So(fmt.Sprintf("%#v", val), ShouldEqual,
			`dynstruct.Abc{StrField:"abc", IntField:123, FloatField:123.456, StructField:dynstruct.Sub{Abc:"sub abc", Def:[]uint{0x1, 0x2, 0x3}}}`)
		So(fmt.Sprintf("%z", val), ShouldEqual,
			"{%!z(string=abc) %!z(int=123) %!z(float64=123.456) {%!z(string=sub abc) [%!z(uint=1) %!z(uint=2) %!z(uint=3)]}}")
	})

}

func TestNewValue(t *testing.T) {
	Convey("new value", t, func() {
		typ, err := Define("Abc").AddField("int", reflect.TypeOf(int(0))).Finish()
		So(err, ShouldBeNil)

		val := typ.New()
		So(val.Get("int"), ShouldBeZeroValue)
		val.Set("int", 123)
		So(val.Get("int"), ShouldEqual, 123)

		val2 := typ.New()
		So(val2.Get("int"), ShouldBeZeroValue)
	})
}

func TestSetValue(t *testing.T) {
	Convey("set value", t, func() {
		typ, err := Define("Abc").
			AddField("int8", reflect.TypeOf(int8(0))).
			AddField("pint", reflect.TypeOf((*int)(nil))).
			AddField("str", reflect.TypeOf("")).
			Finish()
		So(err, ShouldBeNil)
		val := typ.New()

		So(func() { val.Set("int_1", 1) }, ShouldPanic)
		So(func() { val.Set("int8", "123") }, ShouldPanic)
		So(func() { val.Set("pint", 123) }, ShouldPanic)

		p := new(int)
		*p = 456
		val.Set("int8", int8(123)) // TODO support number type conversion
		val.Set("pint", p)
		val.Set("str", "789")
		data, _ := json.Marshal(val)
		So(string(data), ShouldEqual, `{"int8":123,"pint":456,"str":"789"}`)

		val.Set("pint", nil)
		data, _ = json.Marshal(val)
		So(string(data), ShouldEqual, `{"int8":123,"pint":null,"str":"789"}`)

		bigNumber := 123456
		val.Set("int8", int8(bigNumber))
		So(val.Get("int8"), ShouldEqual, int8(bigNumber))
	})
}

func TestInterfaceField(t *testing.T) {
	Convey("interface field", t, func() {
		anyType := reflect.TypeOf((*interface{})(nil)).Elem()
		writerType := reflect.TypeOf((*io.Writer)(nil)).Elem()
		typ, err := Define("Abc").
			AddField("Any", anyType).
			AddField("Writer", writerType).
			Finish()
		So(err, ShouldBeNil)
		val := typ.New()

		Convey("empty interface", func() {
			var i int
			var s string
			var f float64
			var any interface{}

			val.Set("Any", 123)
			val.Scan("Any", &i)
			So(i, ShouldEqual, 123)
			val.Scan("Any", &any)
			So(any, ShouldEqual, 123)
			So(func() { val.Scan("Any", &s) }, ShouldPanic)

			val.Set("Any", "haha")
			val.Scan("Any", &s)
			So(s, ShouldEqual, "haha")
			val.Scan("Any", &any)
			So(any, ShouldEqual, "haha")
			So(func() { val.Scan("Any", &i) }, ShouldPanic)

			val.Set("Any", 1.23)
			val.Scan("Any", &f)
			So(f, ShouldAlmostEqual, 1.23)
			val.Scan("Any", &any)
			So(any, ShouldAlmostEqual, 1.23)
			So(func() { val.Scan("Any", &i) }, ShouldPanic)
		})

		Convey("non-empty interface", func() {
			var w io.Writer
			var b, bb *bytes.Buffer
			b = bytes.NewBuffer(nil)

			val.Set("Writer", b)
			val.Scan("Writer", &bb)
			So(bb, ShouldPointTo, b)
			val.Scan("Writer", &w)
			So(w, ShouldPointTo, b)
			So(val.Get("Writer"), ShouldPointTo, b)

			So(func() { val.Set("Writer", 0) }, ShouldPanic)
		})
	})
}

func TestScanValue(t *testing.T) {
	Convey("scan value", t, func() {
		typ, err := Define("Abc").
			AddField("int", reflect.TypeOf(int(0))).
			AddField("pint", reflect.TypeOf((*int)(nil))).
			AddField("str", reflect.TypeOf("")).
			Finish()
		So(err, ShouldBeNil)
		val := typ.New()

		p := new(int)
		*p = 456
		val.Set("int", int(123))
		val.Set("pint", p)
		val.Set("str", "789")

		var n int
		var pn *int
		var s string

		Convey("normal scan", func() {
			val.Scan("int", &n)
			val.Scan("pint", &pn)
			val.Scan("str", &s)

			So(n, ShouldEqual, 123)
			So(*pn, ShouldEqual, 456)
			So(s, ShouldEqual, "789")

			val.Set("pint", nil)
			val.Scan("pint", &pn)
			So(pn, ShouldBeNil)
		})

		Convey("invalid scan", func() {
			So(func() { val.Scan("missingField", &n) }, ShouldPanic)
			So(func() { val.Scan("int", &pn) }, ShouldPanic)

			var u8 uint8
			So(func() { val.Scan("int", &u8) }, ShouldPanic)

			So(func() { val.Scan("int", n) }, ShouldPanic)
		})
	})
}

func TestGet(t *testing.T) {
	Convey("get", t, func() {
		typ, err := Define("Abc").
			AddField("int", reflect.TypeOf(int(0))).
			AddField("pint", reflect.TypeOf((*int)(nil))).
			AddField("str", reflect.TypeOf("")).
			Finish()
		So(err, ShouldBeNil)
		val := typ.New()

		So(val.Get("int").(int), ShouldBeZeroValue)
		So(val.Get("pint").(*int), ShouldBeZeroValue)
		So(val.Get("str").(string), ShouldBeZeroValue)
		So(func() { val.Get("unknownField") }, ShouldPanic)

		p := new(int)
		*p = 456
		val.Set("int", int(123))
		val.Set("pint", p)
		val.Set("str", "789")

		So(val.Get("int").(int), ShouldEqual, 123)
		So(val.Get("pint").(*int), ShouldEqual, p)
		So(val.Get("str").(string), ShouldEqual, "789")
	})
}

func TestJsonMarshal(t *testing.T) {
	Convey("JSON marshal", t, func() {
		typ, err := Define("Abc").
			AddField("int", reflect.TypeOf(int(0))).
			AddField("pint", reflect.TypeOf((*int)(nil))).
			AddField("str", reflect.TypeOf("")).
			Finish()
		So(err, ShouldBeNil)
		val := typ.New()

		data, err := json.Marshal(val)
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, `{"int":0,"pint":null,"str":""}`)

		p := new(int)
		*p = 456
		val.Set("int", int(123))
		val.Set("pint", p)
		val.Set("str", "789")

		data, err = json.Marshal(val)
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, `{"int":123,"pint":456,"str":"789"}`)
	})
}

func TestJsonUnmarshal(t *testing.T) {
	Convey("JSON unmarshal", t, func() {
		typ, err := Define("Abc").
			AddField("int", reflect.TypeOf(int(0))).
			AddField("pint", reflect.TypeOf((*int)(nil))).
			AddField("str", reflect.TypeOf("")).
			Finish()
		So(err, ShouldBeNil)
		val := typ.New()

		data := []byte(`{"int":123,"pint":456,"str":"789"}`)
		err = json.Unmarshal(data, &val)
		So(err, ShouldBeNil)

		So(val.Get("int").(int), ShouldEqual, 123)
		So(*val.Get("pint").(*int), ShouldEqual, 456)
		So(val.Get("str").(string), ShouldEqual, "789")

		data = []byte(`{}`)
		err = json.Unmarshal(data, &val)
		So(err, ShouldBeNil)

		So(val.Get("int").(int), ShouldEqual, 0)
		So(val.Get("pint").(*int), ShouldBeNil)
		So(val.Get("str").(string), ShouldEqual, "")

		data = []byte(`{"int":123,"abc":"ABC"}`)
		err = json.Unmarshal(data, &val)
		So(err, ShouldBeNil)

		So(val.Get("int").(int), ShouldEqual, 123)
		So(val.Get("pint").(*int), ShouldBeNil)
		So(val.Get("str").(string), ShouldEqual, "")

		data = []byte(`{"int":"abc"}`)
		err = json.Unmarshal(data, &val)
		So(err, ShouldNotBeNil)
	})
}

func BenchmarkMarshalJsonStruct(b *testing.B) {
	b.ReportAllocs()
	val := struct {
		Field1 string
		Field2 int
		Field3 float64
	}{
		"abcdefg",
		1234,
		1234.5678,
	}
	for i := 0; i < b.N; i++ {
		json.Marshal(val)
	}
}

func BenchmarkMarshalJsonDynStruct(b *testing.B) {
	b.ReportAllocs()
	typ, _ := Define("T").
		AddField("Field1", reflect.TypeOf("")).
		AddField("Field2", reflect.TypeOf(int(0))).
		AddField("Field3", reflect.TypeOf(float64(0))).Finish()
	val := typ.New()
	val.Set("Field1", "abcdefg")
	val.Set("Field2", int(1234))
	val.Set("Field3", float64(1234.5678))
	for i := 0; i < b.N; i++ {
		json.Marshal(val)
	}
}

func BenchmarkMarshalJsonMap(b *testing.B) {
	b.ReportAllocs()
	val := map[string]interface{}{
		"Field1": "abcdefg",
		"Field2": int(1234),
		"Field3": float64(1234.5678),
	}
	for i := 0; i < b.N; i++ {
		json.Marshal(val)
	}
}

func BenchmarkUnmarshalJsonStruct(b *testing.B) {
	b.ReportAllocs()
	var val struct {
		Field1 string
		Field2 int
		Field3 float64
	}
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &val)
		if err != nil {
			panic(err)
		}
	}
}
func BenchmarkUnmarshalJsonDynStruct(b *testing.B) {
	b.ReportAllocs()
	typ, _ := Define("T").
		AddField("Field1", reflect.TypeOf("")).
		AddField("Field2", reflect.TypeOf(int(0))).
		AddField("Field3", reflect.TypeOf(float64(0))).
		Finish()
	for i := 0; i < b.N; i++ {
		val := typ.New()
		err := json.Unmarshal(data, &val)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUnmarshalJsonDynStructExplicit(b *testing.B) {
	b.ReportAllocs()
	typ, _ := Define("T").
		AddField("Field1", reflect.TypeOf("")).
		AddField("Field2", reflect.TypeOf(int(0))).
		AddField("Field3", reflect.TypeOf(float64(0))).
		Finish()
	for i := 0; i < b.N; i++ {
		val := typ.New()
		err := val.UnmarshalJSON(data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUnmarshalJsonMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var val map[string]interface{}
		err := json.Unmarshal(data, &val)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUnmarshalJsonRawMessageMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var val map[string]jsoniter.RawMessage
		err := json.Unmarshal(data, &val)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUnmarshalJsonScan(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := jsonscan.Scan(data)
		if err != nil {
			panic(err)
		}
	}
}
