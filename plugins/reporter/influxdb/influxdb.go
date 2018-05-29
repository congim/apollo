package influxdb

import (
	"log"

	"github.com/shaocongcong/apollo/agent"
)

// InfluxDB influxdb server
type InfluxDB struct {
}

// Start start reporter server
func (influxDB *InfluxDB) Start() error {
	log.Println("start influx db")
	return nil
}

// Stop any connections to the Output
func (influxDB *InfluxDB) Stop() error {
	return nil
}

// Write  put data
func (influxDB *InfluxDB) Write(data interface{}) error {
	log.Println("influxdb write", data)
	return nil
}

// Description returns a one-sentence description on the Output
func (influxDB *InfluxDB) Description() string {
	return "InfluxDB"
}

// var Influxdb *InfluxDB

// // InfluxDB influxDB
// type InfluxDB struct {
// 	// URL is only for backwards compatability
// 	URL              string
// 	URLs             []string `toml:"urls"`
// 	Username         string
// 	Password         string
// 	Database         string
// 	UserAgent        string
// 	RetentionPolicy  string
// 	WriteConsistency string
// 	Timeout          common.Duration
// 	UDPPayload       int `toml:"udp_payload"`
// 	// Precision is only here for legacy support. It will be ignored.
// 	Precision string
// 	TimedTask int
// 	DataLen   int
// 	stopc     chan bool
// 	taskC     chan []*system.Metric

// 	conns []client.Client
// }

// // NewInfluxDB new InfluxDB
// func NewInfluxDB(urls string, database string, writeConsistency string, timeout int, timedTask int, datalen int) *InfluxDB {

// 	return &InfluxDB{
// 		URLs:             strings.Split(urls, ","),
// 		Database:         database,
// 		WriteConsistency: writeConsistency,
// 		Timeout:          common.Duration{time.Duration(timeout) * time.Second},
// 		TimedTask:        timedTask,
// 		DataLen:          datalen,
// 		stopc:            make(chan bool, 1),
// 		taskC:            make(chan []*system.Metric, datalen),
// 	}
// }

// var sampleConfig = `
//   ## The full HTTP or UDP endpoint URL for your InfluxDB instance.
//   ## Multiple urls can be specified as part of the same cluster,
//   ## this means that only ONE of the urls will be written to each interval.
//   # urls = ["udp://localhost:8089"] # UDP endpoint example
//   urls = ["http://localhost:8086"] # required
//   ## The target database for metrics (telegraf will create it if not exists).
//   database = "telegraf" # required

//   ## Retention policy to write to. Empty string writes to the default rp.
//   retention_policy = ""
//   ## Write consistency (clusters only), can be: "any", "one", "quorom", "all"
//   write_consistency = "any"

//   ## Write timeout (for the InfluxDB client), formatted as a string.
//   ## If not provided, will default to 5s. 0s means no timeout (not recommended).
//   timeout = "5s"
//   # username = "telegraf"
//   # password = "metricsmetricsmetricsmetrics"
//   ## Set the user agent for HTTP POSTs (can be useful for log differentiation)
//   # user_agent = "telegraf"
//   ## Set UDP payload size, defaults to InfluxDB UDP Client default (512 bytes)
//   # udp_payload = 512

//   ## Optional SSL Config
//   # ssl_ca = "/etc/telegraf/ca.pem"
//   # ssl_cert = "/etc/telegraf/cert.pem"
//   # ssl_key = "/etc/telegraf/key.pem"
//   ## Use SSL but skip chain & host verification
//   # insecure_skip_verify = false
// `

// // Connect connect infuxDB
// func (i *InfluxDB) Connect() error {
// 	var urls []string
// 	for _, u := range i.URLs {
// 		urls = append(urls, u)
// 	}

// 	// Backward-compatability with single Influx URL config files
// 	// This could eventually be removed in favor of specifying the urls as a list
// 	if i.URL != "" {
// 		urls = append(urls, i.URL)
// 	}

// 	var conns []client.Client
// 	for _, u := range urls {
// 		switch {
// 		case strings.HasPrefix(u, "udp"):
// 			parsed_url, err := url.Parse(u)
// 			if err != nil {
// 				return err
// 			}

// 			if i.UDPPayload == 0 {
// 				i.UDPPayload = client.UDPPayloadSize
// 			}
// 			c, err := client.NewUDPClient(client.UDPConfig{
// 				Addr:        parsed_url.Host,
// 				PayloadSize: i.UDPPayload,
// 			})
// 			if err != nil {
// 				return err
// 			}
// 			conns = append(conns, c)
// 		default:
// 			// If URL doesn't start with "udp", assume HTTP client
// 			c, err := client.NewHTTPClient(client.HTTPConfig{
// 				Addr:      u,
// 				Username:  i.Username,
// 				Password:  i.Password,
// 				UserAgent: i.UserAgent,
// 				Timeout:   i.Timeout.Duration,
// 			})
// 			if err != nil {
// 				return err
// 			}

