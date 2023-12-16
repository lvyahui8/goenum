package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Api interface {
	GetName() string
	getAge() int
}

type B struct {
}

func (b B) GetName() string {
	return "b"
}

func (b B) getAge() int {
	return 1
}

type A struct {
	B
}

func (a A) GetName() string {
	return "a"
}

func (a A) getAge() int {
	return -1
}

func Call[T Api]() (string, int) {
	var t T
	return t.GetName(), t.getAge()
}

// TestOver 结论：在没有跨包的情况下，是可以调用到子类的实现方法的，跨包之后则不行，可能不能破坏go的可见性
func TestOver(t *testing.T) {
	var api Api
	api = A{B{}}
	require.Equal(t, "a", api.GetName())
	require.Equal(t, -1, api.getAge())
	name, age := Call[A]()
	require.Equal(t, "a", name)
	require.Equal(t, -1, age)
}
