package internal

import "github.com/lvyahui8/goenum"

type ColorEnum struct {
	*goenum.Enum // 支持组合指针，NewEnum传入相应的也需要指针类型
}

var (
	Red    = goenum.NewEnum[*ColorEnum]("Red")
	Yellow = goenum.NewEnum[*ColorEnum]("Yellow")
)
