package internal_ws

import "sync"

var manager sync.Map

// store 存储
func Store(deviceId int64, ctx *WSConn) {
	manager.Store(deviceId, ctx)
}

// load 获取
func Load(deviceId int64) *WSConn {
	value, ok := manager.Load(deviceId)
	if ok {
		return value.(*WSConn)
	}
	return nil
}

// delete 删除
func Delete(deviceId int64) {
	manager.Delete(deviceId)
}
