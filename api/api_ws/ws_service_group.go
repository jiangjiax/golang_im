package api_ws

import (
	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"google.golang.org/protobuf/proto"
)

type wsServiceGroup struct{}

var WsServiceGroup = new(wsServiceGroup)

// 创建群组
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

// 更新群组
func (*wsServiceGroup) UpGroup(auth *pb.Auth, data []byte) {
	var in pb.CreateGroupReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoGroup.UpGroup(&in, auth)
	if err != nil {
		conn.WriteMSG("upgroup", err, nil)
		return
	}
}

// 删除群组
func (*wsServiceGroup) DelGroup(auth *pb.Auth, data []byte) {
	var in pb.DeleteGroupReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoGroup.DelGroup(&in, auth)
	if err != nil {
		conn.WriteMSG("delgroup", err, nil)
		return
	}
}

// 获取用户加入的所有群组
func (*wsServiceGroup) GroupByUser(auth *pb.Auth, data []byte) {
	var in pb.GetUserGroupsReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoGroup.GroupByUser(&in, auth)
	if err != nil {
		conn.WriteMSG("groupbyuser", err, nil)
		return
	}
}
