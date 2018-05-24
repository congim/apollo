package agent

import "sync"

// Reporters ...
var Reporters map[string]Reporter
var reportLock sync.RWMutex

// Reporter ...
type Reporter interface {
	Write(interface{}) error
}

// AddReporter ....
func AddReporter(name string, reporter Reporter) {
	reportLock.Lock()
	defer reportLock.Unlock()

	if Reporters == nil {
		Reporters = make(map[string]Reporter)
	}
	Reporters[name] = reporter
}
