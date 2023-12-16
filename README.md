## Go 通用枚举实现

[![License](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvyahui8/goenum)](https://goreportcard.com/report/github.com/lvyahui8/goenum)
[![codecov](https://codecov.io/gh/lvyahui8/goenum/graph/badge.svg?token=YBV3TH2HQU)](https://codecov.io/gh/lvyahui8/goenum)

### 怎么定义和使用枚举？

**只需往枚举struct内嵌goenum.Enum, 即可定义一个枚举类型，并获得开箱即用的一组方法。**

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
goenum.ValueOf[State]("Created") // struct instance: Created
Created.Equals(*goenum.ValueOf[State]("Created")) // true
goenum.Values[State]() // equals []State{Created,Running,Success}
```

### 实例方法

```go
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
```

### 工具方法

- ValueOf 根据字符串获取枚举，如果找不到，则返回nil

- ValueOfIgnoreCase 忽略大小写获取枚举, 涉及到一次反射调用，性能比ValueOf略差

- Values 返回所有枚举

- GetEnumMap 获取所有枚举，以name->enum map的形式返回

- EnumNames  获取一批枚举的名称

- IsValidEnum 判断是否是合法的枚举

```go
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
```

### 进阶：复杂枚举初始化

枚举struct 实现Init方法即可，NewEnum方法中的args参数，会完整透传给Init方法，注意，**Init方法需要将receiver返回以确保初始化生效**。

完整例子请看 [gitlab_role_perms](internal/role_enums.go)

```go
type Module struct {
    goenum.Enum
    perms    []Permission
    basePath string
}

func (m Module) Init(args ...any) any {
    m.perms = args[0].([]Permission)
    m.basePath = args[1].(string)
    return m
}

func (m Module) GetPerms() []Permission {
    return m.perms
}

func (m Module) BasePath() string {
    return m.basePath
}


// 定义模块
var (
    Issues        = goenum.NewEnum[Module]("Issues", []Permission{AddLabels, AddTopic}, "/issues/")
    MergeRequests = goenum.NewEnum[Module]("MergeRequests", []Permission{ViewMergeRequest, ApproveMergeRequest, DeleteMergeRequest}, "/merge/")
)
```

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