// 			err = createDatabase(c, i.Database)
// 			if err != nil {
// 				log.Println("Database creation failed: " + err.Error())
// 				continue
// 			}

// 			conns = append(conns, c)
// 		}
// 	}

// 	i.conns = conns
// 	rand.Seed(time.Now().UnixNano())
// 	return nil
// }

// func createDatabase(c client.Client, database string) error {
// 	// Create Database if it doesn't exist
// 	_, err := c.Query(client.Query{
// 		Command: fmt.Sprintf("CREATE DATABASE \"%s\"", database),
// 	})
// 	return err
// }

// // Close stop InfluxDB
// func (i *InfluxDB) Close() error {
// 	var errS string
// 	for j, _ := range i.conns {
// 		if err := i.conns[j].Close(); err != nil {
// 			errS += err.Error()
// 		}
// 	}
// 	if errS != "" {
// 		return fmt.Errorf("output influxdb close failed: %s", errS)
// 	}
// 	close(i.taskC)
// 	close(i.stopc)
// 	return nil
// }

// // Submit 提交任务
// func (i *InfluxDB) Submit(data []*system.Metric) {
// 	i.taskC <- data
// }

// // Choose a random server in the cluster to write to until a successful write
// // occurs, logging each unsuccessful. If all servers fail, return error.
// func (i *InfluxDB) Write(metrics []*system.Metric) error {
// 	if len(i.conns) == 0 {
// 		err := i.Connect()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
// 		Database:         i.Database,
// 		RetentionPolicy:  i.RetentionPolicy,
// 		WriteConsistency: i.WriteConsistency,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	for _, metric := range metrics {

// 		pt, err := client.NewPoint(metric.Name, metric.Tags, metric.Fields, time.Unix(metric.Time, 0))
// 		if err != nil {
// 			alog.Logger.Error("InfluxDB Write", zap.Error(err))
// 			return err
// 		}
// 		// alog.Logger.Debug("InfluxDB Write", zap.Any("@metric", metric))
// 		bp.AddPoint(pt)
// 	}

// 	// This will get set to nil if a successful write occurs
// 	err = errors.New("Could not write to any InfluxDB server in cluster")

// 	p := rand.Perm(len(i.conns))
// 	for _, n := range p {
// 		if e := i.conns[n].Write(bp); e != nil {
// 			alog.Logger.Error("InfluxDB Write", zap.Error(e))
// 			// If the database was not found, try to recreate it
// 			if strings.Contains(e.Error(), "database not found") {
// 				if errc := createDatabase(i.conns[n], i.Database); errc != nil {
// 					alog.Logger.Error("ERROR: Database "+i.Database+" not found and failed to recreate\n", zap.Error(errc))
// 				}
// 			}
// 		} else {
// 			err = nil
// 			break
// 		}
// 	}

// 	return err
// }

// // Start infuxDB
// func (i *InfluxDB) Start() {
// 	if err := i.Connect(); err != nil {
// 		log.Fatal("InfluxDB Connect failed, err message is ", err)
// 	}
// 	Influxdb = i
// 	// 启动定时任务
// 	go i.timeTask()
// }

// // timeTask 定时任务
// func (i *InfluxDB) timeTask() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			alog.Logger.Error("InfluxDB", zap.Any("recover", err))
// 		}
// 	}()

// 	// log.Println("timeTask", i.TimedTask)

// 	jobpool := make([]*system.Metric, 0, i.DataLen*2)
// 	TimeC := time.Tick(time.Duration(i.TimedTask) * time.Millisecond)
// 	for {
// 		select {
// 		case <-i.stopc:
// 			return
// 		case job, ok := <-i.taskC:
// 			if ok {
// 				jobpool = append(jobpool, job...)
// 				if len(jobpool) >= i.DataLen {
// 					err := i.Write(jobpool)
// 					if err != nil {
// 						alog.Logger.Error("InfluxDB", zap.Error(err))
// 					}
// 					jobpool = jobpool[:0]
// 				}
// 			}
// 			break
// 		case <-TimeC:
// 			if len(jobpool) > 0 {
// 				err := i.Write(jobpool)
// 				if err != nil {
// 					alog.Logger.Error("InfluxDB", zap.Error(err))
// 				}
// 				jobpool = jobpool[:0]
// 			}
// 			break
// 		}
// 	}
// }

// // Compute infuxDB
// func (i *InfluxDB) Compute(metrics []*system.Metric) error {
// 	// log.Println("influxDB data is", metrics)
// 	return i.Write(metrics)
// }

// AddReporter(name string, reporter Reporter)

func init() {
	influxdb := &InfluxDB{}
	agent.AddReporter("influxdb", influxdb)
}
