package common

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"math"
)


func Round(f float64) float64 {
    if f < 0 {
        return math.Ceil(f - .5)
    }
    return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) float64 {
    shift := math.Pow(10, float64(places))
    return Round(f * shift) / shift;
}

// defined error
type commonError struct {
	err error
}

func CatchError(err *error) {
	if e := recover(); e != nil {
		ce, ok := e.(error)
		if ok{
			*err = ce
		}else{
			*err = nil
		}
	}
	return
}

func catchError(err *error) {
	if e := recover(); e != nil {
		ce, ok := e.(commonError)
		if !ok {
			fmt.Println("warning, re-panic...")
			panic(e)//if not defined error, then re-panic, sunch as a runtime error
		}
		*err = ce.err
	}
	return
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Current() string {
	return time.Now().Format("2006-01-02 15:04:05.99999")
}

func NewRedisPool(config *Config) *redis.Pool {
	var timeout  int64= int64(config.Redis_IdleTimeout)
	return &redis.Pool{
		MaxIdle: config.Redis_MaxIdle,
		MaxActive: config.Redis_MaxActive,
		IdleTimeout: time.Duration( timeout * int64(time.Second)),
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis_Addr)
			if err != nil {
				//return nil, err
				panic(err)
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil{
				//fmt.Printf("TestOnBorrow-Ping error: %s\n", err )
				panic(err)
			}
			return err
		},
		Wait:true,
	}
}

func InitRedis(config *Config) redis.Conn {
	conn, err := redis.Dial("tcp", config.Redis_Addr)
	if err != nil {
		panic(err.Error())
	}
	return conn
}
