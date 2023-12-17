package goenum

import (
	"encoding/json"
	"fmt"
	"strings"
)

// EnumSet 枚举set，一般在枚举非常多时使用
type EnumSet[E EnumDefinition] interface {
	// Stringer 支持标准输出及格式化
	fmt.Stringer
	// Marshaler 支持json序列化
	json.Marshaler
	// Add 往set添加元素，添加成功则返回true，如果已经存在则返回false
	Add(e E) bool
	// AddRange 按照枚举的序数，连续添加一段枚举，返回实际添加的数量（排除已经存在的）
	AddRange(begin, end E) int
	// Remove 删除元素，删除成功则返回true，如果元素原本不存在则返回false
	Remove(e E) bool
	// RemoveRange  按照枚举的序数，连续删除一段枚举，返回实际删除的数量（排除原本不存在的）
	RemoveRange(begin, end E) int
	// IsEmpty set是否为空
	IsEmpty() bool
	// Clear 清理set
	Clear()
	// Len set内当前的枚举数量
	Len() int
	// Contains 是否包含指定的枚举，只要有1个不存在则返回false
	Contains(enums ...E) bool
	// ContainsAll  判断是否包含另外一个enumSet（子集关系）
	ContainsAll(set EnumSet[E]) bool
	// Equals 判断两个EnumSet是否相同
	Equals(set EnumSet[E]) bool
	// Each set迭代方法, f方法如果返回false，则中止迭代
	Each(f func(e E) bool)
	// Names 返回set中已有枚举的Name表示
	Names() []string
	// Clone 深拷贝一份set
	Clone() EnumSet[E]
}

func NewUnsafeEnumSet[E EnumDefinition]() EnumSet[E] {
	enumSize := Size[E]()
	return &UnsafeEnumSet[E]{
		enumSize: enumSize,
		elements: make([]uint64, (enumSize+63)>>6 /*除以64并向上取整*/),
	}
}

type UnsafeEnumSet[E EnumDefinition] struct {
	enumSize int
	cap      int
	// elements 低位放低位枚举
	// elements[0] high <- low ordinal(63) <- ordinal(0)
	// elements[1] high <- low ordinal(127) <- ordinal(64)
	// ...
	elements []uint64
	// 已存放枚举数量
	len int
}

func (set *UnsafeEnumSet[E]) String() string {
	return "[" + strings.Join(set.Names(), ",") + "]"
}

func (set *UnsafeEnumSet[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Names())
}

func (set *UnsafeEnumSet[E]) Add(e E) bool {
	return set.addIdx(e.Ordinal())
}

func (set *UnsafeEnumSet[E]) addIdx(ordinal int) bool {
	// set.elements[i]
	i := ordinal >> 6
	old := set.elements[i]
	// 等价于 elements[i] |= (uint64(1) << (ordinal % 64))
	set.elements[i] |= uint64(1) << ordinal
	added := old != set.elements[i]
	if added {
		set.len++
	}
	return added
}

func (set *UnsafeEnumSet[E]) AddRange(begin, end E) int {
	cnt := 0
	for i := begin.Ordinal(); i <= end.Ordinal(); i++ {
		if set.addIdx(i) {
			cnt++
		}
	}
	return cnt
}

func (set *UnsafeEnumSet[E]) removeIdx(ordinal int) bool {
	i := ordinal >> 6
	old := set.elements[i]
	set.elements[i] &= ^(uint64(1) << ordinal)
	deleted := old != set.elements[i]
	if deleted {
		set.len--
	}
	return deleted
}

func (set *UnsafeEnumSet[E]) Remove(e E) bool {
	return set.removeIdx(e.Ordinal())
}

func (set *UnsafeEnumSet[E]) RemoveRange(begin, end E) int {
	cnt := 0
	for i := begin.Ordinal(); i <= end.Ordinal(); i++ {
		if set.removeIdx(i) {
			cnt++
		}
	}
	return cnt
}

func (set *UnsafeEnumSet[E]) Len() int {
	return set.len
}

func (set *UnsafeEnumSet[E]) IsEmpty() bool {
	return set.Len() == 0
}

func (set *UnsafeEnumSet[E]) Clear() {
	for i := 0; i < len(set.elements); i++ {
		set.elements[i] = 0
	}
	set.len = 0
}

func (set *UnsafeEnumSet[E]) Contains(enums ...E) bool {
	for _, e := range enums {
		i := e.Ordinal() >> 6
		flag := uint64(1) << e.Ordinal()
		if set.elements[i]&flag != flag {
			return false
		}
	}
	// 如果len(enums)==0， 应该为true还是false？
	return true
}

func (set *UnsafeEnumSet[E]) ContainsAll(enumSet EnumSet[E]) bool {
	if es, ok := enumSet.(*UnsafeEnumSet[E]); ok {
		for i := 0; i < len(es.elements); i++ {
			if es.elements[i]&set.elements[i] != es.elements[i] {
				return false
			}
		}
		return true
	}
	notFound := false
	enumSet.Each(func(e E) bool {
		if !set.Contains(e) {
			notFound = true
			return false
		}
		return true
	})
	return !notFound
}

func (set *UnsafeEnumSet[E]) Equals(enumSet EnumSet[E]) bool {
	return set.ContainsAll(enumSet) && enumSet.ContainsAll(set)
}

func (set *UnsafeEnumSet[E]) Each(f func(e E) bool) {
	allEnums := Values[E]()
	for _, e := range allEnums {
		if set.Contains(e) {
			if !f(e) {
				break
			}
		}
	}
}

func (set *UnsafeEnumSet[E]) Names() []string {
	var list []string
	set.Each(func(e E) bool {
		list = append(list, e.Name())
		return true
	})
	return list
}

func (set *UnsafeEnumSet[E]) Clone() EnumSet[E] {
	res := &UnsafeEnumSet[E]{
		enumSize: set.enumSize,
		len:      set.len,
		elements: make([]uint64, len(set.elements)),
	}
	copy(res.elements, set.elements)
	return res
}
