package goenum

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type EnumDefinition interface {
	fmt.Stringer
	json.Marshaler
	encoding.TextMarshaler
	// Name The name of the enumeration. The enumeration names of the same Type must be unique.
	Name() string
	// Equals Compare with another enumeration. Only return true if both Type and Name are the same
	Equals(other EnumDefinition) bool
	// Type String representation of enumeration type
	Type() string
	// Ordinal Get the ordinal of the enumeration, starting from zero and increasing in declared order.
	Ordinal() int
	// Compare -Compare with the ordinal value of another enumeration
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

// name2enumsMap Name to enumeration instance mapping.
// Value uses slice to store enumeration instances with different types but conflicting names
var name2enumsMap = make(map[string][]EnumDefinition)

// type2enumsMap The mapping from enumeration type to all enumerations,
// where key is the string representation of the enumeration type
var type2enumsMap = make(map[string][]EnumDefinition)

// typeIndexMap Store instance counters of different enumeration types for calculating Ordinal.
var typeIndexMap = make(map[string]int)

// NewEnum Create a new enumeration. If an enumeration instance with the same Type and Name already exists,
// the current method will throw a panic to prevent duplicate enumeration creation
func NewEnum[T EnumDefinition](name string, src ...T) T {
	if IsValidEnum[T](name) {
		panic("Enum must be unique")
	}
	var t T
	if len(src) > 0 {
		t = src[0]
	}
	v := reflect.ValueOf(t)
	tFullName := typeKey(v.Type())

	isPtr := v.Kind() == reflect.Ptr
	if isPtr {
		if v.IsNil() {
			v = reflect.New(v.Type().Elem())
		}
	} else {
		v = reflect.ValueOf(&t)
	}
	elem := reflect.Indirect(v)
	enumFiled := elem.FieldByName(reflect.TypeOf(Enum{}).Name())

	idx := typeIndexMap[tFullName]
	e := Enum{name: name, _type: tFullName, index: idx}
	if enumFiled.Kind() == reflect.Ptr {
		enumFiled.Set(reflect.ValueOf(&e))
	} else {
		enumFiled.Set(reflect.ValueOf(e))
	}
	typeIndexMap[tFullName] = idx + 1
	if isPtr {
		t = v.Interface().(T)
	} else {
		t = reflect.Indirect(v).Interface().(T)
	}
	type2enumsMap[tFullName] = append(type2enumsMap[tFullName], t)
	name2enumsMap[name] = append(name2enumsMap[name], t)
	return t
}

// ValueOf Find an enumeration instance based on the string, and return a zero value if not found
func ValueOf[T EnumDefinition](name string) (t T, valid bool) {
	enums := name2enumsMap[name]
	for _, e := range enums {
		if v, ok := e.(T); ok {
			return v, true
		}
	}
	return
}

// ValueOfIgnoreCase Ignoring case to obtain enumeration instances.
// Note: This method involves one reflection call,
// and its performance is slightly worse than the ValueOf method
func ValueOfIgnoreCase[T EnumDefinition](name string) (t T, valid bool) {
	values := Values[T]()
	for _, e := range values {
		if strings.EqualFold(e.Name(), name) {
			return e, true
		}
	}
	return
}

// Unmarshal Deserialize the enumeration instance from the JSON string,
// and return an error if not found
func Unmarshal[T EnumDefinition](data []byte) (t T, err error) {
	var name string
	err = json.Unmarshal(data, &name)
	if err != nil {
		return
	}
	t, valid := ValueOf[T](name)
	if !valid {
		return t, errors.New("enum not found")
	}
	return
}

// Values Return all enumeration instances. The returned slice are sorted by ordinal
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

// GetEnumMap Get all enumeration instances of the specified type.
// The key is the Name of the enumeration instance, and the value is the enumeration instance.
func GetEnumMap[T EnumDefinition]() map[string]T {
	values := Values[T]()
	res := make(map[string]T)
	for _, e := range values {
		res[e.Name()] = e
	}
	return res
}

// EnumNames Get the names of a batch of enumerations. If no instances are passed in,
// return the Name value of all enumeration instances of the type specified by the generic parameter.
func EnumNames[T EnumDefinition](enums ...T) (names []string) {
	if len(enums) == 0 {
		enums = Values[T]()
	}
	for _, e := range enums {
		names = append(names, e.Name())
	}
	return
}

// GetEnums Obtain a batch of enumeration instances based on the enumeration name list
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

// IsValidEnum Determine if the incoming string is a valid enumeration
func IsValidEnum[T EnumDefinition](name string) (valid bool) {
	_, valid = ValueOf[T](name)
	return
}

func typeKey(t reflect.Type) string {
	return t.String()
}
