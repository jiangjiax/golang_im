package util

import (
	"golang_im/pkg/db"
	"golang_im/pkg/log"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Set 将指定值设置到redis中，使用json的序列化方式
func RSet(key string, value interface{}, duration time.Duration) error {
	bytes, err := jsoniter.Marshal(value)
	if err != nil {
		log.Error(err)
		return err
	}

	err = db.RedisCli.Set(key, bytes, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

// get 从redis中读取指定值，使用json的反序列化方式
func RGet(key string, value interface{}) error {
	bytes, err := db.RedisCli.Get(key).Bytes()
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(bytes, value)
	if err != nil {
		return err
	}
	return nil
}
