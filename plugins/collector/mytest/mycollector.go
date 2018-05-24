package mytest

import (
	"log"

	"github.com/shaocongcong/apollo/agent"
	"github.com/shaocongcong/apollo/plugins/data/system"
)

type MyCollector struct {
	Addr  string
	Descr []string
	// Reporters []string
}

func (mc *MyCollector) Description() string {
	log.Println(mc.Addr)
	log.Println(mc.Descr)
	// log.Println(mc.Reporters)
	return "测试"
}

func (mc *MyCollector) AddData(c agent.Data) error {
	metric := &system.Metric{}
	// c.AddChan(metric)
	log.Println(metric)
	return nil
}

func (mc *MyCollector) Gather() error {
	return nil
}

func (mc *MyCollector) Run() error {
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
