package instance

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"io.xiu/listAllFiles/domian"
	"xorm.io/xorm"
)

func Sqlite3() *xorm.Engine {
	engine, err := xorm.NewEngine("sqlite3", "./file.db")
	if err != nil {
		log.Fatalln(err)
	}
	engine.Sync2(new(domian.File))
	return engine
}
