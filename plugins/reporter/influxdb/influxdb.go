package influxdb

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/shaocongcong/apollo/agent"
	"github.com/shaocongcong/apollo/plugins/data/system"
)

// InfluxDB influxdb server
type InfluxDB struct {
	URL              string
	URLs             []string `toml:"urls"`
	Username         string
	Password         string
	DataBase         string `toml:"data_base"`
	RetentionPolicy  string
	UserAgent        string
	WriteConsistency string `toml:"write_consistency"`
	Timeout          int    `toml:"timeout"`
	UDPPayload       int    `toml:"udp_payload"`
	TimeTask         int    `toml:"time_task"`
	DataLen          int    `toml:"data_len"`

	stopc chan bool
	taskC chan []*system.Metric

	conns []client.Client
}

// Init int reporter server
func (influxDB *InfluxDB) Init() error {
	return nil
}

// Start start reporter server
func (influxDB *InfluxDB) Start() error {
	if err := influxDB.connect(); err != nil {
		log.Fatal("InfluxDB Connect failed, err message is ", err)
	}

	// 启动定时任务
	go influxDB.timeTask()
	return nil
}

// timeTask 定时任务
func (influxDB *InfluxDB) timeTask() {
	defer func() {
		if err := recover(); err != nil {
			// alog.Logger.Error("InfluxDB", zap.Any("recover", err))
			log.Println(err)
		}
	}()

	jobpool := make([]*system.Metric, 0, influxDB.DataLen*2)
	TimeC := time.Tick(time.Duration(influxDB.TimeTask) * time.Millisecond)
	for {
		select {
		case <-influxDB.stopc:
			return
		case job, ok := <-influxDB.taskC:
			if ok {
				jobpool = append(jobpool, job...)
				if len(jobpool) >= influxDB.DataLen {
					err := influxDB.writer(jobpool)
					if err != nil {
						// alog.Logger.Error("InfluxDB", zap.Error(err))
						log.Println(err)
					}
					jobpool = jobpool[:0]
				}
			}
			break
		case <-TimeC:
			if len(jobpool) > 0 {
				err := influxDB.writer(jobpool)
				if err != nil {
					// alog.Logger.Error("InfluxDB", zap.Error(err))
					log.Println(err)
				}
				jobpool = jobpool[:0]
			}
			break
		}
	}
}

func (influxDB *InfluxDB) writer(metrics []*system.Metric) error {
	if len(influxDB.conns) == 0 {
		err := influxDB.connect()
		if err != nil {
			return err
		}
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:         influxDB.DataBase,
		RetentionPolicy:  influxDB.RetentionPolicy,
		WriteConsistency: influxDB.WriteConsistency,
	})
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		pt, err := client.NewPoint(metric.Name, metric.Tags, metric.Fields, time.Unix(metric.Time, 0))
		if err != nil {
			// alog.Logger.Error("InfluxDB Write", zap.Error(err))
			log.Println(err)
			return err
		}
		// alog.Logger.Debug("InfluxDB Write", zap.Any("@metric", metric))
		bp.AddPoint(pt)
	}

	// This will get set to nil if a successful write occurs
	err = errors.New("Could not write to any InfluxDB server in cluster")

	p := rand.Perm(len(influxDB.conns))
	for _, n := range p {
		if e := influxDB.conns[n].Write(bp); e != nil {
			log.Println("InfluxDB Write", e)
			// If the database was not found, try to recreate it
			if strings.Contains(e.Error(), "database not found") {
				if errc := createDatabase(influxDB.conns[n], influxDB.DataBase); errc != nil {
					log.Println("ERROR: Database ", influxDB.DataBase, " not found and failed to recreate\n", errc)
				}
			}
		} else {
			err = nil
			break
		}
	}

	return err
}

func (influxDB *InfluxDB) connect() error {
	influxDB.stopc = make(chan bool, 1)
	influxDB.taskC = make(chan []*system.Metric, influxDB.DataLen)

	var urls []string
	for _, u := range influxDB.URLs {
		urls = append(urls, u)
	}

	// Backward-compatability with single Influx URL config files
	// This could eventually be removed in favor of specifying the urls as a list
	if influxDB.URL != "" {
		urls = append(urls, influxDB.URL)
	}

	var conns []client.Client
	for _, u := range urls {
		switch {
		case strings.HasPrefix(u, "udp"):
			parsed_url, err := url.Parse(u)
			if err != nil {
				return err
			}

			if influxDB.UDPPayload == 0 {
				influxDB.UDPPayload = client.UDPPayloadSize
			}
			c, err := client.NewUDPClient(client.UDPConfig{
				Addr:        parsed_url.Host,
				PayloadSize: influxDB.UDPPayload,
			})
			if err != nil {
				return err
			}
			conns = append(conns, c)
		default:
			// If URL doesn't start with "udp", assume HTTP client
			c, err := client.NewHTTPClient(client.HTTPConfig{
				Addr:      u,
				Username:  influxDB.Username,
				Password:  influxDB.Password,
				UserAgent: influxDB.UserAgent,
				// Timeout:   time.Duration((influxDB.Timeout)), //influxDB.Timeout.Duration,
				Timeout: time.Duration(time.Duration(influxDB.Timeout) * time.Second), //influxDB.Timeout.Duration,
			})
			if err != nil {
				return err
			}

			err = createDatabase(c, influxDB.DataBase)
			if err != nil {
				log.Println("Database creation failed: " + err.Error())
				continue
			}

			conns = append(conns, c)
		}
	}

	influxDB.conns = conns
	rand.Seed(time.Now().UnixNano())

	return nil
}

// Close any connections to the Output
func (influxDB *InfluxDB) Close() error {
	var errS string
	for j, _ := range influxDB.conns {
		if err := influxDB.conns[j].Close(); err != nil {
			errS += err.Error()
		}
	}
	if errS != "" {
		return fmt.Errorf("output influxdb close failed: %s", errS)
	}
	close(influxDB.taskC)
	close(influxDB.stopc)
	return nil
}

// Writer  put data
func (influxDB *InfluxDB) Writer(data interface{}) error {
	switch metrics := data.(type) {
	case []*system.Metric:
		influxDB.taskC <- metrics
	}
	return nil
}

// submit 提交任务
func (influxDB *InfluxDB) submit(data []*system.Metric) error {
	influxDB.taskC <- data
	return nil
}

// Description returns a one-sentence description on the Output
func (influxDB *InfluxDB) Description() string {
	return "influxDB"
}

func createDatabase(c client.Client, database string) error {
	// Create Database if it doesn't exist
	_, err := c.Query(client.Query{
		Command: fmt.Sprintf("CREATE DATABASE \"%s\"", database),
	})
	return err
}

func init() {
	influxdb := &InfluxDB{}
	agent.AddReporter("influxdb", influxdb)
}
