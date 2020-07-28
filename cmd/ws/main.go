package main

import (
	"fmt"
	"golang_im/api/api_ws"
	"golang_im/config"
	"golang_im/internal/internal_ws"
	"golang_im/pkg/db"
	"golang_im/pkg/util"
	"time"
)

func main() {
	// 生成测试用户token
	now := time.Now()
	t := now.Add(time.Hour * 24)
	value, _ := util.GetToken(1, 1, 1, t.Unix(), util.PublicKey)
	fmt.Println("测试用户1：", value)
	value, _ = util.GetToken(1, 2, 2, t.Unix(), util.PublicKey)
	fmt.Println("测试用户2：", value)

	// 启动数据库连接
	db.Mysql_init()
	db.Redis_init()

	// 启动Controller
	api_ws.Controller_init()

	// 启动websocket服务器
	internal_ws.StartWSServer(config.WSConf.WSListenAddr)
}
