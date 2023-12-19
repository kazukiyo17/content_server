package conf

import (
	"content_server/setting"
	"fmt"
	"github.com/gomodule/redigo/redis"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

var (
	once sync.Once
	C    = new(Config) // C 全局配置
)

type Config struct {
	//Server       *serverConfig                    // 服务配置
	MysqlConn *gorm.DB    // 数据库实例
	RedisConn *redis.Pool // redis实例
	NlpConn   *nlp.Client
	//RedisMQConn  *redis.Pool
	//RedisMQConn *redis_mq.RedisStreamMQClient
	//RedisConnMap   map[string]redis.Pool // redis实例
}

// Load 加载服务配置,加载失败直接退出
func Load() {
	once.Do(func() {
		cfg := new(Config)
		err := cfg.loadMysql()
		if err != nil {
			//logger.Errorf(errcode.DBFailError, "loadMysql err: %v", err)
			panic("failed to loadMysql")
		}

		err = cfg.loadRedis()
		if err != nil {
			//logger.Errorf(errcode.RedisFailError, "loadRedis err: %v", err)
			panic("failed to loadRedis")
		}

		//err = cfg.loadRedisMQ()
		//if err != nil {
		//	//logger.Errorf(errcode.RedisFailError, "loadRedis err: %v", err)
		//	panic("failed to loadRedis")
		//}

		C = cfg
	})
}

// loadMysql 初始化MySQL连接池
func (c *Config) loadMysql() error {
	_dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User, setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host, setting.DatabaseSetting.Port,
		setting.DatabaseSetting.Name)
	db, err := gorm.Open(mysql.Open(_dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	c.MysqlConn = db
	return nil
}

func (c *Config) loadRedis() error {
	c.RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

//func (c *Config) loadRedisMQ() error {
//	redisConn := &redis.Pool{
//		MaxIdle:     setting.RedisSetting.MaxIdle,
//		MaxActive:   setting.RedisSetting.MaxActive,
//		IdleTimeout: setting.RedisSetting.IdleTimeout,
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
//			if err != nil {
//				return nil, err
//			}
//			if setting.RedisSetting.Password != "" {
//				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
//					c.Close()
//					return nil, err
//				}
//			}
//			return c, err
//		},
//	}
//	c.RedisMQConn = &redis_mq.RedisStreamMQClient{
//		ConnPool: redisConn,
//	}
//	return nil
//}
