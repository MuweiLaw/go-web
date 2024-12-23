package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	redis2 "github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"go-web/conf"
	"go-web/db/mysql"
	"go-web/db/redis"
	"log"
)

func main() {
	appConf := conf.GetAppConf()
	nConf := conf.GetRemoteConf()

	dbMysql := mysql.GetMysql(*nConf)
	dbRedis := redis.GetRedis(*nConf)

	// 监听端口，为空则随机端口
	if err := gin.Default().Run(appConf.Server.Port); err != nil {
		log.Panicf("gin引擎启动失败! err:%v", err)
	}

	/***************************************************测试***************************************************/
	//测试Mysql
	rows, err := dbMysql.Query("SELECT * FROM `auth`.`auth_user` LIMIT ?,?", 1, 10)
	if err != nil {
		log.Printf("查询失败!")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("rows close err: %s", err.Error())
		}
	}(rows)

	// 查询数据
	var (
		account  string
		password string
	)
	for rows.Next() {
		if err := rows.Scan(&account, &password); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, Name: %s\n", account, password)
	}

	// 检查迭代是否因为错误而提前结束
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	//测试Redis, 设置值
	err = dbRedis.Set("key", "value", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取值
	val, err := dbRedis.Get("key").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("key:", val)
	}

	// 关闭连接
	defer func(dbRedis *redis2.Client) {
		if err := dbRedis.Close(); err != nil {
			log.Printf("rows close err: %s", err.Error())
		}
	}(dbRedis)
}
