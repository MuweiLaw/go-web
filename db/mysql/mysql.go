package mysql

import (
	"database/sql"
	"fmt"
	"go-web/conf"
	"log"
	"sync"
)

var (
	mysql *sql.DB
	once  sync.Once
)

// GetMysql 获取连接MySQL
func GetMysql() *sql.DB {

	if mysql == nil {
		once.Do(func() {
			// 进行单例的初始化工作
			nConf := conf.GetRemoteConf()
			mysql = openMysql(*nConf)
		})
	}
	return mysql
}

// openMysql 获取连接MySQL
func openMysql(nConf conf.RemoteConf) *sql.DB {

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
			nConf.Mysql.User,
			nConf.Mysql.Password,
			nConf.Mysql.Addr,
			nConf.Mysql.Db.Name))
	if err != nil {
		log.Fatalf("打开MySQL失败! err:%v", err)
	}

	if err := db.Ping(); err != nil {
		log.Panicf("MySql连接失败:%v", err)
	}
	log.Printf("MySql连接成功")
	return db
}
