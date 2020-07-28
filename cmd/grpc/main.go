package main

import (
	"golang_im/config"
	"golang_im/internal/internal_grpc"
	"golang_im/pkg/db"
)

func main() {
	// 启动数据库连接
	db.Mysql_init()
	db.Redis_init()

	// 启动grpc服务器
	internal_grpc.StartGRPCServer(config.GRPCConf.GRPCListenAddr)
}
