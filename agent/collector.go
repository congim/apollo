package agent

import "sync"

// Collectors ...
var Collectors map[string]Collector
var clock sync.RWMutex

// Collector ...
type Collector interface {
	Init() error
	Start() error
	Close() error
}

// AddCollector ....
func AddCollector(name string, collector Collector) {
	clock.Lock()
	defer clock.Unlock()

	if Collectors == nil {
		Collectors = make(map[string]Collector)
	}
	Collectors[name] = collector
}
