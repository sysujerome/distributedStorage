package common

import (
	"sync"
)

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

type SMap struct {
	sync.RWMutex
	Map map[string]string
}

func (l *SMap) ReadMap(key string) (string, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *SMap) WriteMap(key string, value string) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}

func (l *SMap) GetLen() int {
	l.RLock()
	len := len(l.Map)
	l.RUnlock()
	return len
}
func (l *SMap) Del(key string) {
	l.Lock()
	delete(l.Map, key)
	l.Unlock()
}
func (l *SMap) Range() (string, string, bool) {
	l.Lock()
	if len(l.Map) <= 0 {
		l.Unlock()
		return "", "", false
	}
	var key string
	var value string
	for k, v := range l.Map {
		key = k
		value = v
		break
	}
	l.Unlock()
	return key, value, true
}

type IMap struct {
	sync.RWMutex
	Map map[int64]string
}

func (l *IMap) ReadMap(key int64) (string, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *IMap) WriteMap(key int64, value string) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}

func (l *IMap) GetLen() int {
	l.RLock()
	len := len(l.Map)
	l.RUnlock()
	return len
}
func (l *IMap) Del(key int64) {
	l.Lock()
	delete(l.Map, key)
	l.Unlock()
}
func (l *IMap) Range() (int64, string, bool) {
	l.Lock()
	if len(l.Map) <= 0 {
		l.Unlock()
		return 0, "", false
	}
	var key int64
	var value string
	for k, v := range l.Map {
		key = k
		value = v
		break
	}
	l.Unlock()
	return key, value, true
}

type IIMap struct {
	sync.RWMutex
	Map map[int64]int64
}

func (l *IIMap) ReadMap(key int64) (int64, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *IIMap) WriteMap(key int64, value int64) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}

func (l *IIMap) GetLen() int {
	l.RLock()
	len := len(l.Map)
	l.RUnlock()
	return len
}
func (l *IIMap) Del(key int64) {
	l.Lock()
	delete(l.Map, key)
	l.Unlock()
}
func (l *IIMap) Range() (int64, int64, bool) {
	l.Lock()
	if len(l.Map) <= 0 {
		l.Unlock()
		return 0, 0, false
	}
	var key int64
	var value int64
	for k, v := range l.Map {
		key = k
		value = v
		break
	}
	l.Unlock()
	return key, value, true
}
