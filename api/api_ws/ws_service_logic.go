package api_ws

import (
	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"google.golang.org/protobuf/proto"
)

type wsServiceLogic struct{}

var WsServiceLogic = new(wsServiceLogic)

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

// 心跳
func (*wsServiceLogic) Heartbeat(auth *pb.Auth, data []byte) {
	WsDaoLogic.Heartbeat(auth)
}

// 发送消息
func (*wsServiceLogic) SendMessage(auth *pb.Auth, data []byte) {
	var in pb.GetUserInfoReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoLogic.GetUserInfo(&in, auth)
	if err != nil {
		conn.WriteMSG("sendmessage", err, nil)
		return
	}
}
