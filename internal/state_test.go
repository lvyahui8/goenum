package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestState_IsFinal(t *testing.T) {
	require.False(t, TradePaid.IsFinal())
	require.True(t, TradeDelivered.IsFinal())

	require.False(t, ReverseCreated.IsFinal())
	require.True(t, ReverseRefunded.IsFinal())

	require.False(t, ReverseCreated.Equals(TradeCreated))        // false， 类型不同
	require.Equal(t, ReverseCreated.Name(), TradeCreated.Name()) // true，实际的枚举值一样的，都是Created
}
