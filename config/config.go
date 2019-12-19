package config

import (
	"errors"
	"flag"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"go2region/utils/logs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

const SERVER_NAME = "go2region"

var CONFIG_KEY = fmt.Sprintf("/configs/eebo.ehr.%s/system", SERVER_NAME)

var Config *config           // 静态配置
var DynamicConf *dynamicConf // 动态配置
var ConsulAddr string
var _path string
var _client *consulapi.Client

// 静态配置，程序启动后无法再做更改的参数配置
type config struct {
	BaseConf BaseConf `yaml:"base"`
}
type BaseConf struct {
	// 当前服务监听的端口
	ServerPort int `yaml:"server_port"`

	// 访问日志和运行日志相关配置
	AccessLogPath string `yaml:"access_log_path"`
	SystemLogPath string `yaml:"error_log_path"`
	LogMaxAge     int    `yaml:"log_max_age"`

	IpdataPath string `yaml:"ipdata_path"`
}

// 动态配置，程序运行过程中，可以动态更改的参数配置
type dynamicConf struct {
	LogLevel string `yaml:"log_level"`
}

// 初始化解析参数
func init() {
	flag.StringVar(&_path, "c", SERVER_NAME+".yml", "default config path")
	flag.StringVar(&ConsulAddr, "consul", os.Getenv("CONSUL"), "default consul address")
}

// 优先从consul中加载配置，没有则从配置文件中加载配置
// consul中的配置文件需为yaml格式
func InitConfig() error {
	var err error
	var content []byte

	if ConsulAddr != "" {
		content, err = fetchConfig(CONFIG_KEY, watchDynamicConfig)
	} else {
		content, err = ioutil.ReadFile(_path)
	}

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return errors.New("not found nothing config")
	}

	Config = &config{}
	if err := yaml.Unmarshal(content, Config); err != nil {
		return err
	}

	DynamicConf = &dynamicConf{}
	if err := yaml.Unmarshal(content, DynamicConf); err != nil {
		return err
	}

	level, err := logrus.ParseLevel(DynamicConf.LogLevel)
	if err == nil {
		logs.LogSystem.SetLevel(level)
	}

	fmt.Printf("static config => [%#v]\n", Config)

	return nil
}

// 从consul中获取配置信息
func fetchConfig(configKey string, watchFn func([]byte)) ([]byte, error) {
	config := consulapi.DefaultConfig()
	config.Address = ConsulAddr
	_client, err := consulapi.NewClient(config)
	if err != nil {
		logs.LogSystem.Error("consul client error : ", err)
	}
	data, meta, err := _client.KV().Get(configKey, nil)

	if watchFn != nil {
		go func() {
			for {
				options := &consulapi.QueryOptions{WaitIndex: meta.LastIndex, WaitTime: time.Minute * 5}
				data, meta, err = _client.KV().Get(configKey, options)
				if err == nil {
					watchFn(data.Value)
				} else {
					for {
						_client, err = consulapi.NewClient(config)
						if err != nil {
							logs.LogSystem.Error("consul client error : ", err)
							time.Sleep(time.Second * 10)
						} else {
							data, meta, err = _client.KV().Get(configKey, nil)
							if err == nil {
								break
							}
						}
					}
				}
			}
		}()
	}
	return data.Value, err
}

// 监控动态配置，并使用值拷贝进行全部替换
func watchDynamicConfig(val []byte) {
	dc := new(dynamicConf)
	*dc = *DynamicConf

	_ = yaml.Unmarshal(val, dc)

	DynamicConf = dc

	// 更新运行日志等级
	level, err := logrus.ParseLevel(DynamicConf.LogLevel)
	if err == nil {
		logs.LogSystem.SetLevel(level)
	}

	logs.LogSystem.Error("test log level, curr_level: ", logs.LogSystem.GetLevel())

	fmt.Printf("Latest dynamic config => [%#v]\n", DynamicConf)
}
