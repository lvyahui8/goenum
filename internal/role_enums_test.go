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
		r, valid := goenum.ValueOf[Role]("Reporter")
		require.True(t, valid)
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
	t.Run("textMarshal", func(t *testing.T) {
		var m = make(map[*Role]int)
		m[&Developer] = Developer.Ordinal()
		bytes, err := json.Marshal(m)
		require.Nil(t, err)
		require.Equal(t, "{\"Developer\":1}", string(bytes))
	})
	t.Run("String", func(t *testing.T) {
		t.Logf("formatted stirng %s", Developer)
	})
}

func TestHelpers(t *testing.T) {
	t.Run("NewEnum", func(t *testing.T) {
		defer func() {
			err := recover()
			require.NotNil(t, err)
			require.Equal(t, "Enum must be unique", err)
		}()
		_ = goenum.NewEnum[Role]("Owner")
	})
	t.Run("ValueOf", func(t *testing.T) {
		r, valid := goenum.ValueOf[Role]("Owner")
		require.True(t, valid)
		require.True(t, Owner.Equals(r))
		r, valid = goenum.ValueOf[Role]("Owner")
		require.True(t, valid)
		require.False(t, Developer.Equals(r))
	})
	t.Run("ValueOfIgnoreCase", func(t *testing.T) {
		r, valid := goenum.ValueOfIgnoreCase[Role]("oWnEr")
		require.True(t, valid)
		require.True(t, Owner.Equals(r))
		r, valid = goenum.ValueOfIgnoreCase[Role]("oWnEr")
		require.True(t, valid)
		require.False(t, Reporter.Equals(r))
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
	t.Run("GetEnums", func(t *testing.T) {
		res, valid := goenum.GetEnums[Role]("Owner", "Developer")
		require.True(t, valid)
		require.True(t, reflect.DeepEqual([]Role{Owner, Developer}, res))
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

// BenchmarkValueOf
//  go test -bench='BenchmarkValueOf'  -benchtime=5s -benchmem -count=3
func BenchmarkValueOf(b *testing.B) {
	n := 1000
	b.Run("ValueOf", func(b *testing.B) {
		for i := 0; i < n; i++ {
			_, _ = goenum.ValueOf[Permission]("AddTopic")
		}
	})
	b.Run("ValueOfIgnoreCase", func(b *testing.B) {
		for i := 0; i < n; i++ {
			_, _ = goenum.ValueOfIgnoreCase[Permission]("AddTopic")
		}
	})
}
