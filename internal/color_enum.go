package internal

import "github.com/lvyahui8/goenum"

type ColorEnum struct {
	*goenum.Enum
}

var (
	Red    = goenum.NewEnum[*ColorEnum]("Red")
	Yellow = goenum.NewEnum[*ColorEnum]("Yellow")
)
