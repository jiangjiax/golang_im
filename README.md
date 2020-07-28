#### 项目需求

使用 Golang 完成一个即时通讯服务器，支持功能：
1. 支持前端 websocket 接入
2. 支持其他服务通过 GRPC 调用
3. 多业务接入
4. 单用户多设备同时在线，接收消息
5. 支持单聊，群聊，朋友圈
6. 可扩展

#### 使用技术

数据库：Mysql+Redis
通讯：GRPC+Protobuf
语言：Golang

#### 目录结构

- api: 服务对外提供的接口
    - api_grpc: grpc接口
    - api_ws: websocket接口
- cmd: 服务启动入口
    - grpc: grpc服务启动
    - ws: websocket服务启动
- config: 服务配置
- internal: 每个服务私有代码
    - internal_grpc: grpc服务代码
    - internal_ws: websocket服务代码
- pkg: 服务共有代码
    - log: 日志代码
    - errs: 错误回调代码
    - util: 通用方法
    - proto: proto文件
    - pb: protp生成的go代码
	- db: 连接数据库代码
	- models: 结构体代码
- sql: 项目sql文件
- test: 测试脚本
- docs: 文档

#### 环境配置
创建文件：

conf.go
dev_conf.go
prod_conf.go
test_conf.go

在 conf.go 文件中编写配置如下：

``` Go
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
```

分别编写开发环境、生产环境、测试环境的配置，以下以开发环境为例：

``` Go
package config

func initDevConf() {
	DBConf = dbConf{
		MySQL:    "root:1111112@tcp(localhost:3306)/golang_im?charset=utf8&parseTime=true&loc=Local",
		RedisIP:  "127.0.0.1:6379",
		RedisPwd: "",
	}

	WSConf = wsConf{
		WSListenAddr: ":9091",
	}

	GRPCConf = grpcConf{
		GRPCListenAddr: ":9092",
	}
}
```