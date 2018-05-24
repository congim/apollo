package agent

import (
	"sync"
)

// Datas ...
var Datas map[string]Data
var dataLock sync.RWMutex

// Data ...
type Data interface {
	Init() error
	Close() error
	Description() string
	Encoder() error
	Decoder() error
	Reader() (interface{}, error)
	Writer(interface{}) error
	GetReporters() []string
}

// AddData ....
func AddData(name string, data Data) {
	dataLock.Lock()
	defer dataLock.Unlock()

	if Datas == nil {
		Datas = make(map[string]Data)
	}
	Datas[name] = data
}
