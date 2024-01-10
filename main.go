package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:admin123456@tcp(localhost:3306)/auth?charset=utf8mb4&parseTime=True")
	if err := db.Ping(); err != nil {
		fmt.Println("连接失败")
	}
	fmt.Println("连接成功")

	// 1.创建路由
	r := gin.Default()

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	err := r.Run(":8000")
	if err != nil {
		return
	}
}
