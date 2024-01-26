package internal

import "github.com/lvyahui8/goenum"

type ColorEnum struct {
	*goenum.Enum
}

var (
	red    = goenum.NewEnum[*ColorEnum]("red")
	yellow = goenum.NewEnum[*ColorEnum]("yellow")
)
