## Go 通用枚举实现

[![License](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvyahui8/goenum)](https://goreportcard.com/report/github.com/lvyahui8/goenum)
[![codecov](https://codecov.io/gh/lvyahui8/goenum/graph/badge.svg?token=YBV3TH2HQU)](https://codecov.io/gh/lvyahui8/goenum)

### 怎么定义和使用枚举？

**只需往枚举struct内嵌goenum.Enum, 即可定义一个枚举类型，并获得开箱即用的一组方法。**

```shell
go get github.com/lvyahui8/goenum
```

```go
import "github.com/lvyahui8/goenum"

// 声明枚举类型
type State struct {
    goenum.Enum
}

// 定义枚举
var (
    Created = goenum.NewEnum[State]("Created")
    Running = goenum.NewEnum[State]("Running")
    Success = goenum.NewEnum[State]("Success")
)

// Usage
Created.Name() // string "Created"
Created.Ordinal() // int 0
Created.Compare(Running) //  < 0 : Created < Running
Created.String() // Created
json.Marshal(Created) // \"Created\"
IsValidEnum[State]("Created") // true
s,valid := goenum.ValueOf[State]("Created") // s: Created(struct instance) ,valid = true
Created.Equals(s) // true
s,valid := goenum.ValueOf[State]("cReaTed") //  s: Created(struct instance) ,valid = true
Created.Equals(s) // true
goenum.Values[State]() // equals []State{Created,Running,Success}
Size[State]() // = 3
```

### 实例方法

```go
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
```

### 工具方法

- ValueOf 根据字符串获取枚举，如果找不到，则返回nil
- ValueOfIgnoreCase 忽略大小写获取枚举, 涉及到一次反射调用，性能比ValueOf略差
- Values 返回所有枚举
- GetEnumMap 获取所有枚举，以name->enum map的形式返回
- EnumNames  获取一批枚举的名称
- GetEnums 根据枚举名字列表获得一批枚举
- IsValidEnum 判断是否是合法的枚举

```go
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
		_, valid = goenum.GetEnums[Role]("a", "b")
		require.False(t, valid)
	})
	t.Run("IsValidEnum", func(t *testing.T) {
		require.True(t, goenum.IsValidEnum[Role]("Owner"))
		require.False(t, goenum.IsValidEnum[Role]("Test"))
	})
}
```

### 更多特性支持

#### 复杂枚举初始化

NewEnum方法中的第二个参数，如果有传入一个src，NewEnum将使用传入的对象构造枚举对象。

完整例子请看 [gitlab_role_perms](internal/role_enums.go)

```go
type Module struct {
	goenum.Enum
	perms    []Permission
	basePath string
}

func (m Module) GetPerms() []Permission {
	return m.perms
}

func (m Module) BasePath() string {
	return m.basePath
}

// 定义模块
var (
	Issues        = goenum.NewEnum[Module]("Issues", Module{perms: []Permission{AddLabels, AddTopic}, basePath: "/issues/"})
	MergeRequests = goenum.NewEnum[Module]("MergeRequests", Module{perms: []Permission{ViewMergeRequest, ApproveMergeRequest, DeleteMergeRequest}, basePath: "/merge/"})
)
```

#### 枚举json序列化和反序列化支持

默认已支持了json序列化，受限于go json库的反序列化机制限制，需要枚举类自行实现json.Unmarshaler接口，在接口中调用工具方法即可。

```go
type Member struct {
	Roles []Role
}

type Role struct {
	goenum.Enum
	perms []Permission
}
// UnmarshalJSON 枚举类需要自行实现json.Unmarshaler接口
func (r *Role) UnmarshalJSON(data []byte) error {
	role, err := goenum.Unmarshal[Role](data)
	if err != nil {
		return err
	}
	*r = role
	return nil
}

t.Run("jsonMarshal", func(t *testing.T) {
    bytes, err := json.Marshal(Developer)
    require.Nil(t, err)
    require.Equal(t, "\"Developer\"", string(bytes))
    // 测试枚举序列化功能
    member := Member{Roles: []Role{Reporter, Owner}}
    bytes, err = json.Marshal(member)
    require.Nil(t, err)
    require.Equal(t, "{\"Roles\":[\"Reporter\",\"Owner\"]}", string(bytes))
})
t.Run("jsonUnmarshal", func(t *testing.T) {
    newMember := Member{}
    err := json.Unmarshal([]byte("{\"Roles\":[\"Reporter\",\"Owner\"]}"), &newMember)
    require.Nil(t, err)
    require.True(t, reflect.DeepEqual([]Role{Reporter, Owner}, newMember.Roles))
})
```

#### EnumSet

api声明
```go
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
```

EnumSet usage

```go
stmtSet := NewUnsafeEnumSet[Statement]()

stmtSet.Add(Decl)
stmtSet.IsEmpty()
stmtSet.Contains(Decl)
stmtSet.Len() == 1
stmtSet.AddRange(Comm, Range)
stmtSet.Remove(Decl)
```

完整例子请看 [enum_set_test](enum_set_test.go)

### ValueOf性能测试

不用担心任何性能问题，反射调用基本集中在NewEnum方法中，其他方法尽量避免反射调用。

```text
goos: linux
goarch: arm64
pkg: github.com/lvyahui8/goenum/internal
BenchmarkValueOf/ValueOf-4             1000000000             0.0002492 ns/op           0 B/op           0 allocs/op
BenchmarkValueOf/ValueOf-4             1000000000             0.0002966 ns/op           0 B/op           0 allocs/op
BenchmarkValueOf/ValueOf-4             1000000000             0.0002713 ns/op           0 B/op           0 allocs/op
BenchmarkValueOf/ValueOfIgnoreCase-4             1000000000             0.002228 ns/op           0 B/op           0 allocs/op
BenchmarkValueOf/ValueOfIgnoreCase-4             1000000000             0.002611 ns/op           0 B/op           0 allocs/op
BenchmarkValueOf/ValueOfIgnoreCase-4             1000000000             0.002606 ns/op           0 B/op           0 allocs/op
PASS
ok      github.com/lvyahui8/goenum/internal    0.097s
```
