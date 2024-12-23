package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type LocalConf struct {
	Nacos Nacos `yaml:"nacos"`
}

type AppConf struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Nacos struct {
	Discovery Discovery `yaml:"discovery"`
	Config    Config    `yaml:"config"`
}
type Discovery struct {
	Addr      string `yaml:"addr"`
	Namespace string `yaml:"namespace"`
}

type Config struct {
	Host      string `yaml:"host"`
	Port      uint64 `yaml:"port"`
	Namespace string `yaml:"namespace"`
	LogDir    string `yaml:"log-dir"`
	CacheDir  string `yaml:"cache-dir"`
	LogLevel  string `yaml:"log-level"`
	DataId    string `yaml:"data-id"`
	Group     string `yaml:"group"`
}

// GetAppConf 读取app.yml配置
func GetAppConf() AppConf {
	//获取当前项目运行文件路径
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//读取app.yml配置文件
	appYml, err := os.ReadFile(fmt.Sprintf("%s/conf/yml/app.yml", wd))
	if err != nil {
		log.Fatalf("读取app.yml配置文件失败! err: %v", err.Error())
	}
	//将yml配置反序列化对象
	var c AppConf
	err = yaml.Unmarshal(appYml, &c)
	if err != nil {
		log.Fatalf("解码app.yml失败: %v", err)
	}
	log.Printf("解码app.yml成功: %v", c)
	return c
}

// GetLocalConf 根据环境变量读取本地yml配置
func getLocalConf() LocalConf {

	//读取环境变量
	runMode, present := os.LookupEnv("GoWebRunMode")
	if present {
		log.Printf("环境变量 GoWebRunMode 的值是: %s", runMode)
	} else {
		log.Fatalf("环境变量 GoWebRunMode 未找到!")
	}

	//获取当前项目运行文件路径
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//读取本地yml配置文件
	configYml, err := os.ReadFile(fmt.Sprintf("%s/conf/yml/%s.yml", wd, runMode))
	if err != nil {
		log.Fatalf("读取配置文件失败! err: %v", err.Error())
	}

	//将yml配置反序列化对象
	var c LocalConf
	err = yaml.Unmarshal(configYml, &c)
	if err != nil {
		log.Fatalf("解码Yml失败: %v", err)
	}
	log.Printf("解码Yml成功: %v", c)
	return c
}
