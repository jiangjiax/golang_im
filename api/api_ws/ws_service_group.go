package api_ws

import (
	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"google.golang.org/protobuf/proto"
)

type wsServiceGroup struct{}

var WsServiceGroup = new(wsServiceGroup)

// 获取用户信息
func (*wsServiceGroup) CreatGroup(auth *pb.Auth, data []byte) {
	var in pb.CreateGroupReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoGroup.CreatGroup(&in, auth)
	if err != nil {
		conn.WriteMSG("creatgroup", err, nil)
		return
	}
}
