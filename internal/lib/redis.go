package lib

import (  
	"os"  
	"strconv"  
	"time"  

	"github.com/redis/go-redis/v9"  
)  

type RedisConfig struct {  
	Addr         string  
	Password     string  
	DB           int  
	PoolSize     int  
	ReadTimeout  time.Duration  
	WriteTimeout time.Duration  
}  

func NewRedisConfig() *RedisConfig {  
	return &RedisConfig{  
		Addr:         getEnv("REDIS_ADDR", "localhost:6379"),  
		Password:     getEnv("REDIS_PASSWORD", ""),  
		DB:           getEnvInt("REDIS_DB", 0),  
		PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),  
		ReadTimeout:  time.Duration(getEnvInt("REDIS_READ_TIMEOUT", 3)) * time.Second,  
		WriteTimeout: time.Duration(getEnvInt("REDIS_WRITE_TIMEOUT", 3)) * time.Second,  
	}  
}  

func (rc *RedisConfig) ToClientOptions() *redis.Options {  
	return &redis.Options{  
		Addr:         rc.Addr,  
		Password:     rc.Password,  
		DB:           rc.DB,  
		PoolSize:     rc.PoolSize,  
		ReadTimeout:  rc.ReadTimeout,  
		WriteTimeout: rc.WriteTimeout,  
	}  
}  

func getEnv(key, defaultValue string) string {  
	if value := os.Getenv(key); value != "" {  
		return value  
	}  
	return defaultValue  
}  

func getEnvInt(key string, defaultValue int) int {  
	valueStr := os.Getenv(key)  
	if valueStr == "" {  
		return defaultValue  
	}  
	value, err := strconv.Atoi(valueStr)  
	if err != nil {  
		return defaultValue  
	}  
	return value  
}  