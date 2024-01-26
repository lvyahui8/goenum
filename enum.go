package goenum

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type EnumDefinition interface {
	// Stringer 支持打印输出
	fmt.Stringer
	// Marshaler 支持枚举序列化
	json.Marshaler
	// TextMarshaler 支持枚举序列化
	encoding.TextMarshaler
	// Name 枚举名称，同一类型枚举应该唯一
	Name() string
	// Equals 枚举对比
	Equals(other EnumDefinition) bool
	// Type 实际的枚举类型
	Type() string
	// Ordinal 获取枚举序数
	Ordinal() int
	// Compare 枚举比较方法
	Compare(other EnumDefinition) int
}

type Enum struct {
	name  string
	_type string
	index int
}

func (e Enum) Name() string {
	return e.name
}

func (e Enum) Equals(other EnumDefinition) bool {
	if e.Name() != other.Name() {
		return false
	}
	// 比较类型
	return e._type == other.Type()
}

func (e Enum) String() string {
	return e.Name()
}

func (e Enum) Ordinal() int {
	return e.index
}

func (e Enum) Type() string {
	return e._type
}

func (e Enum) Compare(other EnumDefinition) int {
	return e.Ordinal() - other.Ordinal()
}

func (e Enum) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Name())
}

func (e Enum) MarshalText() (text []byte, err error) {
	return []byte(e.Name()), nil
}

// name2enumsMap name到枚举实例的映射，不同的枚举，name可能冲突，所以value是slice
var name2enumsMap = make(map[string][]EnumDefinition)

// type2enumsMap 枚举类型到所有枚举的映射，key为枚举类型的全路径名称
var type2enumsMap = make(map[string][]EnumDefinition)

// typeIndexMap 类型枚举索引
var typeIndexMap = make(map[string]int)

// NewEnum 新建枚举, 如果枚举（同类型）已经存在，则会抛出panic，禁止重复创建枚举
func NewEnum[T EnumDefinition](name string, src ...T) T {
	if IsValidEnum[T](name) {
		panic("Enum must be unique")
	}
	var t T
	if len(src) > 0 {
		t = src[0]
	}
	v := reflect.ValueOf(t)
	var elem reflect.Value
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			elem = reflect.Indirect(reflect.New(v.Type().Elem()))
		} else {
			elem = reflect.Indirect(v)
		}
	} else {
		elem = reflect.ValueOf(&t).Elem()
	}
	enumFiled := elem.FieldByName(reflect.TypeOf(Enum{}).Name())

	// 获取泛型具体类型名
	tFullName := typeKey(v.Type())
	idx := typeIndexMap[tFullName]
	e := Enum{name: name, _type: tFullName, index: idx}
	if enumFiled.Kind() == reflect.Ptr {
		enumFiled.Set(reflect.ValueOf(&e))
	} else {
		enumFiled.Set(reflect.ValueOf(e))
	}
	typeIndexMap[tFullName] = idx + 1

	type2enumsMap[tFullName] = append(type2enumsMap[tFullName], t)
	name2enumsMap[name] = append(name2enumsMap[name], t)
	return t
}

// ValueOf 根据字符串获取枚举，如果找不到，则返回nil
func ValueOf[T EnumDefinition](name string) (t T, valid bool) {
	enums := name2enumsMap[name]
	for _, e := range enums {
		if v, ok := e.(T); ok {
			return v, true
		}
	}
	return
}

// ValueOfIgnoreCase 忽略大小写获取枚举, 涉及到一次反射调用，性能比ValueOf略差
func ValueOfIgnoreCase[T EnumDefinition](name string) (t T, valid bool) {
	values := Values[T]()
	for _, e := range values {
		if strings.EqualFold(e.Name(), name) {
			return e, true
		}
	}
	return
}

// Values 返回所有可用枚举，返回slice是有序的，按照ordinal排序
func Values[T EnumDefinition]() []T {
	var t T
	var res []T
	tName := typeKey(reflect.TypeOf(t))
	for _, e := range type2enumsMap[tName] {
		if v, ok := e.(T); ok {
			res = append(res, v)
		}
	}
	return res
}

func Size[T EnumDefinition]() int {
	var t T
	tName := typeKey(reflect.TypeOf(t))
	return len(type2enumsMap[tName])
}

// GetEnumMap 获取所有枚举，以name->enum map的形式返回
func GetEnumMap[T EnumDefinition]() map[string]T {
	values := Values[T]()
	res := make(map[string]T)
	for _, e := range values {
		res[e.Name()] = e
	}
	return res
}

// EnumNames 获取一批枚举的名称
func EnumNames[T EnumDefinition](enums ...T) (names []string) {
	for _, e := range enums {
		names = append(names, e.Name())
	}
	return
}

// GetEnums 根据枚举名字列表获得一批枚举
func GetEnums[T EnumDefinition](names ...string) (res []T, valid bool) {
	for _, n := range names {
		t, valid := ValueOf[T](n)
		if !valid {
			return nil, false
		}
		res = append(res, t)
	}
	valid = true
	return
}

// IsValidEnum 判断是否是合法的枚举
func IsValidEnum[T EnumDefinition](name string) (valid bool) {
	_, valid = ValueOf[T](name)
	return
}

func typeKey(t reflect.Type) string {
	return t.String()
}
