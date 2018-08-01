package common

import (
	"github.com/Terry-Mao/goconf"
)


type Config struct {
	GO_MAX_PROCS               int `goconf:"golang:go-max-procs"`
	THREAD_POOL                int `goconf:"golang:thread-pool"`
	PIPELINE_JOB               int `goconf:"golang:pipeline-job"`
	FILE_READBUFFER            int `goconf:"golang:file-readbuffer"`

	//redis配置
	Redis_Addr                 string `goconf:"db:redis-addr"`
	Redis_MaxIdle              int `goconf:"db:redis-maxidle"`
	Redis_MaxActive            int `goconf:"db:redis-maxactive"`
	Redis_IdleTimeout          int `goconf:"db:redis-idletimeout"`
}

var myConfig *Config


func InitConfig(configFile string) *Config {
	conf := goconf.New()

	if err := conf.Parse(configFile); err !=nil {
		panic(err)
	}

	myConfig := &Config{}

	if err := conf.Unmarshal(myConfig); err!=nil {
		panic(err)
	}

	return myConfig
}