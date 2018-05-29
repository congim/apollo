package system

import (
	"fmt"
	"time"

	"github.com/shaocongcong/apollo/agent"
	"github.com/shaocongcong/apollo/common"
	"github.com/shaocongcong/apollo/plugins/data/system"
	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
)

// CPUStats ....
type CPUStats struct {
	ps        PS
	lastStats []cpu.TimesStat
	stop      chan bool
	data      agent.Data
	PerCPU    bool `toml:"percpu"`
	TotalCPU  bool `toml:"totalcpu"`
	Interval  int  `toml:"interval"`
}

// NewCPUStats ...
func NewCPUStats(ps PS) *CPUStats {
	return &CPUStats{
		ps: ps,
	}
}

// Gather ....
func (cpu *CPUStats) Gather() error {
	times, err := cpu.ps.CPUTimes(cpu.PerCPU, cpu.TotalCPU)
	if err != nil {
		return fmt.Errorf("error getting CPU info: %s", err)
	}
	now := time.Now()

	for i, cts := range times {
		tags := map[string]string{
			"cpu": cts.CPU,
		}
		total := totalCpuTime(cts)
		fields := map[string]interface{}{}
		// Add in percentage
		if len(cpu.lastStats) == 0 {
			metric := &system.Metric{
				Name:     "cpu",
				Tags:     tags,
				Fields:   fields,
				Time:     now.Unix(),
				Interval: cpu.Interval,
			}
			cpu.data.Writer([]*system.Metric{metric})
			continue
		}
		lastCts := cpu.lastStats[i]
		lastTotal := totalCpuTime(lastCts)
		totalDelta := total - lastTotal

		if totalDelta < 0 {
			cpu.lastStats = times
			return fmt.Errorf("Error: current total CPU time is less than previous total CPU time")
		}

		if totalDelta == 0 {
			continue
		}

		fields["usage_user"] = 100 * (cts.User - lastCts.User) / totalDelta
		fields["usage_idle"] = 100 * (cts.Idle - lastCts.Idle) / totalDelta
		fields["usage_iowait"] = 100 * (cts.Iowait - lastCts.Iowait) / totalDelta
		metric := &system.Metric{
			Name:     "cpu",
			Tags:     tags,
			Fields:   fields,
			Time:     now.Unix(),
			Interval: cpu.Interval,
		}
		cpu.data.Writer([]*system.Metric{metric})
	}

	cpu.lastStats = times

	return nil
}

// Init ...
func (cpu *CPUStats) Init() error {
	go cpu.start()
	return nil
}

func (cpu *CPUStats) start() error {
	ticker := time.NewTicker(time.Duration(cpu.Interval) * time.Second)
	defer func() {
		if err := recover(); err != nil {
			common.PrintStack(true)
			agent.Logger.Error("cpu init", zap.Any("err", err))
		}
	}()
	defer ticker.Stop()

	for {
		select {
		case <-cpu.stop:
			return nil
		case <-ticker.C:
			cpu.Gather()
			continue
		}
	}
}

// Stop ...
func (cpu *CPUStats) Stop() error {
	cpu.stop <- true
	return nil
}

// AddData ...
func (cpu *CPUStats) AddData(data agent.Data) error {
	cpu.data = data
	return nil
}

// Description ...
func (cpu *CPUStats) Description() string {
	return "Read metrics about cpu usage"
}

func totalCpuTime(t cpu.TimesStat) float64 {
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal +
		t.Guest + t.GuestNice + t.Idle
	return total
}

func init() {
	agent.AddCollector("cpu", &CPUStats{
		stop:     make(chan bool, 1),
		PerCPU:   true,
		TotalCPU: true,
		ps:       &systemPS{},
	})
}
