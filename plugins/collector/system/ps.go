package system

type PS interface {
	//CPUTimes(perCPU, totalCPU bool) ([]cpu.TimesStat, error)
	// DiskUsage(mountPointFilter []string, fstypeExclude []string) ([]*disk.UsageStat, error)
	// NetIO() ([]net.IOCountersStat, error)
	// NetProto() ([]net.ProtoCountersStat, error)
	// DiskIO() (map[string]disk.IOCountersStat, error)
	// VMStat() (*mem.VirtualMemoryStat, error)
	// SwapStat() (*mem.SwapMemoryStat, error)
	// NetConnections() ([]net.ConnectionStat, error)
}

// func add(acc agent.Accumulator, name string, val float64, tags map[string]string) {
// 	if val >= 0 {
// 		acc.Add(name, val, tags)
// 	}
// }

//type systemPS struct{}
//
//func (s *systemPS) CPUTimes(perCPU, totalCPU bool) ([]cpu.TimesStat, error) {
//	var cpuTimes []cpu.TimesStat
//	if perCPU {
//		if perCPUTimes, err := cpu.Times(true); err == nil {
//			cpuTimes = append(cpuTimes, perCPUTimes...)
//		} else {
//			return nil, err
//		}
//	}
//	if totalCPU {
//		if totalCPUTimes, err := cpu.Times(false); err == nil {
//			cpuTimes = append(cpuTimes, totalCPUTimes...)
//		} else {
//			return nil, err
//		}
//	}
//	return cpuTimes, nil
//}
