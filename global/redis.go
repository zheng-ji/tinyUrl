package global

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var GRedisPool *redis.Pool

func NewRedisPool(cfg RedisConf) *redis.Pool {
	GRedisPool = &redis.Pool{
		//MaxIdle:     cfg.MaxIdle,
		IdleTimeout: 1 * time.Hour,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.Addr)
			if nil != err {
				Runlogger.Errorf("[Redis] Dial erro: %v", err)
				return nil, err
			}
			if _, err := c.Do("SELECT", cfg.DB); nil != err {
				c.Close()
				Runlogger.Errorf("[Redis] Select DB error: %v", err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return GRedisPool
}

// TODO 未来需要支持多个缓存实例
