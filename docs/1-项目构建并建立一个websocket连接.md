## 项目构建

### 项目需求

使用 Golang 完成一个即时通讯服务器，支持功能：
1. 支持前端 websocket 接入
2. 支持其他服务通过 GRPC 调用
3. 多业务接入
4. 单用户多设备同时在线，接收消息
5. 支持单聊，群聊，朋友圈
6. 可扩展

### 使用技术

数据库：Mysql+Redis
通讯：GRPC+Protobuf
语言：Golang

### 目录结构

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

### 环境配置
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

---

## 建立一个websocket连接

### websocket服务器

在目录 internal/ws 下创建文件 server.go 用于配置 websocket 的路由：

``` Go
package internal_ws

import (
	"fmt"
	"golang_im/pkg/log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	// CheckOrigin 检查源，如果 CheckOrigin 函数返回 false，则 Upgrade 方法会使 WebSocket 握手失败，HTTP 状态为403
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级协议
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("连接错误:", err)
		return
	}

	ctx := NewWSConn(Conn, 0, 0, 0)
	ctx.DoConn()
}

// websocket服务路由
func StartWSServer(address string) {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("websocket server start")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Panic(err)
		return
	}
}

```

客户端将通过路径 /ws 建立连接，接下来创建文件 ctx_conn.go 进行建立连接后的进一步处理：

``` Go
package internal_ws

import (
	"fmt"
	"io"
	"strings"
	"time"

	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// Controllers用于放置所有websocket接口的结构体
var Controllers = make(map[string]func(auth *pb.Auth, data []byte))

type WSConn struct {
	Conn     *websocket.Conn // websocket连接
	AppId    int64           // AppId
	DeviceId int64           // 设备id
	UserId   int64           // 用户id
}

func NewWSConn(conn *websocket.Conn, appId, userId, deviceId int64) *WSConn {
	return &WSConn{
		Conn:     conn,
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
	}
}

// DoConn 处理连接
func (c *WSConn) DoConn() {
	for {
		// 接收方法
		err := c.Conn.SetReadDeadline(time.Now().Add(time.Minute))
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		// 解析上行数据
		var input pb.Input
		err = proto.Unmarshal(data, &input)
		if err != nil {
			log.Error("解析上行数据失败:", err)
			c.Release()
		}
	}
}

// HandleReadErr 读取conn错误
func (c *WSConn) HandleReadErr(err error) {
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}
	c.Release()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		fmt.Println("客户端主动关闭连接或者异常程序退出")
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		fmt.Println("SetReadDeadline 之后，超时返回的错误")
		return
	}
}

// Release 释放连接
func (c *WSConn) Release() {
	// 关闭连接
	err := c.Conn.Close()
	if err != nil {
		log.Error("close err:", err)
	}
}
```

ctx_conn.go 就像一个中转站，每当用户发送请求时，都会通过接收方法 Conn.ReadMessage() 接收到。
注意，在接收方法上有一行代码 c.Conn.SetReadDeadline(time.Now().Add(time.Minute))，他设置了 time.Minute 也就是1分钟的时间，如果在这个通道上的用户超过1分钟没有发送数据，即关闭通道。这是为了验证该通道是否因意外而断开。因此，用户需要每隔一段时间（1分钟内）向通道发送一条消息，证明这个通道没有因意外而断开，这就是 心跳 机制。

最后，在 cmd/ws 下创建启动文件 main.go:

``` Go
package main

import (
	"golang_im/config"
	"golang_im/internal/internal_ws"
)

func main() {
	// 启动websocket服务器
	ws.StartWSServer(config.WSConf.WSListenAddr)
}
```

在目录下执行 go run main.go 即可启动服务。

### websocket客户端

在目录 test/ws 下创建 index.html 用于测试 websocket 服务器：

``` HTML
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>测试websocket连接</title>
</head>
<body>
    <button onclick="linkws()">连接websocket</button>
</body>
<script language="javascript" type="text/javascript">
    function linkws() {
        // 建立连接
        WebSockets = new WebSocket(`ws://127.0.0.1:9091/ws`);

        // 连接开启监听
        WebSockets.onopen = () => {
            console.log("open")
            WebSockets.send("111");
            WebSockets.onmessage = function (evt) {
                console.log("message")
            };
        },

        // 发生错误监听
        WebSockets.onerror = function(event) {
            console.log("Connected to WebSocket server error");
        },

        //连接关闭监听
        WebSockets.onclose = function(event) {
            console.log('WebSocket Connection Closed.');
        }
    }
</script>
</html>
```

[图片]
在浏览器的控制台中，可以看到连接 websocket 服务器成功了。