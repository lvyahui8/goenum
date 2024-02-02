package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCode_Code(t *testing.T) {
	require.Equal(t, "Success", Success.Name())
	require.Equal(t, 0, Success.Code())
}
