package agent

import "sync"

// Datas ...
var Datas map[string]Data
var dlock sync.RWMutex

// Data ...
type Data interface {
	Encoder() error
	Decoder() error
}

// AddData ....
func AddData(name string, data Data) {
	dlock.Lock()
	defer dlock.Unlock()

	if Datas == nil {
		Datas = make(map[string]Data)
	}
	Datas[name] = data
}
