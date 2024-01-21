package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var cfg *ini.File

type Server struct {
	RunMode     string
	HttpPort    string
	Domain      string
	AuthExpire  int
	SceneExpire int
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type RedisMQ struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisMQSetting = &RedisMQ{}

type TencentCloud struct {
	SecretID  string
	SecretKey string
}

var TencentCloudSetting = &TencentCloud{}

type DashScope struct {
	Url    string
	ApiKey string
}

var DashScopeSetting = &DashScope{}

func Setup() {
	var err error
	cfg, err = ini.Load("app.ini")
	if err != nil {
		log.Printf("setting.Setup, fail to parse 'app.ini': %v", err)
	}

	//mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("dash-scope", DashScopeSetting)
	mapTo("tencent-cloud", TencentCloudSetting)
	mapTo("redis-mq", RedisMQSetting)

	//ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	//ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo 用于映射配置文件中的各个 section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Printf("Cfg.MapTo %s err: %v", section, err)
	}
}
