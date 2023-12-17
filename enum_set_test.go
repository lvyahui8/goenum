// goenum
// 测试复杂枚举使用set，复杂枚举例子
// - com.sun.tools.doclint.HtmlTag
// - org.apache.logging.log4j.spi.StandardLevel
// - javax.lang.model.SourceVersion
package goenum

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

type Statement struct {
	Enum
}

var (
	Decl       = NewEnum[Statement]("Decl")
	Empty      = NewEnum[Statement]("Empty")
	Labeled    = NewEnum[Statement]("Labeled")
	Expr       = NewEnum[Statement]("Expr")
	Send       = NewEnum[Statement]("Send")
	IncDec     = NewEnum[Statement]("IncDec")
	Assign     = NewEnum[Statement]("Assign")
	Go         = NewEnum[Statement]("Go")
	Defer      = NewEnum[Statement]("Defer")
	Return     = NewEnum[Statement]("Return")
	Branch     = NewEnum[Statement]("Branch")
	Block      = NewEnum[Statement]("Block")
	If         = NewEnum[Statement]("If")
	Case       = NewEnum[Statement]("Case")
	Switch     = NewEnum[Statement]("Switch")
	TypeSwitch = NewEnum[Statement]("TypeSwitch")
	Comm       = NewEnum[Statement]("Comm")
	Select     = NewEnum[Statement]("Select")
	For        = NewEnum[Statement]("For")
	Range      = NewEnum[Statement]("Range")
)

func TestEnumSet_Basic(t *testing.T) {
	stmtSet := NewUnsafeEnumSet[Statement]()
	require.True(t, stmtSet.IsEmpty())
	require.True(t, stmtSet.Add(Decl))
	require.False(t, stmtSet.IsEmpty())
	require.True(t, stmtSet.Contains(Decl))
	require.False(t, stmtSet.Contains(Select))
	require.True(t, stmtSet.Len() == 1)
	require.True(t, stmtSet.Add(Select))
	require.True(t, stmtSet.Len() == 2)
	require.True(t, stmtSet.Contains(Select))
	require.False(t, stmtSet.Contains(For))
	// string、json
	require.Equal(t, "[Decl,Select]", stmtSet.String())
	bytes, err := json.Marshal(stmtSet)
	require.Nil(t, err)
	require.Equal(t, "[\"Decl\",\"Select\"]", string(bytes))
	// Equals
	sameStmtSet := NewUnsafeEnumSet[Statement]()
	sameStmtSet.Add(Decl)
	require.False(t, stmtSet.Equals(sameStmtSet))
	sameStmtSet.Add(Select)
	require.True(t, stmtSet.Equals(sameStmtSet))
	// Clone
	copiedSet := stmtSet.Clone()
	require.True(t, stmtSet.Equals(copiedSet))
	copiedSet.Add(If)
	require.True(t, copiedSet.Contains(If))
	require.False(t, stmtSet.Contains(If))
	require.False(t, stmtSet.Equals(copiedSet))
	require.True(t, copiedSet.Len()-stmtSet.Len() == 1)
	// addRange
	require.True(t, stmtSet.AddRange(Comm, Range) == 3) // Comm -> Range一共4个，但是select已经添加过，所以实际添加3个
	require.True(t, stmtSet.Len() == 5)                 // 2+3
	// 验证contains
	subStmt := NewUnsafeEnumSet[Statement]()
	require.True(t, subStmt.Add(Decl))
	require.True(t, subStmt.AddRange(Select, Range) == 3)
	require.True(t, stmtSet.ContainsAll(subStmt))
	otherStmt := NewUnsafeEnumSet[Statement]()
	otherStmt.Add(Decl)
	otherStmt.Add(Switch)
	require.False(t, stmtSet.ContainsAll(otherStmt))
	// 删除
	require.True(t, stmtSet.Remove(Decl))
	require.False(t, stmtSet.Remove(Branch)) // 删除不存在的枚举，应该返回false
	require.False(t, stmtSet.Contains(Decl))
	require.True(t, stmtSet.AddRange(Expr, Assign) == 4)
	require.True(t, stmtSet.RemoveRange(IncDec, Go) == 2) // 实际只能删除2个
	stmtSet.Clear()
	require.False(t, stmtSet.Contains(Select))
	require.True(t, stmtSet.Len() == 0)
}

// BenchmarkUnsafeEnumSet_Contains 测试下来性能反而更差一些，原因应该是枚举数量太少，随着枚举数量的增加，map性能肯定会变差的，而EnumSet可以维持性能。
// java EnumSet，按照是否超过64个，划分了2个实现， 估计也是这个原因
func BenchmarkUnsafeEnumSet_Contains(b *testing.B) {
	stmtSet := NewUnsafeEnumSet[Statement]()
	stmtSet.Add(TypeSwitch)
	enumMap := GetEnumMap[Statement]()
	b.Run("map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = enumMap[TypeSwitch.Name()]
		}
	})
	b.Run("enumSet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = stmtSet.Contains(TypeSwitch)
		}
	})
}
