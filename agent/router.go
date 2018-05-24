package agent

import "log"

// Router router collecter's data for reporter
type Router struct {
	datas []Data

	data      Data
	reporters []Reporter
}

//

func (router *Router) Run() {
	for _, data := range router.datas {
		go func() {
			for {
				// sleep (flush time)
				log.Println(data)
			}
		}()
	}
}

func (router *Router) start() {

	for {
		data, err := router.data.Reader()
		if err != nil {
			continue
		}
		for _, reporter := range router.reporters {
			reporter.Write(data)
		}
	}
}
