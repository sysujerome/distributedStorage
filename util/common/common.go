package common

type Status struct {
	Ok     string
	Failed string
	Stored string
}

type Error struct {
	NotFound   string
	NotDefined string
}

func (sta *Status) Init() {
	sta.Ok = "OK"
	sta.Failed = "failed"
	sta.Stored = "stored"
}
func (err *Error) Init() {
	err.NotFound = "not found!"
	err.NotDefined = "not defined!"
}
