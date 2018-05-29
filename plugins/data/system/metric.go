package system

import (
	"github.com/shaocongcong/apollo/agent"
)

// Metricer metricer
type Metricer struct {
	Collectors []string
	Reporters  []string
	ChanSize   int
	Name       string
	tchan      chan []*Metric
}

// 16:58:30 udp.go:142: docker 监控测试打印 Name cpu
// 16:58:30 udp.go:143: docker 监控测试打印 Tags map[app:zeus cpu:cpu-total host:10.7.14.201]
// 16:58:30 udp.go:144: docker 监控测试打印 Fields map[usage_idle:91.98396789442992 usage_iowait:0 usage_user:7.014028052713285]
// 16:58:30 udp.go:145: docker 监控测试打印 Time 1527497900
// 16:58:30 udp.go:146: docker 监控测试打印 Interval 15

// Metric ...
type Metric struct {
	Name     string                 `msg:"n"`
	Tags     map[string]string      `msg:"ts"`
	Fields   map[string]interface{} `msg:"f"`
	Time     int64                  `msg:"t"`
	Interval int                    `msg:"i"`
}

// Init 初始化缓存管道
func (m *Metricer) Init() error {
	// init chan
	if m.ChanSize <= 0 {
		m.ChanSize = 1000
	}
	m.tchan = make(chan []*Metric, m.ChanSize)
	return nil
}

// Close 释放资源
func (m *Metricer) Close() error {
	if m.tchan != nil {
		close(m.tchan)
	}
	return nil
}

// Reader Reporter插件服务读取数据使用
func (m *Metricer) Reader() (interface{}, error) {
	mterics, ok := <-m.tchan
	if !ok {
		return nil, nil
	}
	return mterics, nil
}

// Writer collector插件写数据专用
func (m *Metricer) Writer(metric interface{}) error {
	switch data := metric.(type) {
	case []*Metric:
		m.tchan <- data
	}
	return nil
}

// Description 描述
func (m *Metricer) Description() string {
	return m.Name
}

// Encoder 编码
func (m *Metricer) Encoder() error {
	return nil
}

// Decoder 编码
func (m *Metricer) Decoder() error {
	return nil
}

// GetCollectors get collectors
func (m *Metricer) GetCollectors() []string {
	return m.Collectors
}

// GetReporters get reporters
func (m *Metricer) GetReporters() []string {
	return m.Reporters
}

func init() {
	metricer := &Metricer{
		Name: "metric",
	}
	metricer.Init()
	agent.AddData("metric", metricer)
}
