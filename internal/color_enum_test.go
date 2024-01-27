package internal

import (
	"github.com/lvyahui8/goenum"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPtrEnum(t *testing.T) {
	t.Run("ValueOf", func(t *testing.T) {
		c, valid := goenum.ValueOf[*ColorEnum]("Red")
		require.True(t, valid)
		require.True(t, c.Equals(Red))
		_, valid = goenum.ValueOf[ColorEnum]("Red")
		require.False(t, valid)
	})
}
