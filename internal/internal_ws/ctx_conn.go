package internal_ws

import (
	"fmt"
	"io"
	"strings"
	"time"

	"golang_im/internal/internal_ws/cache"
	"golang_im/pkg/db"
	"golang_im/pkg/log"
	"golang_im/pkg/pb"
	"golang_im/pkg/util"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/status"
)

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

		if c.AppId == 0 && c.UserId == 0 && c.DeviceId == 0 && input.Type == "sign" {
			// 登录
			c.Sign(input.Auth[0].AppId, input.Auth[0].UserId, input.Auth[0].DeviceId, input.Auth[0].Token)
		} else if input.Type != "sign" {
			fmt.Println("type:", input.Type)
			// 执行controller
			Controllers[input.Type](input.Auth[0], input.Data)
		}
	}
}

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

	if c.AppId != 0 && c.UserId != 0 && c.DeviceId != 0 {
		// 将设备设置为离线
		err := c.DeviceOnline(c.AppId, c.DeviceId, c.UserId, 0)
		if err != nil {
			log.Error(err)
		}
	}
}

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
