package agent

import (
	"sync"

	"go.uber.org/zap"
)

// Router router collecter's data for reporter
type Router struct {
	sync.RWMutex
	name      string
	data      Data
	reporters map[string]Reporter
}

// NewRouter new router
func NewRouter(name string, data Data) *Router {
	return &Router{
		name:      name,
		data:      data,
		reporters: make(map[string]Reporter),
	}
}

// AddReporter add reporter
func (router *Router) AddReporter(name string, reporter Reporter) {
	router.Lock()
	defer router.Unlock()
	if _, ok := router.reporters[name]; ok {
		Logger.Info("AddReporter", zap.String("exist reporter", name))
		return
	}
	router.reporters[name] = reporter
}

// Run run
func (router *Router) Run() {
	// read data from router.data
	// for  reporters and write data
	// for _, data := range router.datas {
	// 	go func() {
	// 		for {
	// 			// sleep (flush time)
	// 			log.Println(data)
	// 		}
	// 	}()
	// }
}

func (router *Router) start() {
	// for {
	// 	data, err := router.data.Reader()
	// 	if err != nil {
	// 		continue
	// 	}
	// 	for _, reporter := range router.reporters {
	// 		reporter.Write(data)
	// 	}
	// }
}
