## 配置Protobuf

#### 编写proto文件

Protobuf可以将编写的proto文件转变为其他语言可以使用的协议，可以作为设计安全的跨语言RPC接口的基础工具。
proto文件中最基本的数据单元是message，是类似Go语言中结构体的存在。在message中可以嵌套message或其他的基础数据类型的成员。

在 pkg/proto 目录下创建 client_message.proto 文件：

``` protobuf
syntax = "proto3";
package pb;

// 上行数据
message Input {
    string type = 1; // 包的类型
    bytes data = 2; // 数据
    repeated Auth auth = 3;
}

// 下行数据
message Output {
    string type = 1; // 包的类型
    int32 code = 2; // 错误码
    string message = 3; // 错误信息
    bytes data = 4; // 数据
}

// 权限判断
message Auth {
    int64 app_id = 1; // app_id
    int64 device_id = 2; // 设备id
    int64 user_id = 3; // 用户id
    string token = 4; // 秘钥
}

// 获取用户信息
message GetUserInfoReq {
    int64 user_id = 1;
}
message GetUserInfoResp {
    int64 user_id = 1;
    string nickname = 2;
    int32 sex = 3; // 性别 0 未知 1 男 2 女
    string avatar_url = 4;
    string sign = 5;
    string account = 6;
}
```

client_message.proto 文件里定义了很多 message，进入项目查看全部的：https://github.com/jiangjiax/golang_im

#### 编译proto文件

安装编译工具 protoc，然后使用以下命令将 proto 文件编译为 Go 语言可以使用的协议：

``` shell
# 编译目录下所有protp文件
$ protoc --go_out=. *.proto

# 编译client_message.proto文件
$ protoc -I . --go_out=plugins=grpc:. client_message.proto
```

使用以下命令将 proto 文件编译为 Javascript 语言可以使用的协议：

``` shell
# 将proto编译为js文件
$ protoc --js_out=import_style=commonjs,binary:./ *.proto

# 将编译好的js文件转换为浏览器可以直接使用的js文件
$ browserify client_message_pb.js -o client_message.js
```
---

## 编写Websocket接口

#### 管理用户建立的连接

在 internal/internal_ws 目录下新建 manager.go 文件，使用数据结构 sync.Map 管理用户建立的连接，sync.Map 是一个有读写锁机制的 Map 结构：

``` Go
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

```

#### 心跳

在 internal/internal_ws/ctx_conn.go 文件的 DoConn() 函数中加入调用接口的代码：

``` Go
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

        // 执行controller调用接口
        Controllers[input.Type](input.Auth[0], input.Data)
	}
}
```

Controllers 是一个 map 结构，以字符串为健，以函数为值。input.Type 就相当于接口名，我们通过 Controllers[input.Type](input.Auth[0], input.Data) 调用对应接口的函数。

在 api/api_ws 目录下新建 ws_controller.go 文件：

``` Go
package api_ws

import (
	"fmt"
	"golang_im/internal/internal_ws"
)

func Controller_init() {
	// 基础
	logic_Controller()
	// 单聊、群聊
	chat_Controller()
	// 群组
	group_Controller()
	// 好友
	friend_Controller()
	// 朋友圈
	circle_Controller()

	fmt.Println("controller start")
}

// 基础
func logic_Controller() {
	// 心跳
	internal_ws.Controllers["heartbeat"] = WsServiceLogic.Heartbeat
}

// 单聊，群聊，好友/加群请求
func chat_Controller() {
}

// 群组
func group_Controller() {
}

// 好友
func friend_Controller() {
}

// 朋友圈
func circle_Controller() {
}
```

在同级目录下新建 ws_service_logic.go 文件：

``` Go
package api_ws

import (
	"golang_im/internal/internal_ws"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/pb"
)

type wsDaoLogic struct{}

var WsDaoLogic = new(wsDaoLogic)

// 心跳
func (*wsDaoLogic) Heartbeat(auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	conn.WriteMSG("heartbeat", nil, nil)
	return conn, nil
}
```

在 ws_service_logic.go 文件中，我们使用 internal_ws.Load(设备号) 获取对应设备的连接，然后我们调用 conn.WriteMSG() 方法往对应连接的通道中发送数据。

我们将在 internal/internal_ws/ctx_conn.go 文件中编写 conn.WriteMSG() 方法：

