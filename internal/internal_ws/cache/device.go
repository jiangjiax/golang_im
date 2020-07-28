package cache

import (
	"golang_im/pkg/db"
	"golang_im/pkg/log"
	"golang_im/pkg/models"
	"golang_im/pkg/util"
	"strconv"
	"time"
)

const (
	ListOnlineKey = "ListOnline:"
	deviceExpire  = 2 * time.Hour
)

type deviceCache struct{}

var DeviceCache = new(deviceCache)

func (c *deviceCache) Key(appId, userId int64) string {
	return ListOnlineKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)
}

// Get 获取缓存
func (c *deviceCache) Get(appId, userId int64) ([]models.Device, error) {
	var Devices []models.Device
	err := util.RGet(ListOnlineKey+strconv.FormatInt(appId, 10)+":"+strconv.FormatInt(userId, 10), Devices)
	if err != nil {
		return nil, err
	}
	return Devices, nil
}

// Set 设置缓存
func (c *deviceCache) Set(appId, userId int64, devices []models.Device) error {
	err := util.RSet(ListOnlineKey+strconv.FormatInt(appId, 10)+":"+strconv.FormatInt(userId, 10), devices, deviceExpire)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Del 删除缓存
func (c *deviceCache) Del(appId, userId int64) error {
	_, err := db.RedisCli.Del(ListOnlineKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
