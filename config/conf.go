package config

import (
	"os"
)

var (
	DBConf   dbConf
	WSConf   wsConf
	GRPCConf grpcConf
)

// 数据库配置
type dbConf struct {
	MySQL    string
	RedisIP  string
	RedisPwd string
}

// WS配置
type wsConf struct {
	WSListenAddr string
}

// GRPC配置
type grpcConf struct {
	GRPCListenAddr string
}

func init() {
	// 获取环境变量
	env := os.Getenv("golang_im_env")
	switch env {
	case "dev":
		initDevConf() // 开发环境
	case "prod":
		initProdConf() // 生产环境
	case "test":
		initTestConf() // 测试环境
	default:
		initDevConf()
	}
}
