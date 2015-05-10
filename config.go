package main

import (
	"github.com/vaughan0/go-ini"
	"errors"
	"strconv"
)

type GeneralConfig struct {
	Host string
	Port int
	Driver string
}

func NewGeneralConfig(file ini.File) (conf GeneralConfig, err error) {
	conf = GeneralConfig{}
	section := file.Section("general")

	// @TODO remove copy-paste
	if val, ok := section["host"]; ok {
		conf.Host = val
	} else {
		return conf, errors.New("empty 'host' field in general config")
	}

	if val, ok := section["port"]; ok {
		conf.Port, err = strconv.Atoi(val)
		if err != nil {
			return conf, err
		}
	} else {
		return conf, errors.New("empty 'port' field in general config")
	}

	if val, ok := section["driver"]; ok {
		conf.Driver = val
	} else {
		return conf, errors.New("empty 'driver' field in general config")
	}
	return conf, nil
}
