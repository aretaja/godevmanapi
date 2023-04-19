package config

import (
	"os"

	"github.com/tkanos/gonfig"
)

// API configuration sruct
type Configuration struct {
	DbURL     string `env:"GODEVMANAPI_DBURL"`
	ApiListen string `env:"GODEVMANAPI_LISTEN"`
	Salt      string `env:"GODEVMANAPI_SALT"`
}

// Fills Configuration struct. Prefers environment variables
func GetConfig() (*Configuration, error) {
	conf := new(Configuration)

	f := "/usr/local/etc/godevmanapi.conf"
	if os.Getenv("GODEVMAN_TESTDB") != "" {
		f = "/usr/local/etc/godevmanapi_testdb.conf"
	}

	err := gonfig.GetConf(f, conf)

	return conf, err
}
