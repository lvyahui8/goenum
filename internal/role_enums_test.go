package internal

import (
	"encoding/json"
	"github.com/lvyahui8/goenum"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type Member struct {
	RoleName Role
}

func TestRoleBasic(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		require.NotNil(t, Owner.perms)
	})
	t.Run("OrdinalAndCompare", func(t *testing.T) {
		// ordinal 方法
		require.True(t, Reporter.Ordinal() == 0)
		require.True(t, Owner.Compare(Reporter) > 0)
	})
	t.Run("Equals", func(t *testing.T) {
		// valueOf测试
		r := *goenum.ValueOf[Role]("Reporter")
		require.True(t, Reporter.Equals(r))
		require.True(t, r.Equals(r))
		require.False(t, Developer.Equals(Reporter))
		require.False(t, Developer.Equals(r))
	})
	t.Run("Name", func(t *testing.T) {
		require.Equal(t, "Owner", Owner.Name())
	})
	t.Run("jsonMarshal", func(t *testing.T) {
		bytes, err := json.Marshal(Developer)
		require.Nil(t, err)
		require.Equal(t, "\"Developer\"", string(bytes))
		// 测试枚举序列化功能
		member := Member{RoleName: Reporter}
		bytes, err = json.Marshal(member)
		require.Nil(t, err)
		require.Equal(t, "{\"RoleName\":\"Reporter\"}", string(bytes))
	})
}

func TestHelpers(t *testing.T) {
	t.Run("ValueOf", func(t *testing.T) {
		require.True(t, Owner.Equals(*goenum.ValueOf[Role]("Owner")))
		require.False(t, Developer.Equals(*goenum.ValueOf[Role]("Owner")))
	})
	t.Run("ValueOfIgnoreCase", func(t *testing.T) {
		require.True(t, Owner.Equals(*goenum.ValueOfIgnoreCase[Role]("oWnEr")))
		require.False(t, Reporter.Equals(*goenum.ValueOfIgnoreCase[Role]("oWnEr")))
	})
	t.Run("Values", func(t *testing.T) {
		require.True(t, reflect.DeepEqual([]Role{Reporter, Developer, Owner}, goenum.Values[Role]()))
	})
	t.Run("GetEnumMap", func(t *testing.T) {
		enumMap := goenum.GetEnumMap[Role]()
		require.True(t, len(enumMap) == 3)
		role, exist := enumMap["Owner"]
		require.True(t, exist)
		require.True(t, role.Equals(Owner))
	})
	t.Run("EnumNames", func(t *testing.T) {
		require.True(t, reflect.DeepEqual([]string{"Owner", "Developer"}, goenum.EnumNames(Owner, Developer)))
	})
	t.Run("IsValidEnum", func(t *testing.T) {
		require.True(t, goenum.IsValidEnum[Role]("Owner"))
		require.False(t, goenum.IsValidEnum[Role]("Test"))
	})
}

func TestRole_HasPerm(t *testing.T) {
	require.True(t, Owner.HasPerm(DeleteMergeRequest))
	require.True(t, Owner.HasPerm(ViewMergeRequest))
	require.False(t, Reporter.HasPerm(DeleteMergeRequest))
	require.False(t, Developer.HasPerm(DeleteMergeRequest))
}

func TestModule_BasePath(t *testing.T) {
	require.True(t, len(Issues.basePath) > 0)
	require.True(t, len(MergeRequests.basePath) > 0)
}

func BenchmarkValueOf(b *testing.B) {
	b.Run("ValueOf", func(b *testing.B) {
		_ = goenum.ValueOf[Permission]("AddTopic")
	})
	b.Run("ValueOfIgnoreCase", func(b *testing.B) {
		_ = goenum.ValueOfIgnoreCase[Permission]("AddTopic")
	})
}
