package drivers

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/vaughan0/go-ini"
	"errors"
	"github.com/biogo/store/interval"
	"database/sql"
	"fmt"
	"github.com/bo0rsh201/intervals/common"
)

type Sqlite struct {
	file string
	table string
}

func (p *Sqlite) Init(section ini.Section) error {
	driver := *p
	// @TODO remove copy-paste
	if val, ok := section["file"]; ok {
		driver.file = val
	} else {
		return errors.New("Empty 'file' field in config")
	}

	if val, ok := section["table"]; ok {
		driver.table = val
	} else {
		return errors.New("Empty 'table' field in config")
	}
	*p = driver
	return nil
}

func (driver Sqlite) Load() (tree interval.IntTree, err error) {
	db, err := sql.Open("sqlite3", driver.file)
	if err != nil {
		return tree, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", driver.table))
	if err != nil {
		return tree, err
	}
	defer rows.Close()

	for rows.Next() {
		iv := common.IntInterval{}
		var id int
		err := rows.Scan(&id, &iv.Start, &iv.End)
		if err != nil {
			return tree, err
		}
		iv.Id = uintptr(id)
		err = tree.Insert(iv, false)
		if err != nil {
			return tree, err
		}
	}
	return tree, nil
}