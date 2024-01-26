package pkga

import (
	"github.com/lvyahui8/goenum"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestStatus(t *testing.T) {
	require.NotNil(t, Created)
	require.Equal(t, "Created", Created.Name())
	e, valid := goenum.ValueOf[Status]("Created")
	require.True(t, valid)
	require.Equal(t, Created, e)
	e, valid = goenum.ValueOf[Status]("Created")
	require.True(t, valid)
	require.True(t, Created == e)
	require.False(t, Created.Equals(Failed))
	require.True(t, Created.Equals(Created))
	e, valid = goenum.ValueOf[Status]("Created")
	require.True(t, valid)
	require.True(t, Created.Equals(e))
	require.True(t, reflect.DeepEqual([]Status{Created, Pending, Success, Failed}, goenum.Values[Status]()))
}
