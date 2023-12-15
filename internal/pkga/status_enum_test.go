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
	require.Equal(t, Created, *goenum.ValueOf[Status]("Created"))
	require.True(t, Created == *goenum.ValueOf[Status]("Created"))
	require.False(t, Created.Equals(Failed))
	require.True(t, Created.Equals(Created))
	require.True(t, Created.Equals(*goenum.ValueOf[Status]("Created")))
	require.True(t, reflect.DeepEqual([]Status{Created, Pending, Success, Failed}, goenum.Values[Status]()))
}
