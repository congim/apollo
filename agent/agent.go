package agent

import (
	"log"
	"sync"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
	"go.uber.org/zap"
)

// Agent agent for collect system msgs
type Agent struct {
	// Collectors []Collector
	// Reporters  []Reporter
	Collectors map[string]Collector
	Reporters  map[string]Reporter
	Datas      map[string]Data
}

var gAgent *Agent

// CreateAgent create new agent
func CreateAgent() (*Agent, error) {
	agent := &Agent{
		// Collectors: make([]Collector, 0),
		// Reporters:  make([]Reporter, 0),
		Collectors: make(map[string]Collector),
		Reporters:  make(map[string]Reporter),
		Datas:      make(map[string]Data),
	}
	gAgent = agent
	return agent, nil
}

// Run run agent server
func (agent *Agent) Run() error {
	// Logger.Info("Run", zap.Int("len", int(len(agent.Collectors))))
	// for _, collector := range agent.Collectors {
	// 	collector.Description()
	// }
	//

	log.Println("Run")
	agent.runCollectors()
	agent.runData()
	log.Println("Run 2", agent.Datas)
	return nil
}

func (agent *Agent) runCollectors() error {
	// Logger.Info("Run", zap.Int("len", int(len(agent.Collectors))))
	// for _, collector := range agent.Collectors {
	// 	collector.Description()
	// }
	//
	return nil
}

func (agent *Agent) runData() error {
	for _, data := range agent.Datas {
		// data.Reader()
		log.Println("get reporters", data.GetReporters())
	}
	return nil
}

// Stop stop agent server
func (agent *Agent) Stop() error {
	return nil
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
	// agent.Collectors = append(agent.Collectors, collector)
}

// AddData ...
func (agent *Agent) AddData(name string, iTbl *ast.Table) {
	data, ok := Datas[name]
	if !ok {
		return
	}

	err := toml.UnmarshalTable(iTbl, data)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal collector: ", err)
	}

	Logger.Info("Add", zap.String("data", name))

	// 添加到Agent Datas
	agent.Datas[name] = data
	// agent.Datas = append(agent.Datas, data)
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
	// agent.Reporters = append(agent.Reporters, reporter)
}

func (agent *Agent) flush() {
	// log.Println("flush")
	// for _, reporter := range agent.Reporters {
	// 	reporter.Write("")
	// }
}

// flusher monitors the metrics input channel and flushes on the minimum interval
func (agent *Agent) flusher(wg *sync.WaitGroup, shutdown chan struct{}, metricC chan []byte) {
	// defer func() {
	// 	wg.Done()
	// 	if err := recover(); err != nil {
	// 		Logger.Fatal("flush fatal error ", zap.Error(err.(error)))
	// 	}
	// }()

	// ticker := time.NewTicker(time.Duration(Conf.AgentC.FlushInterval) * time.Second)
	// for {
	// 	select {
	// 	case <-shutdown:
	// 		agent.flush()
	// 		return
	// 	case <-ticker.C:
	// 		agent.flush()
	// 	case m := <-metricC:
	// 		log.Println(m)
	// 		// for _, o := range Conf.Outputs {
	// 		// 	o.AddMetric(m)
	// 		// }
	// 	}
	// }
}
