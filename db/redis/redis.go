package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"go-web/conf"
)

// GetRedis 获取redis客户端
func GetRedis(nConf conf.RemoteConf) *redis.Client {

	// 创建redis连接池
	client := redis.NewClient(&redis.Options{
		Addr:     nConf.Redis.Addr,
		Password: nConf.Redis.Password, // 设置密码
		DB:       nConf.Redis.Db,       // 选择数据库
		PoolSize: nConf.Redis.PoolSize, // 设置连接池大小
	})

	// 测试连接
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}
