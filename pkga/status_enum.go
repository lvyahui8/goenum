package pkga

import "github.com/lvyahui8/goenum"

type Status struct {
	goenum.Enum
}

var (
	Created = goenum.NewEnum[Status]("Created")
	Pending = goenum.NewEnum[Status]("Pending")
	Success = goenum.NewEnum[Status]("Success")
	Failed  = goenum.NewEnum[Status]("Failed")
)
