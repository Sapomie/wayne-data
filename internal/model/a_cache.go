package model

import (
	"encoding/json"
	"github.com/Sapomie/wayne-data/pkg/log"
	"github.com/Sapomie/wayne-data/pkg/setting"
	"github.com/garyburd/redigo/redis"
)

type Cache struct {
	*redis.Pool
}

func NewCacheEngine(s *setting.RedisSettingS) (*redis.Pool, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(s.Network, s.Address+":"+s.Port, redis.DialDatabase(0))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	//have a try
	_, err := pool.Get().Do("ping")
	return pool, err
}

func NewCache(pool *redis.Pool) *Cache {
	return &Cache{pool}
}

func (c *Cache) getConn() redis.Conn {
	return c.Pool.Get()
}

func (c *Cache) GetString(key string) (string, bool, error) {
	conn := c.getConn()
	defer conn.Close()
	val, err := redis.String(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func (c *Cache) Get(key string, v interface{}) (exist bool, err error) {
	val, exist, err := c.GetString(key)
	if err != nil {
		return false, err
	}
	if exist {
		err = json.Unmarshal([]byte(val), v)
		if err != nil {
			return false, err
		}
	}

	return exist, nil
}

func (c *Cache) SetString(key, value string, expire int) error {
	conn := c.getConn()
	defer conn.Close()
	if expire > 0 {
		_, err := conn.Do("SET", key, value, expire)
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("SET", key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Set(key string, value interface{}, expire int) (err error) {
	byt, err := json.Marshal(value)
	if err != nil {
		log.Error("json marshal err", err, "marshal key", key)
	}
	err = c.SetString(key, string(byt), expire)
	if err != nil {
		return err
	}
	return nil
}

func (r *Cache) FlushDb() (err error) {
	conn := r.getConn()
	defer conn.Close()
	if _, err = conn.Do("flushdb"); err != nil {
		return err
	}
	return nil
}
