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
