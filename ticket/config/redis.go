package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConn struct {
	Conn *redis.Client
}

var (
	DB *RedisConn
)


func ConfigDB() *redis.Options {
	REDIS_URL := fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	log.Printf("Redis URL: %v\n", REDIS_URL)

	redisConfig := &redis.Options {
		Addr: REDIS_URL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
		PoolSize: 10,
		MinIdleConns: 2,
	}

	return redisConfig
}

func (db *RedisConn) ConnectDB() error {
	redisConfig := ConfigDB()
	rdb := redis.NewClient(redisConfig)

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(context).Result()
	if err != nil {
		log.Printf("Could not connect to Redis: %v", err)
		return err
	}

	log.Println("Connected to Redis")
	db.Conn = rdb
	return nil
}

func (db *RedisConn) Set(key, value string) error {
	err := db.Conn.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %v to value %v: %v", key, value, err)
		return err
	}
	return nil
}

func (db *RedisConn) Get(key string) (string, error) {
	val, err := db.Conn.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("Failed to get key %v: %v", key, err)
		return "", err
	}
	return val, nil
}