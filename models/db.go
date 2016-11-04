package models

import (
	"github.com/tidwall/buntdb"
)

// 定义全局的DB变量
var (
	HistoryDB *buntdb.DB
)

// InitHistoryDB 初始化下载历史DB
func InitHistoryDB(filename string) {
	db, err := buntdb.Open(filename)
	if err != nil {
		panic(err)
	}
	err = db.CreateIndex("time", "*", buntdb.IndexJSON("Time"))
	if err != nil {
		panic(err)
	}
	HistoryDB = db
}
