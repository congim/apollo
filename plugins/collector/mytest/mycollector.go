package mytest

import (
	"log"

	"github.com/shaocongcong/apollo/agent"
)

type MyCollector struct {
	Addr  string
	Descr []string
	data  agent.Data
	// Reporters []string
}

func (mc *MyCollector) Description() string {
	log.Println(mc.Addr)
	log.Println(mc.Descr)
	// log.Println(mc.Reporters)
	return "测试"
}

func (mc *MyCollector) AddData(c agent.Data) error {
	mc.data = c
	// metric := &system.Metric{}
	// // c.AddChan(metric)
	// log.Println(metric)

	return nil
}

func (mc *MyCollector) Gather() error {
	return nil
}

func (mc *MyCollector) Init() error {
	return nil
}
func (mc *MyCollector) Stop() error {
	return nil
}

func init() {
	mc := &MyCollector{
		Descr: make([]string, 0),
	}
	agent.AddCollector("mycollector", mc)
}
