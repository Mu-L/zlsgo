package ztype

import (
	"fmt"
	"strconv"
	"testing"

	zls "github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zjson"
)

type st interface {
	String() string
	Set(string)
}

type (
	type1 struct {
		A  int
		B  string
		C1 float32
	}
	type2 struct {
		D bool
		E *uint
		F []string
		G map[string]int
		type1
		S1 type1
		S2 *type1
	}
	type3 struct {
		Name string
	}
)

var ni interface{}

type j struct {
	Name string
	Key  string
	Age  int `json:"age"`
}

var (
	str          = "123"
	i            = 123
	i8   int8    = 123
	i16  int16   = 123
	i32  int32   = 123
	i64  int64   = 123
	ui8  uint8   = 123
	ui   uint    = 123
	ui16 uint16  = 123
	ui32 uint32  = 123
	ui64 uint64  = 123
	f3   float32 = 123
	f6   float64 = 123
	b            = true
)

func (s *j) String() string {
	return ToString(s.Key)
}

func (s *j) Set(v string) {
	s.Key = v
}

func TestTo(t *testing.T) {
	T := zls.NewTest(t)
	var sst st = new(j)
	sst.Set(str)
	jj := j{Name: "123"}

	T.Equal([]byte(str), ToByte(str))
	T.Equal([]byte(str), ToByte(i))

	T.Equal(0, ToInt(ni))
	T.Equal(i, ToInt(str))
	T.Equal(i, ToInt(i))
	T.Equal(i8, ToInt8(str))
	T.Equal(i8, ToInt8(i8))
	T.Equal(i16, ToInt16(str))
	T.Equal(i16, ToInt16(i16))
	T.Equal(i32, ToInt32(str))
	T.Equal(i32, ToInt32(i32))

	T.Equal(i64, ToInt64(str))
	T.Equal(i64, ToInt64(i))
	T.Equal(i64, ToInt64(i8))
	T.Equal(i64, ToInt64(i16))
	T.Equal(i64, ToInt64(i32))
	T.Equal(i64, ToInt64(i64))
	T.Equal(i64, ToInt64(ui8))
	T.Equal(i64, ToInt64(ui))
	T.Equal(i64, ToInt64(ui16))
	T.Equal(i64, ToInt64(ui32))
	T.Equal(i64, ToInt64(ui64))
	T.Equal(i64, ToInt64(f3))
	T.Equal(i64, ToInt64(f6))
	// 无法转换直接换成0
	T.Equal(ToInt64(0), ToInt64(jj))
	T.Equal(i64, ToInt64("0x7b"))
	T.Equal(i64, ToInt64("0173"))
	T.Equal(ToInt64(1), ToInt64(b))
	T.Equal(ToInt64(0), ToInt64(false))

	T.Equal(ToUint(0), ToUint(ni))
	T.Equal(ui, ToUint(str))
	T.Equal(ui, ToUint(ui))
	T.Equal(ui8, ToUint8(str))
	T.Equal(ui8, ToUint8(ui8))
	T.Equal(ui16, ToUint16(str))
	T.Equal(ui16, ToUint16(ui16))
	T.Equal(ui32, ToUint32(str))
	T.Equal(ui32, ToUint32(ui32))

	T.Equal(ui64, ToUint64(i64))
	T.Equal(ui64, ToUint64(str))
	T.Equal(ui64, ToUint64(i))
	T.Equal(ui64, ToUint64(i8))
	T.Equal(ui64, ToUint64(i16))
	T.Equal(ui64, ToUint64(i32))
	T.Equal(ui64, ToUint64(ui))
	T.Equal(ui64, ToUint64(ui8))
	T.Equal(ui64, ToUint64(ui16))
	T.Equal(ui64, ToUint64(ui32))
	T.Equal(ui64, ToUint64(ui64))
	T.Equal(ui64, ToUint64(f3))
	T.Equal(ui64, ToUint64(f6))
	// 无法转换直接换成0
	T.Equal(ToUint64(0), ToUint64(jj))
	T.Equal(ui64, ToUint64("0x7b"))
	T.Equal(ui64, ToUint64("0173"))
	T.Equal(ToUint64(1), ToUint64(b))
	T.Equal(ToUint64(0), ToUint64(false))

	T.Equal(str, ToString(sst))
	T.Equal("", ToString(ni))
	T.Equal("true", ToString(b))
	T.Equal(str, ToString(str))
	T.Equal(str, ToString(i8))
	T.Equal(str, ToString(ui))
	T.Equal(str, ToString(i))
	T.Equal(str, ToString(i8))
	T.Equal(str, ToString(i16))
	T.Equal(str, ToString(i32))
	T.Equal(str, ToString(i64))
	T.Equal(str, ToString(ui8))
	T.Equal(str, ToString(ui16))
	T.Equal(str, ToString(ui32))
	T.Equal(str, ToString(ui64))
	T.Equal(str, ToString(f6))
	T.Equal(str, ToString(f3))
	T.Equal(str, ToString(ToByte(i)))
	T.Equal("{\"Name\":\"123\",\"Key\":\"\",\"age\":0}", ToString(jj))
	T.Equal(f6, ToFloat64(i))
	T.Equal(f6, ToFloat64(f3))
	T.Equal(f6, ToFloat64(f6))
	T.Equal(ToFloat64(0), ToFloat64(ni))

	T.Equal(f3, ToFloat32(i))
	T.Equal(f3, ToFloat32(f3))
	T.Equal(f3, ToFloat32(f6))
	T.Equal(ToFloat32(0), ToFloat32(ni))

	T.Equal(true, ToBool(b))
	T.Equal(true, ToBool(str))
	T.Equal(false, ToBool(ni))

}

