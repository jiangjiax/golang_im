package api_ws

import (
	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"google.golang.org/protobuf/proto"
)

type wsServiceFriend struct{}

var WsServiceFriend = new(wsServiceFriend)

// 获取好友列表
func (*wsServiceFriend) GetFriendList(auth *pb.Auth, data []byte) {
	var in pb.GetFriendListReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoFriend.GetFriendList(&in, auth)
	if err != nil {
		conn.WriteMSG("getfriends", err, nil)
		return
	}
}

// 搜索好友
func (*wsServiceFriend) SearchFriendList(auth *pb.Auth, data []byte) {
	var in pb.GetFriendListReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoFriend.SearchFriendList(&in, auth)
	if err != nil {
		conn.WriteMSG("searchfriends", err, nil)
		return
	}
}

// 更新好友备注
func (*wsServiceFriend) UpRemarkFriend(auth *pb.Auth, data []byte) {
	var in pb.AddFriendReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoFriend.UpRemarkFriend(&in, auth)
	if err != nil {
		conn.WriteMSG("upremarkfriend", err, nil)
		return
	}
}

// 删除好友
func (*wsServiceFriend) DeleteFriend(auth *pb.Auth, data []byte) {
	var in pb.DeleteFriendReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoFriend.DeleteFriend(&in, auth)
	if err != nil {
		conn.WriteMSG("deletefriend", err, nil)
		return
	}
}
