package goenum

import (
	"reflect"
)

type EnumDefinition interface {
	// Name 枚举名称，同一类型枚举应该唯一
	Name() string
	// Init 枚举初始化
	Init(args ...any) any
	// Equals 枚举对比
	Equals(other EnumDefinition) bool
	// Type 实际的枚举类型
	Type() string
}

type Enum struct {
	name  string
	_type string
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

func (e Enum) Type() string {
	return e._type
}

func (e Enum) Init(args ...any) any { return e }

// name2enumsMap name到枚举实例的映射，不同的枚举，name可能冲突，所以value是slice
var name2enumsMap = make(map[string][]EnumDefinition)

// type2enumsMap 枚举类型到所有枚举的映射，key为枚举类型的全路径名称
var type2enumsMap = make(map[string][]EnumDefinition)

func NewEnum[T EnumDefinition](name string, args ...interface{}) T {
	var t T
	elem := reflect.ValueOf(&t).Elem()
	enumFiled := elem.FieldByName(reflect.TypeOf(Enum{}).Name())
	// 获取泛型具体类型名
	tFullName := typeKey(elem.Type())
	enumFiled.Set(reflect.ValueOf(Enum{name: name, _type: tFullName}))
	type2enumsMap[tFullName] = append(type2enumsMap[tFullName], t)
	name2enumsMap[name] = append(name2enumsMap[name], t)
	return t
}

func ValueOf[T EnumDefinition](name string) *T {
	enums := name2enumsMap[name]
	for _, e := range enums {
		if v, ok := e.(T); ok {
			return &v
		}
	}
	return nil
}

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

func EnumNames(enums ...EnumDefinition) (names []string) {
	for _, e := range enums {
		names = append(names, e.Name())
	}
	return
}

func typeKey(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}