// func BenchmarkToString0(b *testing.B) {
// s := 123
// for i := 0; i < b.N; i++ {
// _ = strconv.Itoa(s)
// }
// }

func BenchmarkToString1(b *testing.B) {
	s := true
	for i := 0; i < b.N; i++ {
		_ = ToString(s)
	}
	// type a struct {
	// Na string `json:"n"`
	// }

	// n := &a{
	// Na: "hi",
	// }
	// b.Log(ToString(n))
	// b.Log(ToString(*n))
}
func BenchmarkToString2(b *testing.B) {
	s := true
	for i := 0; i < b.N; i++ {
		_ = String(s)
	}
	// type a struct {
	// Na string `json:"n"`
	// }

	// n := &a{
	// Na: "hi",
	// }
	// b.Log(String(n))
	// b.Log(String(*n))
}
func String(val interface{}) string {
	if val == nil {
		return ""
	}

	switch t := val.(type) {
	case bool:
		return strconv.FormatBool(t)
	case int:
		return strconv.FormatInt(int64(t), 10)
	case int8:
		return strconv.FormatInt(int64(t), 10)
	case int16:
		return strconv.FormatInt(int64(t), 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case []byte:
		return string(t)
	case string:
		return t
	default:
		return fmt.Sprintf("%v", val)
	}
}

func TestStructToMap(tt *testing.T) {
	e := uint(8)
	t := zls.NewTest(tt)
	v := &type2{
		D: true,
		E: &e,
		F: []string{"f1", "f2"},
		G: map[string]int{"G1": 1, "G2": 2},
		type1: type1{
			A: 1,
			B: "type1",
		},
		S1: type1{
			A: 2,
			B: "S1",
		},
		S2: &type1{
			A: 3,
			B: "Ss",
		},
	}
	r := StructToMap(v)
	t.Log(v, r)
	j, err := zjson.Marshal(r)
	t.EqualNil(err)
	t.EqualExit(`{"D":true,"E":8,"F":["f1","f2"],"G":{"G1":1,"G2":2},"S1":{"A":2,"B":"S1"},"S2":{"A":3,"B":"Ss"},"type1":{"A":1,"B":"type1"}}`, string(j))

	v2 := []string{"1", "2", "more"}
	r = StructToMap(v2)
	t.Log(v2, r)
	j, err = zjson.Marshal(v2)
	t.EqualNil(err)
	t.EqualExit(`["1","2","more"]`, string(j))

	v3 := "ok"
	r = StructToMap(v3)
	t.Log(v3, r)
	j, err = zjson.Marshal(v3)
	t.EqualNil(err)
	t.EqualExit(`"ok"`, string(j))
}
