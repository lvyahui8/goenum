package pkgb

import "github.com/lvyahui8/goenum"

type State struct {
	goenum.Enum
}

var (
	Created = goenum.NewEnum[State]("Created")
	Running = goenum.NewEnum[State]("Running")
	Success = goenum.NewEnum[State]("Success")
)
