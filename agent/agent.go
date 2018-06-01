package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
	"github.com/shaocongcong/xxxx/common"
	"go.uber.org/zap"
)

// Agent ...
type Agent struct {
	Name       string
	Collectors map[string]Collector
	Reporters  map[string]Reporter
	Processer  *Processer
}

var gAgent *Agent

/************************************************************************************************/
/*									Agent	Add												*/
/************************************************************************************************/

// AddData ...
func (agent *Agent) AddData(name string, iTbl *ast.Table) {
	// 查看是否有该类数据采集多实现
	data, ok := Datas[name]
	if !ok {
		return
	}

	dataC := &DataConf{}

	err := toml.UnmarshalTable(iTbl, dataC)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal collector: ", err)
	}

	Logger.Info("Add", zap.String("data", name), zap.Any("dataC", dataC))

	// 添加到Conf 中
	Conf.DataC[name] = dataC
	// 添加并初始化到 Agent processer
	agent.Processer.InitChannels(name, dataC, data)
}

// AddReporter ...
func (agent *Agent) AddReporter(name string, iTbl *ast.Table) {
	// 过滤
	if _, ok := Conf.ReporterFilters[name]; ok {
		Logger.Info("Pass", zap.String("reporter", name))
		return
	}

	reporter, ok := Reporters[name]
	if !ok {
		return
	}

	err := toml.UnmarshalTable(iTbl, reporter)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal reporter: ", err)
	}

	Logger.Info("Add", zap.String("reporter", name))
	// 添加到Agent Reporters
	agent.Reporters[name] = reporter
}

// AddCollector ...
func (agent *Agent) AddCollector(name string, iTbl *ast.Table) {
	// 过滤
	if _, ok := Conf.CollectorFilters[name]; ok {
		Logger.Info("Pass", zap.String("collector", name))
		return
	}

	collector, ok := Collectors[name]
	if !ok {
		return
	}

	err := toml.UnmarshalTable(iTbl, collector)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal collector: ", err)
	}

	Logger.Info("Add", zap.String("collector", name))

	// 添加到Agent Collectors
	agent.Collectors[name] = collector
}

// writer ...
func (agent *Agent) writer(name string, datas interface{}) error {
	return agent.Processer.Route(name, datas)
}

func (agent *Agent) init() error {
	agent.setAgentName()
	// init reporters
	for name, reporter := range agent.Reporters {
		if err := reporter.Init(); err != nil {
			Logger.Fatal("reporter Init", zap.Error(err), zap.String("name", name))
		}
		if err := reporter.Start(); err != nil {
			Logger.Fatal("reporter Start", zap.Error(err), zap.String("name", name))
		}
		Logger.Info("reporter start", zap.String("name", name))
	}

	// init processer
	{
		// 为processer的channel添加reporter映射关系
		for name, dataC := range Conf.DataC {
			for _, rname := range dataC.Reporters {
				reporter, ok := agent.Reporters[rname]
				if ok {
					agent.Processer.AddReporter(name, reporter)
				} else {
					Logger.Warn("add reporter", zap.String("err", "unfind reporter"), zap.String("name", name))
				}
			}
		}
		if err := agent.Processer.Start(); err != nil {
			Logger.Fatal("Processer init", zap.Error(err))
		}
	}

	// init collectors
	for name, collector := range agent.Collectors {
		if err := collector.Init(); err != nil {
			Logger.Fatal("collector Init", zap.Error(err), zap.String("name", name))
		}
		if err := collector.Start(); err != nil {
			Logger.Fatal("collector Start", zap.Error(err), zap.String("name", name))
		}
		Logger.Info("collector start", zap.String("name", name))
	}
	return nil
}

func (agent *Agent) setAgentName() error {
	// 使用默认设置
	if !Conf.AgentC.UseEnv {
		if agent.Name == "" {
			// 获取主机名
			host, err := common.GetHostName()
			if err != nil {
				Logger.Fatal("setAgentName", zap.Error(err))
			}
			agent.Name = host
		}
	} else {
		// 使用环境变量
		if Conf.AgentC.Env == "" {
			Logger.Fatal("setAgentName", zap.Error(fmt.Errorf("env is nil")))
		}
		name := os.Getenv(Conf.AgentC.Env)
		if name == "" {
			Logger.Fatal("setAgentName", zap.Error(fmt.Errorf("get env is nil")), zap.String("env", Conf.AgentC.Env))
		}
		agent.Name = name
	}
	Logger.Info("AgentName", zap.String("name", agent.Name))
	return nil
}

// New new agent
func New() *Agent {
	return &Agent{
		Collectors: make(map[string]Collector),
		Reporters:  make(map[string]Reporter),
		Processer:  NewProcesser(),
	}
}

// Start start agent server
func Start() error {
	if err := gAgent.init(); err != nil {
		Logger.Fatal("Start", zap.Error(err))
	}

	Logger.Info("Start")
	return nil
}

// Writer ...
func Writer(name string, in interface{}) error {
	if gAgent != nil {
		return gAgent.writer(name, in)
	}
	return nil
}

// Name ...
func Name() string {
	if gAgent != nil {
		return gAgent.Name
	}
	return ""
}

// Stop stop agent server
func Stop() error {
	if gAgent != nil {
		for name, collector := range gAgent.Collectors {
			collector.Close()
			Logger.Info("stop collector", zap.String("name", name))
		}
		for name, reporter := range gAgent.Reporters {
			reporter.Close()
			Logger.Info("stop reporter", zap.String("name", name))
		}
		gAgent.Processer.Close()
		Logger.Info("stop Processer")
	}
	Logger.Info("Stop")
	return nil
}

// init get new agent
func init() {
	gAgent = New()
}
