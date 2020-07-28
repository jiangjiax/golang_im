package models

import "time"

// 设备
type Device struct {
	AppId          int64     `json:"app_id"`
	UserId         int64     `json:"user_id"`
	Type           int32     `json:"type"`           // 设备类型：1 Android 2 IOS 3 Windows  4 MacOS 5 Web
	Brand          string    `json:"brand"`          // 手机厂商
	Model          string    `json:"model"`          // 机型
	SystemVersion  string    `json:"system_version"` // 系统版本
	Status         int32     `json:"status"`         // 在线状态：0 离线 1 在线
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
	Del            int32     `json:"del"`            // 1 正常 0 删除
	Identification string    `json:"identification"` // 唯一标识
	Id             int64     `json:"id"`
}
