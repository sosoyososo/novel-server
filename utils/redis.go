package utils

import (
	"time"

	"github.com/go-redis/redis"
)

var (
	redisConf RedisConf
	RedisC    *redis.Client
)

type RedisConf struct {
	Host string `json:"host"`
	Pswd string `json:"pswd"`
}

func newRedis(hostAddr, pswd string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     hostAddr,
		Password: pswd, // no password set
		DB:       0,    // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

/**
 * 判断一个字符串值在redis里是否存在
 */
func RedisListJudgeExist(
	key, value string,
	dataLoader func() ([]string, error),
	expire time.Duration,
	resultAutoSaveIfNotExist bool,
) (bool, error) {
	redisK := key
	c := RedisC
	v, err := c.Exists(redisK).Result()
	if nil != err {
		return false, err
	}

	if v == 0 {
		list, err := dataLoader()
		if nil != err {
			return false, err
		}
		if len(list) > 0 {
			err = c.SAdd(redisK, list).Err()

			if nil != err {
				return false, err
			}
		}
	}
	ret, err := c.SIsMember(redisK, value).Result()

	if !ret && resultAutoSaveIfNotExist {
		c.SAdd(redisK, value)
	}
	if expire > 0 {
		c.PExpire(redisK, expire)
	}

	return ret, err
}
