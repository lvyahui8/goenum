package internal

import "github.com/lvyahui8/goenum"

type State struct {
	goenum.Enum
	final bool
}

func (s State) IsFinal() bool {
	return s.final
}

// TradeState 不能使用type TradeState State定义类型，这样TradeState无法访问IsFinal方法
type TradeState struct {
	State
}

var (
	TradeCreated   = goenum.NewEnum[TradeState]("Created")
	TradeFailed    = goenum.NewEnum[TradeState]("Failed", TradeState{State: State{final: true}})
	TradePaid      = goenum.NewEnum[TradeState]("Paid")
	TradeShipped   = goenum.NewEnum[TradeState]("Shipped")
	TradeDelivered = goenum.NewEnum[TradeState]("Delivered", TradeState{State: State{final: true}})
)

type ReverseState struct {
	State
}

var (
	ReverseCreated  = goenum.NewEnum[ReverseState]("Created")
	ReverseFailed   = goenum.NewEnum[ReverseState]("Failed", ReverseState{State: State{final: true}})
	ReverseRefunded = goenum.NewEnum[ReverseState]("Refunded", ReverseState{State: State{final: true}})
)
