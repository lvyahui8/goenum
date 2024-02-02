package internal

import "github.com/lvyahui8/goenum"

type Code struct {
	goenum.Enum
	code int
	desc string
}

func (c Code) Code() int {
	return c.code
}

func (c Code) Desc() string {
	return c.desc
}

// ErrorCode 这里加等号与不加等号含义是不同的，
//加等号，ErrorCode = Code， 表示完全是相同的类型，reflect.Type(Success) = Code。
//而不加等号，则是创建了一个新的类型， reflect.Type(Success) = ErrorCode
//所以，如果加等号，就不能出现name相同的枚举，否则会破坏同类型枚举不能出现同名枚举的原则。类似的示例可以参考 state.go
type ErrorCode = Code

var (
	Success      = goenum.NewEnum[ErrorCode]("Success", ErrorCode{code: 0, desc: "成功"})
	Failed       = goenum.NewEnum[ErrorCode]("Failed", ErrorCode{code: -1, desc: "未知异常"})
	NetworkError = goenum.NewEnum[ErrorCode]("NetworkError", ErrorCode{code: 500, desc: "网络错误"})
	EncodeError  = goenum.NewEnum[ErrorCode]("EncodeError", ErrorCode{code: 600, desc: "编码错误"})
)

type BizCode = Code

var (
	Payment  = goenum.NewEnum[BizCode]("Member", BizCode{code: 2, desc: "支付服务"})
	Trade    = goenum.NewEnum[BizCode]("Trade", BizCode{code: 2, desc: "交易服务"})
	Delivery = goenum.NewEnum[BizCode]("Delivery", BizCode{code: 2, desc: "履约服务"})
)
