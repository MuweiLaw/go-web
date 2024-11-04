package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:admin123456@tcp(localhost:3306)/auth?charset=utf8mb4&parseTime=True")
	if err := db.Ping(); err != nil {
		fmt.Println("MySql连接失败")
		panic("启动失败")
	}
	fmt.Println("MySql连接成功")

	// 创建redis连接池
	client := redis.NewClient(&redis.Options{
		Addr:     "139.159.191.200:6379",
		Password: "admin123", // 设置密码
		DB:       0,          // 选择数据库
		PoolSize: 10,         // 设置连接池大小
	})

	// 测试连接
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	// 设置值
	err = client.Set("key", "value", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取值
	val, err := client.Get("key").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("key:", val)
	}

	// 关闭连接
	defer client.Close()
	// 1.创建路由
	r := gin.Default()

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	err = r.Run(":8000")
	if err != nil {
		return
	}
}
