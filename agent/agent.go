package agent

import (
	"fmt"
	"log"
	"time"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
	"go.uber.org/zap"
)

// Agent agent for collect system msgs
type Agent struct {
	Collectors map[string]Collector
	Reporters  map[string]Reporter
	Datas      map[string]Data
	Routers    map[string]*Router
}

var gAgent *Agent

// CreateAgent create new agent
func CreateAgent() (*Agent, error) {
	agent := &Agent{
		Collectors: make(map[string]Collector),
		Reporters:  make(map[string]Reporter),
		Datas:      make(map[string]Data),
		Routers:    make(map[string]*Router),
	}
	gAgent = agent
	return agent, nil
}

// Run run agent server
func (agent *Agent) Run() error {

	// init Data
	if err := agent.initData(); err != nil {
		Logger.Fatal("Run", zap.Error(err))
		return err
	}

	// init collector
	if err := agent.initCollector(); err != nil {
		Logger.Fatal("Run", zap.Error(err))
		return err
	}

	// init reporter
	if err := agent.initReporter(); err != nil {
		Logger.Fatal("Run", zap.Error(err))
		return err
	}

	// init router
	if err := agent.initRouter(); err != nil {
		Logger.Fatal("Run", zap.Error(err))
		return err
	}

	// start consume
	agent.Consume()

	Logger.Info("Agent Run")
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

/************************************************************************************************/
/*									Agent	Init												*/
/************************************************************************************************/

// initData  init data
func (agent *Agent) initData() error {
	for name, data := range agent.Datas {
		if err := data.Init(); err != nil {
			Logger.Fatal("initData", zap.Error(err), zap.String("data name", name))
		}

		// add data into Collector
		for _, cname := range data.GetCollectors() {
			if err := agent.collectorAddData(cname, data); err != nil {
				Logger.Warn("collectorAddData", zap.Error(err), zap.String("data name", name))
				continue
			}
		}

		Logger.Info("initData", zap.String("name", name))
	}
	return nil
}

func (agent *Agent) initReporter() error {
	for name, reporter := range agent.Reporters {
		log.Println(name, reporter)
		if err := reporter.Start(); err != nil {
			Logger.Error("initReporter", zap.String("name", name), zap.Error(err))
		}
	}
	return nil
}

func (agent *Agent) collectorAddData(cname string, data Data) error {
	collector, ok := agent.Collectors[cname]
	if !ok {
		return fmt.Errorf("unfind collector, collector name is %s", cname)
	}
	collector.AddData(data)
	return nil
}

// initCollector init collector
func (agent *Agent) initCollector() error {
	for name, collector := range agent.Collectors {
		if err := collector.Init(); err != nil {
			Logger.Fatal("initCollector", zap.Error(err), zap.String("name", name))
		}
		Logger.Info("initCollector", zap.String("name", name))
	}
	return nil
}

// initRouter ...
func (agent *Agent) initRouter() error {
	for _, data := range agent.Datas {
		router, ok := agent.Routers[data.Description()]
		if !ok {
			router = NewRouter(data.Description(), data)
			agent.Routers[data.Description()] = router
		}
		for _, reporterName := range data.GetReporters() {
			// get reporter
			report, ok := agent.Reporters[reporterName]
			if ok {
				Logger.Info("initRouter", zap.String("add reporter", reporterName))
				router.AddReporter(reporterName, report)
			} else {
				Logger.Warn("initRouter", zap.String("unfind reporter", reporterName))
			}
		}
	}
	return nil
}

/************************************************************************************************/
/*									Agent	Run													*/
/************************************************************************************************/

// Consume ....
func (agent *Agent) Consume() {
	for name, router := range agent.Routers {
		go agent.consume(router.data)
		Logger.Info("Consume", zap.String("router", name))
	}
}

func (agent *Agent) consume(data Data) {
	for {
		mdata, err := data.Reader()
		if err != nil {
			Logger.Warn("consume", zap.String("data name", data.Description()), zap.Error(err))
			time.Sleep(1 * time.Second)
		}

		for _, name := range data.GetReporters() {
			reporter, ok := agent.Reporters[name]
			if !ok {
				continue
			}
			reporter.Write(mdata)
		}
	}
}
