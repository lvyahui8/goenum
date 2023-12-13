package pkgb

import (
	"github.com/lvyahui8/goenum"
	"github.com/lvyahui8/goenum/pkga"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestState(t *testing.T) {
	require.True(t, reflect.DeepEqual([]State{Created, Running, Success}, goenum.Values[State]()))
}

func TestNameConflict(t *testing.T) {
	require.False(t, pkga.Created.Equals(Created))
	require.True(t, pkga.Created.Name() == Created.Name())
	require.False(t, pkga.Created.Type() == Created.Type())
}
