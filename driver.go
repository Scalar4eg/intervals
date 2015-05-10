package main

import (
	"github.com/biogo/store/interval"
	"github.com/vaughan0/go-ini"
	"github.com/bo0rsh201/intervals/drivers"
	"errors"
	"fmt"
)

var supported map[string]Driver

func init() {
	supported = make(map[string]Driver)
	supported["sqlite"] = new(drivers.Sqlite)
}

func InitDriver(general GeneralConfig, file ini.File) (driver Driver, err error) {
	driver, ok := supported[general.Driver]
	if !ok {
		return driver, errors.New(fmt.Sprintf("Unsupported driver %s", general.Driver))
	}
	err = driver.Init(file.Section(general.Driver))
	return driver, err
}

type Driver interface {
	Init(section ini.Section) error
	Load() (interval.IntTree, error)
}