package toolbox

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"mini_cache/cache/provider"
	"mini_cache/config"
	"os"
)

func LoadConfig() *config.Configs {
	dataBytes, err := os.ReadFile("./conf/config.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return nil
	}
	fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	conf := config.Configs{}
	err = yaml.Unmarshal(dataBytes, &conf)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return nil
	}
	fmt.Printf("config → %+v\n", conf) // config → {Mysql:{Url:127.0.0.1 Port:3306} Redis:{Host:127.0.0.1 Port:6379}}
	return &conf
}

func LoadOpts() provider.CacheOptions {
	dataBytes, err := os.ReadFile("./conf/config.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return provider.CacheOptions{}
	}
	fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	opts := provider.CacheOptions{}
	err = yaml.Unmarshal(dataBytes, &opts)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return provider.CacheOptions{}
	}
	fmt.Printf("opts → %+v\n", opts) // config → {Mysql:{Url:127.0.0.1 Port:3306} Redis:{Host:127.0.0.1 Port:6379}}
	return opts
}
