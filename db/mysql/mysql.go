package mysql

import (
	"database/sql"
	"fmt"
	"go-web/conf"
	"log"
)

// GetMysql 获取连接MySQL
func GetMysql(nConf conf.RemoteConf) *sql.DB {
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
