package main

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type Conf struct {
	logFile string
	verbose bool
}

func getConf() (Conf, error) {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("etc/rpkgm.yaml")
	if err != nil {
		return Conf{}, err
	}
	confData := config.Data()
	conf := Conf{
		logFile: confData["logFile"].(string),
		verbose: confData["verbose"].(bool),
	}
	return conf, nil
}
