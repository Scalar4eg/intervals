package drivers

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/vaughan0/go-ini"
	"errors"
	"github.com/biogo/store/interval"
	"database/sql"
	"fmt"
	"github.com/bo0rsh201/intervals/common"
	"sync"
	"log"
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

func (driver Sqlite) loadChunk(lock *sync.Mutex, t *interval.IntTree, start int ,limit int) error {
	db, err := sql.Open("sqlite3", driver.file)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", driver.table, limit, start))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		iv := common.IntInterval{}
		var id int
		err := rows.Scan(&id, &iv.Start, &iv.End)
		if err != nil {
			return err
		}
		iv.Id = uintptr(id)
		lock.Lock()
		err = t.Insert(iv, false)
		lock.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (driver Sqlite) Load() (tree interval.IntTree, err error) {
	var count int
	db, err := sql.Open("sqlite3", driver.file)
	if err != nil {
		return tree, err
	}
	defer db.Close()
	row := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", driver.table))
	if err != nil {
		return tree, err
	}

	err = row.Scan(&count)
	if err != nil {
		return tree, err
	}

	lock := &sync.Mutex{}

	workerCount := int(count / 10)

	resChan := make(chan error)
	threadCount := 0
	for i:=0; i < count; i = i + workerCount {
		threadCount++
		go func(start int) {
			log.Printf("start %d:%d threads=%d", start, workerCount, threadCount)
			resChan <- driver.loadChunk(lock, &tree, start, workerCount)
		}(i)
	}
	for threadCount > 0 {
		err := <-resChan
		threadCount--
		log.Printf("done threads=%d", threadCount)
		if err != nil {
			return tree, err
		}
	}

	return tree, nil
}