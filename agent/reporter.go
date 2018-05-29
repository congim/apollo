package agent

import "sync"

// Reporters ...
var Reporters map[string]Reporter
var reportLock sync.RWMutex

// Reporter ...
type Reporter interface {
	// Start start reporter server
	Start() error
	// Stop any connections to the Output
	Stop() error
	// Write  put data
	Write(interface{}) error
	// Description returns a one-sentence description on the Output
	Description() string
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
