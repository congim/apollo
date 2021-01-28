package system

import (
	"github.com/congim/apollo/agent"
)

// Metric ...
type Metric struct {
	Name     string                 `msg:"n"`
	Tags     map[string]string      `msg:"ts"`
	Fields   map[string]interface{} `msg:"f"`
	Time     int64                  `msg:"t"`
	Interval int                    `msg:"i"`
}

// Encoder 编码
func (m *Metric) Encoder() error {
	return nil
}

// Decoder 编码
func (m *Metric) Decoder() error {
	return nil
}

func init() {
	agent.AddData("metric", &Metric{})
}
