package service

import (
	"log"
	"sync"

	"io.xiu/listAllFiles/domian"
	"io.xiu/listAllFiles/instance"
)

func Insert(files []domian.File, lock *sync.RWMutex) {
	lock.Lock()
	defer lock.Unlock()
	engine := instance.Sqlite3()
	_, err := engine.Insert(&files)
	if err != nil {
		log.Fatal(err.Error())
	}
}
