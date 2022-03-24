package main

import "fmt"

type Status struct {
	Ok string `default:"default-a"`

	Failed string `default:"default-f"`
	Stored string
}

func NewStatus() Status {
	status := Status{}
	status.Ok = "OK"
	status.Failed = "failed"
	status.Stored = "stored"
	return status
}

func (sta *Status) fill_defaults() {

	// setting default values
	// if no values present
	if sta.Ok == "" {
		sta.Ok = "OK"
	}

	if sta.Failed == "" {
		sta.Failed = "failed"
	}

	if sta.Stored == "" {
		sta.Stored = "stored"
	}
}
func main() {
	var sta Status
	sta.Ok = "???"
	fmt.Println(sta.Ok)
}
