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

// Group 群组
type Group struct {
	Id           int64     `json:"id"`
	AppId        int64     `json:"app_id"`   // appId
	GroupId      int64     `json:"group_id"` // 群组id
	UserId       int64     `json:"user_id"`
	DeviceId     int64     `json:"device_id"`
	Name         string    `json:"name"`         // 组名
	Introduction string    `json:"introduction"` // 群简介
	UserNum      int32     `json:"user_num"`     // 群组人数
	Type         int32     `json:"type"`         // 群组类型
	Privacy      int32     `json:"privacy"`      // 1 公开群 2 隐私群
	Avatar       string    `json:"avatar"`       // 头像
	Extra        string    `json:"extra"`        // 附加属性
	CreateTime   time.Time `json:"create_time"`  // 创建时间
	UpdateTime   time.Time `json:"update_time"`  // 更新时间
	Way          int32     `json:"way"`
	Coordinatex  float64   `json:"coordinatex"`
	Coordinatey  float64   `json:"coordinatey"`
	Commandword  string    `json:"commandword"`
	UserType     int32     `json:"user_type"` // 1 群成员 2 群管理 3 群主
	Label        string    `json:"label"`
	Ban          int32     `json:"ban"` // 群聊邀请确认 1 不接受邀请请求
}
