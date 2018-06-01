package agent

import "sync"

// Reporters ...
var Reporters map[string]Reporter
var rlock sync.RWMutex

// Reporter ...
type Reporter interface {
	Init() error
	Start() error
	Close() error
	Writer(interface{}) error
}

// AddReporter ....
func AddReporter(name string, reporter Reporter) {
	rlock.Lock()
	defer rlock.Unlock()

	if Reporters == nil {
		Reporters = make(map[string]Reporter)
	}
	Reporters[name] = reporter
}
