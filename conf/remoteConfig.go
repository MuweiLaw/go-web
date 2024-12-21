package conf

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
	"log"
)

type RemoteConf struct {
	Mysql Mysql `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
}

type Mysql struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       Db     `yaml:"db"`
}
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	PoolSize int    `yaml:"pool-size"`
}

type Db struct {
	Name string `yaml:"name"`
}

// GetRemoteConf 根据环境变量读取本地yml配置
func GetRemoteConf() *RemoteConf {
	//获取本地配置
	localConf := getLocalConf()

	serverConf := []constant.ServerConfig{
		{
			IpAddr: localConf.Nacos.Config.Host, //nacos 地址
			Port:   localConf.Nacos.Config.Port, //nacos 端口
		},
	}

	clientConf := &constant.ClientConfig{
		NamespaceId:         localConf.Nacos.Config.Namespace, //命名空间 比较重要 拿取刚才创建的命名空间ID
		TimeoutMs:           3000,
		NotLoadCacheAtStart: true,
		LogDir:              localConf.Nacos.Config.LogDir,
		CacheDir:            localConf.Nacos.Config.CacheDir,
		LogLevel:            localConf.Nacos.Config.LogLevel,
	}

	confClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  clientConf,
			ServerConfigs: serverConf,
		},
	)
	if err != nil {
		log.Fatalf("创建Nacos客户端失败! err:%v", err)
	}

	nConfStr, err := confClient.GetConfig(vo.ConfigParam{
		DataId: localConf.Nacos.Config.DataId, //配置文件名
		Group:  localConf.Nacos.Config.Group,  //配置文件分组
	})
	if err != nil {
		log.Fatalf("读取nacos配置失败! err:%v", err)
	}

	// 解析配置文件
	var rConf *RemoteConf
	if err := yaml.Unmarshal([]byte(nConfStr), &rConf); err != nil {
		log.Fatalf("序列化远程配置失败! err:%v", err.Error())
	}
	log.Printf("读取nacos配置成功! \nnacos-conf:%v", rConf)

	// 监听配置文件变化
	if err := confClient.ListenConfig(vo.ConfigParam{
		DataId: localConf.Nacos.Config.DataId, //配置文件名
		Group:  localConf.Nacos.Config.Group,
		OnChange: func(namespace, group, dataId, data string) {
			log.Printf("监听Nacos配置，DataID:%s，Group:%s，Namespace:%s", dataId, group, namespace)
			// 解析配置文件
			var nrConf *RemoteConf
			if err := yaml.Unmarshal([]byte(data), &nrConf); err != nil {
				log.Printf("监听Nacos中配置异常! \nerr:%v", err.Error())
				return
			}
			//这里我们打印一下解析后的配置文件内容
			log.Printf("配置详情：\n%v", data)
		},
	}); err != nil {
		log.Printf("打开监听Nacos配置异常! \nerr:%v", err.Error())
	}
	return rConf
}
