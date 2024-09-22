## Go enumeration template implementation

[![License](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvyahui8/goenum)](https://goreportcard.com/report/github.com/lvyahui8/goenum)
[![codecov](https://codecov.io/gh/lvyahui8/goenum/graph/badge.svg?token=YBV3TH2HQU)](https://codecov.io/gh/lvyahui8/goenum)

[中文文档](./README_CN.md)

### How to define and use enumerations?

**Simply embed goenum.Enum into a struct to define an enumeration type and automatically obtain a set of instance methods and utility functions.**

```shell
go get github.com/lvyahui8/goenum
```

```go
import "github.com/lvyahui8/goenum"

// Declaring an enumeration type
type State struct {
    goenum.Enum
}

// Defines a set of enumeration instances
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

More Examples [examples](internal)

### Instance Methods

```go
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
```

### Utility functions

- ValueOf: Find an enumeration instance based on the string, and return a zero value if not found
- ValueOfIgnoreCase: Ignoring case to obtain enumeration instances.
- Values: Return all enumeration instances. The returned slice are sorted by ordinal
- Size: Number of instances of specified enumeration type
- GetEnumMap: Get all enumeration instances of the specified type.
- EnumNames: Get the names of a batch of enumerations.
- GetEnums: Obtain a batch of enumeration instances based on the enumeration name list
- IsValidEnum: Determine if the incoming string is a valid enumeration 

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

### More features

#### Complex enumeration initialization

The second parameter in the NewEnum function, if a Src is passed in, NewEnum will use the passed object to construct an enumeration instance.

example code [gitlab_role_perms](internal/role_enums.go)

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

#### JSON serialization and deserialization support for enumerating instances

JSON serialization is already supported by default. Restricted by the implementation of the go JSON library, enumeration classes are required to implement JSON Unmarshaler interface to realize deserialization. Call Unmarshal tool function in the interface to obtain enumeration instances.

```go
type Member struct {
	Roles []Role
}

type Role struct {
	goenum.Enum
	perms []Permission
}

func (r *Role) UnmarshalJSON(data []byte) (err error) {
	role, err := goenum.Unmarshal[Role](data)
	if err == nil {
		*r = role
	}
	return
}

t.Run("jsonMarshal", func(t *testing.T) {
    bytes, err := json.Marshal(Developer)
    require.Nil(t, err)
    require.Equal(t, "\"Developer\"", string(bytes))
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

Example code [enum_set_test](enum_set_test.go)

### ValueOf Performance

Don't worry about any performance issues, reflection calls are mostly only used in NewEnum methods, and other methods will try to avoid reflection calls as much as possible.

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
