package api_ws

import (
	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"google.golang.org/protobuf/proto"
)

type wsServiceCircle struct{}

var WsServiceCircle = new(wsServiceCircle)

// 发朋友圈
func (*wsServiceChat) AddTrend(auth *pb.Auth, data []byte) {
	var in pb.AddTrend
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoCircle.AddTrend(&in, auth)
	if err != nil {
		conn.WriteMSG("addtrend", err, nil)
		return
	}
}

// 获取朋友圈列表（包括内容、点赞人列表、评论列表，下拉加载）
func (*wsServiceChat) GetTrends(auth *pb.Auth, data []byte) {
	var in pb.GetTrendsReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoCircle.GetTrends(&in, auth)
	if err != nil {
		conn.WriteMSG("gettrends", err, nil)
		return
	}
}

// 点赞与取消点赞
func (*wsServiceChat) Thumb(auth *pb.Auth, data []byte) {
	var in pb.Thumb
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoCircle.Thumb(&in, auth)
	if err != nil {
		conn.WriteMSG("thumb", err, nil)
		return
	}
}

// 评论与回复
func (*wsServiceChat) AddTrendsComment(auth *pb.Auth, data []byte) {
	var in pb.AddTrendsComment
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoCircle.AddTrendsComment(&in, auth)
	if err != nil {
		conn.WriteMSG("addtrendscomment", err, nil)
		return
	}
}
