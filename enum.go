package goenum

import (
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
	// Init 枚举初始化。使用方不应该直接调用这个方法。
	Init(args ...any) any
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

func (e Enum) Init(args ...any) any { return e }

func (e Enum) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Name())
}

// name2enumsMap name到枚举实例的映射，不同的枚举，name可能冲突，所以value是slice
var name2enumsMap = make(map[string][]EnumDefinition)

// type2enumsMap 枚举类型到所有枚举的映射，key为枚举类型的全路径名称
var type2enumsMap = make(map[string][]EnumDefinition)

// typeIndexMap 类型枚举索引
var typeIndexMap = make(map[string]int)

// NewEnum 新建枚举
func NewEnum[T EnumDefinition](name string, args ...any) T {
	var t T
	elem := reflect.ValueOf(&t).Elem()
	enumFiled := elem.FieldByName(reflect.TypeOf(Enum{}).Name())
	// 获取泛型具体类型名
	tFullName := typeKey(elem.Type())
	idx := typeIndexMap[tFullName]
	enumFiled.Set(reflect.ValueOf(Enum{name: name, _type: tFullName, index: idx}))
	typeIndexMap[tFullName] = idx + 1
	res := t.Init(args...)
	if updated, ok := res.(T); ok {
		t = updated
	}
	type2enumsMap[tFullName] = append(type2enumsMap[tFullName], t)
	name2enumsMap[name] = append(name2enumsMap[name], t)
	return t
}

// ValueOf 根据字符串获取枚举，如果找不到，则返回nil
func ValueOf[T EnumDefinition](name string) *T {
	enums := name2enumsMap[name]
	for _, e := range enums {
		if v, ok := e.(T); ok {
			return &v
		}
	}
	return nil
}

// ValueOfIgnoreCase 忽略大小写获取枚举, 涉及到一次反射调用，性能比ValueOf略差
func ValueOfIgnoreCase[T EnumDefinition](name string) *T {
	values := Values[T]()
	for _, e := range values {
		if strings.EqualFold(e.Name(), name) {
			return &e
		}
	}
	return nil
}

// Values 返回所有可用枚举
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

func GetEnumMap[T EnumDefinition]() map[string]T {
	values := Values[T]()
	res := make(map[string]T)
	for _, e := range values {
		res[e.Name()] = e
	}
	return res
}

// EnumNames 获取一批枚举的名称
func EnumNames(enums ...EnumDefinition) (names []string) {
	for _, e := range enums {
		names = append(names, e.Name())
	}
	return
}

func IsValidEnum[T EnumDefinition](name string) bool {
	return ValueOf[T](name) != nil
}

func typeKey(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}
