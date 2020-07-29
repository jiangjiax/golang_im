package api_ws

import (
	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"google.golang.org/protobuf/proto"
)

type wsServiceChat struct{}

var WsServiceChat = new(wsServiceChat)

// 获取消息会话列表（及未读数，分页）
func (*wsServiceChat) GetConversationList(auth *pb.Auth, data []byte) {
	var in pb.ConversationReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.GetConversationList(&in, auth)
	if err != nil {
		conn.WriteMSG("getconversationlist", err, nil)
		return
	}
}

// 根据会话id获取消息
func (*wsServiceChat) Sync(auth *pb.Auth, data []byte) {
	var in pb.SyncReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.Sync(&in, auth)
	if err != nil {
		conn.WriteMSG("sync", err, nil)
		return
	}
}

// 发送消息
func (*wsServiceChat) SendMessage(auth *pb.Auth, data []byte) {
	var in pb.SendMessage
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.SendMessage(&in, auth)
	if err != nil {
		conn.WriteMSG("sendmessage", err, nil)
		return
	}
}

// 回执
func (*wsServiceChat) MessageRead(auth *pb.Auth, data []byte) {
	var in pb.MessageRead
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.MessageRead(&in, auth)
	if err != nil {
		conn.WriteMSG("messageread", err, nil)
		return
	}
}

// 发送好友/加群请求
func (*wsServiceChat) AddExamine(auth *pb.Auth, data []byte) {
	var in pb.AddExamine
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.AddExamine(&in, auth)
	if err != nil {
		conn.WriteMSG("addexamine", err, nil)
		return
	}
}

// 获取好友/加群请求
func (*wsServiceChat) GetExamine(auth *pb.Auth, data []byte) {
	var in pb.GetExamineReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.GetExamine(&in, auth)
	if err != nil {
		conn.WriteMSG("getexamine", err, nil)
		return
	}
}

// 处理好友/加群请求
func (*wsServiceChat) UpExamine(auth *pb.Auth, data []byte) {
	var in pb.UpExamineReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.UpExamine(&in, auth)
	if err != nil {
		conn.WriteMSG("upexamine", err, nil)
		return
	}
}

// 好友/加群请求未读数
func (*wsServiceChat) ExamineReadNum(auth *pb.Auth, data []byte) {
	conn, err := WsDaoChat.ExamineReadNum(auth)
	if err != nil {
		conn.WriteMSG("examinereadnum", err, nil)
		return
	}
}

// 聊天未读数
func (*wsServiceChat) ChatReadNum(auth *pb.Auth, data []byte) {
	conn, err := WsDaoChat.ChatReadNum(auth)
	if err != nil {
		conn.WriteMSG("chatreadnum", err, nil)
		return
	}
}

// 添加会话
func (*wsServiceChat) AddConversation(auth *pb.Auth, data []byte) {
	var in pb.AddConversationReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.AddConversationDao(&in, auth)
	if err != nil {
		conn.WriteMSG("addconversation", err, nil)
		return
	}
}

// 会话置顶
func (*wsServiceChat) UpConversation(auth *pb.Auth, data []byte) {
	var in pb.ConversationSettingReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.UpConversation(&in, auth)
	if err != nil {
		conn.WriteMSG("upconversation", err, nil)
		return
	}
}

// 会话免打扰
func (*wsServiceChat) DisturbConversation(auth *pb.Auth, data []byte) {
	var in pb.ConversationSettingReq
	err := proto.Unmarshal(data, &in)
	if err != nil {
		log.Warn(err)
		return
	}

	conn, err := WsDaoChat.DisturbConversation(&in, auth)
	if err != nil {
		conn.WriteMSG("disturbconversation", err, nil)
		return
	}
}
