package agent

import (
	"fmt"

	"go.uber.org/zap"
)

// Processer ...
type Processer struct {
	Channels     map[string]chan interface{}
	consumeRoute map[string][]Reporter
	stop         chan bool
}

// NewProcesser ....
func NewProcesser() *Processer {
	return &Processer{
		Channels:     make(map[string]chan interface{}),
		consumeRoute: make(map[string][]Reporter),
		stop:         make(chan bool, 1),
	}
}

// InitChannels 。。。
func (pro *Processer) InitChannels(name string, dataC *DataConf, data Data) error {
	channel, ok := pro.Channels[name]
	if ok {
		return nil
	}
	channel = make(chan interface{}, dataC.ChanSize)
	pro.Channels[name] = channel
	Logger.Info("initChannel", zap.String("name", name), zap.Int("size", dataC.ChanSize))
	return nil
}

// Route 。。。
func (pro *Processer) Route(name string, datas interface{}) error {
	channel, ok := pro.Channels[name]
	if !ok {
		return fmt.Errorf("unfind channel, collector name is %s", name)
	}
	channel <- datas
	return nil
}

// Start 。。。
func (pro *Processer) Start() error {
	for name, channel := range pro.Channels {
		go pro.consume(name, channel)
	}
	return nil
}

// AddReporter 。。。
func (pro *Processer) AddReporter(name string, reporter Reporter) {
	consumeRoute, ok := pro.consumeRoute[name]
	if !ok {
		consumeRoute = make([]Reporter, 0)
		consumeRoute = append(consumeRoute, reporter)
		pro.consumeRoute[name] = consumeRoute
	} else {
		consumeRoute = append(consumeRoute, reporter)
	}
}

// Close stop processer
func (pro *Processer) Close() {
	pro.stop <- true
	defer close(pro.stop)
}

func (pro *Processer) consume(name string, channel chan interface{}) {
	defer func() {
		if err := recover(); err != nil {
			Logger.Error("cosume", zap.Any("err", err))
		}
	}()
	for {
		select {
		case datas, ok := <-channel:
			if ok {
				reporters, ok := pro.consumeRoute[name]
				if ok {
					for _, reporter := range reporters {
						if err := reporter.Writer(datas); err != nil {
							Logger.Warn("reporter Writer", zap.Error(err), zap.String("name", name))
						}
					}
				}
			}
			break
		case <-pro.stop:
			break
		}
	}
}
