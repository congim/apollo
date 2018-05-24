package system

import "github.com/shaocongcong/apollo/agent"

// Metricer metricer
type Metricer struct {
	Collectors []string
	Reporters  []string
	ChanSize   int
	tchan      chan []*Metric
}

// Metric ...
type Metric struct {
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
	return nil, nil
}

// Writer collector插件写数据专用
func (m *Metricer) Writer(interface{}) error {
	return nil
}

// Description 描述
func (m *Metricer) Description() string {
	return "Metricer"
}

// Encoder 编码
func (m *Metricer) Encoder() error {
	return nil
}

// Decoder 编码
func (m *Metricer) Decoder() error {
	return nil
}

// GetReporters get reporters
func (m *Metricer) GetReporters() []string {
	return m.Reporters
}

func init() {
	metricer := &Metricer{}
	metricer.Init()
	agent.AddData("metric", metricer)
}
