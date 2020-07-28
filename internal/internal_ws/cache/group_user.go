package cache

import (
	"golang_im/pkg/db"
	"golang_im/pkg/log"
	"golang_im/pkg/pb"
	"golang_im/pkg/util"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	GroupUserKey    = "GroupUser:"
	GroupUserExpire = 2 * time.Hour
)

type groupUserCache struct{}

var GroupUserCache = new(groupUserCache)

func (c *groupUserCache) Key(appId, groupId int64) string {
	return ListOnlineKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(groupId, 10)
}

// Get 获取缓存
func (c *groupUserCache) Get(appId, groupId int64) (*pb.GroupUserResp, error) {
	var GroupUser *pb.GroupUserResp
	err := util.RGet(ListOnlineKey+strconv.FormatInt(appId, 10)+":"+strconv.FormatInt(groupId, 10), GroupUser)
	if err != nil && err != redis.Nil {
		log.Error(err)
		return nil, err
	}
	if err == redis.Nil {
		return nil, nil
	}
	return GroupUser, nil
}

// Set 设置缓存
func (c *groupUserCache) Set(appId, groupId int64, GroupUser *pb.GroupUserResp) error {
	err := util.RSet(ListOnlineKey+strconv.FormatInt(appId, 10)+":"+strconv.FormatInt(groupId, 10), GroupUser, GroupUserExpire)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Del 删除缓存
func (c *groupUserCache) Del(appId, userId int64) error {
	_, err := db.RedisCli.Del(ListOnlineKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
