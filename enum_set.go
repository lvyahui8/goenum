package goenum

import (
	"encoding/json"
	"fmt"
	"strings"
)

// EnumSet Enumerate sets. Usually used for high-performance lookup when
// there are many instances of a certain type of enumeration.
type EnumSet[E EnumDefinition] interface {
	fmt.Stringer
	json.Marshaler
	// Add -Add an element to Set. If successful, return true. If it already exists, return false
	Add(e E) bool
	// AddRange According to the ordinal of the enumeration,
	// add a continuous section of the enumeration and
	// return the actual number added (excluding those that already exist)
	AddRange(begin, end E) int
	// Remove Delete element. If the deletion is successful, return true.
	// If the element does not exist, return false
	Remove(e E) bool
	// RemoveRange  According to the ordinal of the enumeration,
	// continuously delete a segment of the enumeration and
	// return the actual number of deletions  (excluding those that non-exist)
	RemoveRange(begin, end E) int
	// IsEmpty Is Set empty
	IsEmpty() bool
	// Clear -Clear set
	Clear()
	// Len The current number of enumeration instances within the Set
	Len() int
	// Contains Does it contain the specified enumeration?
	// Returns false if there is only one that does not exist in the Set
	Contains(enums ...E) bool
	// ContainsAll  Determine if it contains another enumSet (subset relationship)
	ContainsAll(set EnumSet[E]) bool
	// Equals Determine if two EnumSets are the same
	Equals(set EnumSet[E]) bool
	// Each Set iteration method, if f method returns false, abort iteration
	Each(f func(e E) bool)
	// Names Returns the Name representation of an existing enumeration instance in the set
	Names() []string
	// Clone Deep copy to obtain a new set
	Clone() EnumSet[E]
}

func NewUnsafeEnumSet[E EnumDefinition]() *UnsafeEnumSet[E] {
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
		if set.elements[e.Ordinal()>>6]&(uint64(1)<<e.Ordinal()) == 0 {
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
