package agent

import (
	"encoding/json"
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger ....
var Logger *zap.Logger

// LogInit ...
func LogInit(lp string, lv string, isDebug bool) {
	var js string
	if isDebug {
		js = fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stdout"]
		}`, lv)
	} else {
		js = fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["%s"],
		"errorOutputPaths": ["%s"]
		}`, lv, lp, lp)
	}

	var cfg zap.Config
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
	}
}

// InitLogger  init log
func InitLogger(logpath, loglevel string, IsDebug bool) {
	isDebug := true
	if IsDebug != true {
		isDebug = false
	}
	LogInit(logpath, loglevel, isDebug)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
}