``` Go
// Output
func (c *WSConn) WriteMSG(pt string, err error, msgBytes []byte) {
	var output pb.Output

	output = pb.Output{
		Type: pt,
		Data: msgBytes,
	}

	if msgBytes != nil {
		output.Data = msgBytes
	}

	if err != nil {
		status, _ := status.FromError(err)
		output.Code = int32(status.Code())
		output.Message = status.Message()
	} else {
		output.Code = 200
		output.Message = "ok"
	}

	outputBytes, err := proto.Marshal(&output)
	if err != nil {
		return
	}
	err = c.Conn.WriteMessage(websocket.BinaryMessage, outputBytes)
	if err != nil {
		return
	}
}
```

#### 登录&设备在线/离线的设置

登录的接口我直接在 internal/internal_ws/ctx_conn.go 文件中写了，不放到 Controllers 中，先在 DoConn() 方法中判断：

``` Go
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

		if c.AppId == 0 && c.UserId == 0 && c.DeviceId == 0 && input.Type == "sign" {
			// 登录
			c.Sign(input.Auth[0].AppId, input.Auth[0].UserId, input.Auth[0].DeviceId, input.Auth[0].Token)
		} else if input.Type != "sign" {
			// 执行controller
			Controllers[input.Type](input.Auth[0], input.Data)
		}
	}
}
```

编写验证登录的方法：

``` Go
// 登录
func (c *WSConn) Sign(appid, userid, deviceid int64, token string) {
	ctx := NewWSConn(c.Conn, appid, userid, deviceid)

	// 验证
	err := util.VerifyToken(appid, userid, deviceid, token)
	if err != nil {
		c.Release()
	}

	// 将设备设置为在线
	err = c.DeviceOnline(appid, deviceid, userid, 1)
	if err != nil {
		log.Error(err)
	}

	// 断开这个设备之前的连接
	preCtx := Load(deviceid)
	if preCtx != nil {
		preCtx.DeviceId = -1
	}

	Store(deviceid, ctx)
}
```

编写将设备设置为在线或离线的方法，在释放连接时调用的 Release() 方法中也要判断：

``` Go
// 将设备设置为在线或离线
func (c *WSConn) DeviceOnline(appid, deviceid, user_id int64, status int) error {
	sql_str := `
		UPDATE im_device SET status = 1 
		WHERE app_id = ? and id = ? and user_id = ? and del = 1 
	`
	_, err := db.DBCli.Exec(sql_str, appid, deviceid, user_id)
	if err != nil {
		return err
	}

	// 清除该用户在线设备缓存
	cache.DeviceCache.Del(appid, user_id)

	return nil
}

// Release 释放连接
func (c *WSConn) Release() {
	// 关闭连接
	err := c.Conn.Close()
	if err != nil {
		log.Error("close err:", err)
	}

	if c.AppId != 0 && c.UserId != 0 && c.DeviceId != 0 {
		// 将设备设置为离线
		err := c.DeviceOnline(c.AppId, c.DeviceId, c.UserId, 0)
		if err != nil {
			log.Error(err)
		}
	}
}
```

#### 获取用户信息

更改 api/api_ws/ws_controller.go 文件，添加一个获取用户信息的接口：

``` Go
// 基础
func logic_Controller() {
	// 获取用户信息 [已测试]
	internal_ws.Controllers["getuserinfo"] = WsServiceLogic.GetUserInfo
	// 心跳 [已测试]
	internal_ws.Controllers["heartbeat"] = WsServiceLogic.Heartbeat
}
```

在 ws_service_logic.go 文件中加入对应的方法：

``` Go
// 获取用户信息
func (*wsServiceLogic) GetUserInfo(auth *pb.Auth, data []byte) {
	var in pb.GetUserInfoReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoLogic.GetUserInfo(&in, auth)
	if err != nil {
		conn.WriteMSG("getuserinfo", err, nil)
		return
	}
}
```

同样，在 ws_dao_logic.go 中加入对应的方法：

``` Go
// 获取用户信息
func (*wsDaoLogic) GetUserInfo(in *pb.GetUserInfoReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	var uid = in.UserId
	if uid == 0 {
		uid = auth.UserId
	}
	var GetUserInfo pb.GetUserInfoResp
	sql_str := `
	select user_id, nickname, sex, avatar_url, sign, account 
	from im_user 
	where app_id = ? and user_id = ?
	`
	err := db.DBCli.QueryRow(sql_str, auth.AppId, uid).
		Scan(&GetUserInfo.UserId, &GetUserInfo.Nickname, &GetUserInfo.Sex, &GetUserInfo.AvatarUrl, &GetUserInfo.Sign, &GetUserInfo.Account)
	if err != nil {
		return conn, err
	}

	msgBytes, err := proto.Marshal(&GetUserInfo)
	if err != nil {
		return conn, err
	}

	conn.WriteMSG("getuserinfo", nil, msgBytes)
	return conn, nil
}
```