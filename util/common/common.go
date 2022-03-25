package common

type Status struct {
	Ok     string
	Failed string
	Stored string
	Moved  string
}

type Error struct {
	NotFound   string
	NotDefined string
	NotWorking string
	Stored     string
	Moved      string
}

type ServerStatus struct {
	Working  string
	Sleep    string
	Spliting string
	Full     string
}

func (sta *Status) Init() {
	sta.Ok = "OK"
	sta.Failed = "failed"
	sta.Stored = "stored"
	sta.Moved = "moved"
}
func (err *Error) Init() {
	err.NotFound = "not found!"
	err.NotDefined = "not defined!"
	err.NotWorking = "the server is sleeping now."
	err.Stored = "the operation is stored."
	err.Moved = "moved"
}

func (ss *ServerStatus) Init() {
	ss.Working = "working"
	ss.Sleep = "sleep"
	ss.Spliting = "spliting"
	ss.Full = "can not split"
}
