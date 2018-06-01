package agent

import (
	"io/ioutil"
	"log"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
)

// CommonConfig ...
type CommonConfig struct {
	Version  string
	IsDebug  bool
	LogLevel string
	LogPath  string
	Hostname string
}

type EtcdConf struct {
	Addrs     []string
	Dltimeout int
	Rqtimeout int
	WatchDir  string
	TTL       int64
}

// AgentConf ...
type AgentConf struct {
	Name   string
	UseEnv bool
	Env    string
	Size   int
	// Interval        int // Interval at which to gather information
	// FlushInterval   int // FlushInterval is the Interval at which to flush data

}

// DataConf ...
type DataConf struct {
	ChanSize   int
	Collectors []string
	Reporters  []string
}

// Config ...
type Config struct {
	Common           *CommonConfig
	AgentC           *AgentConf
	EtcdC            *EtcdConf
	DataC            map[string]*DataConf
	CollectorFilters map[string]struct{}
	ReporterFilters  map[string]struct{}
}

var Conf *Config

func LoadConfig(file string) {
	// init the new  config params
	Conf = &Config{
		Common:           &CommonConfig{},
		AgentC:           &AgentConf{},
		EtcdC:            &EtcdConf{},
		DataC:            make(map[string]*DataConf),
		CollectorFilters: make(map[string]struct{}),
		ReporterFilters:  make(map[string]struct{}),
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("[FATAL] load config: ", err)
	}

	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("[FATAL] parse config: ", err)
	}

	// parse common config
	parseCommon(tbl)

	// init log logger
	InitLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)

	// parse agent config
	parseAgent(tbl)

	// // parse etcd config
	// parseEtcd(tbl)

	// parse datas config
	parseDatas(tbl)

	// parse filters config
	parseFilters(tbl)

	// parse collector config
	parseCollector(tbl)

	// parse reporter config
	parseReporter(tbl)

	// log.Println("agent conf : ", Conf.AgentC)

}

func parseFilters(tbl *ast.Table) {
	if val, ok := tbl.Fields["global_filters"]; ok {
		if subTbl, ok := val.(*ast.Table); ok {
			if node, ok := subTbl.Fields["collectorsdrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.CollectorFilters[str.Value] = struct{}{}
							}
						}
					}
				}
			}
		}
	}

	if val, ok := tbl.Fields["global_filters"]; ok {
		if subTbl, ok := val.(*ast.Table); ok {
			if node, ok := subTbl.Fields["reportersdrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.ReporterFilters[str.Value] = struct{}{}
							}
						}
					}
				}
			}
		}
	}
}

func parseAgent(tbl *ast.Table) {
	if val, ok := tbl.Fields["agent"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}
		err := toml.UnmarshalTable(subTbl, Conf.AgentC)
		if err != nil {
			log.Fatalln("[FATAL] parseAgent: ", err, subTbl)
		}
	}
}

// func parseEtcd(tbl *ast.Table) {
// 	if val, ok := tbl.Fields["etcd"]; ok {
// 		subTbl, ok := val.(*ast.Table)
// 		if !ok {
// 			log.Fatalln("[FATAL] : ", subTbl)
// 		}
// 		err := toml.UnmarshalTable(subTbl, Conf.EtcdC)
// 		if err != nil {
// 			log.Fatalln("[FATAL] parseEtcd: ", err, subTbl)
// 		}
// 	}
// }

func parseCommon(tbl *ast.Table) {
	if val, ok := tbl.Fields["common"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}
		err := toml.UnmarshalTable(subTbl, Conf.Common)
		if err != nil {
			log.Fatalln("[FATAL] parseCommon: ", err, subTbl)
		}
	}
}

// func parseDatas(tbl *ast.Table) {

// 	if val, ok := tbl.Fields["datas"]; ok {
// 		subTbl, _ := val.(*ast.Table)
// 		for pn, pt := range subTbl.Fields {
// 			switch iTbl := pt.(type) {
// 			case *ast.Table:
// 				gAgent.AddData(pn, iTbl)
// 			case []*ast.Table:
// 				for _, t := range iTbl {
// 					gAgent.AddData(pn, t)
// 				}
// 			default:
// 				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
// 			}
// 		}
// 	}
// }

func parseDatas(tbl *ast.Table) {
	if val, ok := tbl.Fields["datas"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			switch iTbl := pt.(type) {
			case *ast.Table:
				gAgent.AddData(pn, iTbl)
			case []*ast.Table:
				for _, t := range iTbl {
					gAgent.AddData(pn, t)
				}
			default:
				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
			}
		}
	}
}

func parseCollector(tbl *ast.Table) {
	if val, ok := tbl.Fields["collectors"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			switch iTbl := pt.(type) {
			case *ast.Table:
				gAgent.AddCollector(pn, iTbl)
			case []*ast.Table:
				for _, t := range iTbl {
					gAgent.AddCollector(pn, t)
				}
			default:
				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
			}
		}
	}
}

func parseReporter(tbl *ast.Table) {
	if val, ok := tbl.Fields["reporters"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			switch iTbl := pt.(type) {
			case *ast.Table:
				gAgent.AddReporter(pn, iTbl)
			case []*ast.Table:
				for _, t := range iTbl {
					gAgent.AddReporter(pn, t)
				}
			default:
				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
			}
		}
	}
}
