package agent

import "sync"

// Collectors ...
var Collectors map[string]Collector
var lock sync.RWMutex

// Collector ...
type Collector interface {
	Description() string
	Gather() error
	Run() error
	Stop() error
	AddData(Data) error
}

// AddCollector ....
func AddCollector(name string, collector Collector) {
	lock.Lock()
	defer lock.Unlock()

	if Collectors == nil {
		Collectors = make(map[string]Collector)
	}
	Collectors[name] = collector
}

// type AdminCollector struct {
// }

// func NewAdminCollector() *AdminCollector {
// 	return &AdminCollector{}
// }

// Collector

// Collector服务要为Agent提供采集以及控制采集频率
// 同理Report也需要为Agent提供上报服务以及上报频率
