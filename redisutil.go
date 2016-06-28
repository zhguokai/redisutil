package redisutil

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
	"strconv"
	"errors"
)

//Redis连接池
var pool *redis.Pool

type RedisConfig struct {
	Host      string
	Port      int
	MaxActive int
	MaxIdle   int
	TimeOut   time.Duration
	PassWord  string
}

//注册Redis连接
func RegisterRedis(rc *RedisConfig) {

	if rc == nil {
		rc = &RedisConfig{}
		rc.Host = "127.0.0.1"
		rc.Port = 6379
		rc.MaxActive = 100
		rc.MaxIdle = 100
		rc.TimeOut = 30 * time.Minute
	}

	if rc.Host == "" {
		rc.Host = "127.0.0.1"
	}
	if rc.Port == 0 {
		rc.Port = 6379
	}
	if rc.MaxActive == 0 {
		rc.MaxActive = 100
	}
	if rc.MaxIdle == 0 {
		rc.MaxIdle = 0
	}
	if rc.TimeOut == 0 {
		rc.TimeOut = 30 * time.Minute
	}

	url := rc.Host + ":" + strconv.Itoa(rc.Port)

	pool = &redis.Pool{
		//最大空闲连接
		MaxIdle: rc.MaxIdle,
		//最大活动连接
		MaxActive: rc.MaxActive,
		//空闲连接过期时间
		IdleTimeout: rc.TimeOut,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}
			if rc.PassWord != "" {
				if _, err := c.Do("AUTH", rc.PassWord); err != nil {
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
}

//选择数据库
func SelectDb(con redis.Conn, db int) bool {
	if con == nil {
		return false
	}
	_, err := con.Do("SELECT", db)
	if err != nil {
		log.Println("choose db failer,msg:", err.Error())
		return false
	}
	return true
}

//获取连接
func GetConn() redis.Conn {
	if pool == nil {
		log.Println("RegisterDB with the default Config value")
		RegisterRedis(nil)
	}
	return pool.Get()
}

func SetKeyValue(db int, key interface{}, value interface{}) bool {
	con := GetConn()
	defer con.Close()
	if !SelectDb(con, db) {
		return false
	}
	_, err := con.Do("SET", key, value)
	if err != nil {
		log.Println("set value failer,msg:", err.Error())
		return false
	}
	return true
}
//获取
func GetValue(db int, key interface{}) (value interface{}, err error) {
	con := GetConn()
	defer con.Close()
	if !SelectDb(con, db) {
		return "", errors.New("choose db failer")
	}
	v, errDo := con.Do("GET", key)
	if errDo != nil {
		return "", errDo
	}
	return v, nil
}

//获取字符串Value
func GetStrValue(db int, key interface{}) (value string, err error) {
	con := GetConn()
	defer con.Close()
	if !SelectDb(con, db) {
		return "", errors.New("choose db failer")
	}
	v, errDo := con.Do("GET", key)
	if errDo != nil {
		return "", errDo
	}
	if v == nil {
		return "", errors.New("the key doesn't exist")
	}

	value, err = redis.String(v, errDo)
	if err != nil {
		return "", err
	}

	return value, nil
}

//存储集合数据
func HSetKeyFieldValue(db int, key, field, value interface{}) bool {
	con := GetConn()
	defer con.Close()
	if !SelectDb(con, db) {
		return false
	}
	_, err := con.Do("HSET", key, field, value)
	if err != nil {
		log.Println("HSET failer,msg:", err.Error())
		return false
	}
	return true
}



//获取哈希数据
func HGetKeyFieldValue(db int, key, field interface{}) (value interface{}, err error) {
	con := GetConn()
	defer con.Close()
	if !SelectDb(con, db) {
		return "", errors.New("choose db failer")
	}
	reply, err := con.Do("HGET", key, field)
	if err != nil {
		log.Fatalln("HGET value failer,msg:", err.Error())
		return "", err
	}
	return reply, nil
}

//获取哈希数据
func HGetKeyFieldStrValue(db int, key, field interface{}) (value string, err error) {
	con := GetConn()
	defer con.Close()
	if !SelectDb(con, db) {
		return "", errors.New("choose db failer")
	}
	reply, err := con.Do("HGET", key, field)
	if err != nil {
		log.Fatalln("HGET value failer,msg:", err.Error())
		return "", err
	}
	if reply == nil {
		return "", errors.New("the key doesn't exist")
	}
	v, err := redis.String(reply, err)
	if err != nil {
		log.Fatalln("HGET value can't convert to string,msg:", err.Error())
		return "", err
	}
	return v, nil
}



//删除键值数据
func DropKey(db int, key string) bool {
	con := GetConn()
	defer con.Close()
	if !SelectDb(con, db) {
		log.Fatalln("Select Redis Db failer")
		return false
	}
	_, err := con.Do("DEL", key)
	if err != nil {
		log.Fatalln("DEL Redis key ", key, " failer,msg:", err.Error())
		return false
	}
	return true
}
