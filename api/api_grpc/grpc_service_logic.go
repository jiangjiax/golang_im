package api_grpc

import (
	"context"
	"golang_im/api/api_ws"
	"golang_im/pkg/log"
	"golang_im/pkg/pb"
	"golang_im/pkg/util"
	"sync"
)

type GrpcService struct{}

var wg sync.WaitGroup

// 向会话发送系统消息
func (*GrpcService) SystemMsgByConversation(ctx context.Context, in *pb.SendMessage) (*pb.SystemMessageResp, error) {
	SystemMessageResp := &pb.SystemMessageResp{
		ErrCode: 500,
		ErrMsg:  "解析错误",
	}
	appId, userId, deviceId, err := util.GetCtxData(ctx)
	if err != nil {
		return SystemMessageResp, err
	}
	token, err := util.GetCtxToken(ctx)
	if err != nil {
		return SystemMessageResp, err
	}
	auth := &pb.Auth{
		AppId:    appId,
		DeviceId: deviceId,
		UserId:   userId,
		Token:    token,
	}

	// 发送消息
	if in.ReceiverType == 1 {
		// 单聊消息发送
		SystemMessageResp = &pb.SystemMessageResp{
			ErrCode: 500,
			ErrMsg:  "单聊消息发送错误",
		}
		if in.SenderType == 2 {
			// 发送者
			err := api_ws.WsDaoChat.Send(in, auth, in.SenderId)
			if err != nil {
				log.Error(err)
				return SystemMessageResp, err
			}
			// 接收者
			err = api_ws.WsDaoChat.Send(in, auth, in.ReceiverId)
			if err != nil {
				log.Error(err)
				return SystemMessageResp, err
			}
		} else if in.SenderType == 1 {
			// 接收者
			err := api_ws.WsDaoChat.Send(in, auth, in.ReceiverId)
			if err != nil {
				log.Error(err)
				return SystemMessageResp, err
			}
		}
	} else if in.ReceiverType == 2 {
		// 群聊消息发送
		SystemMessageResp = &pb.SystemMessageResp{
			ErrCode: 500,
			ErrMsg:  "群聊消息发送错误",
		}
		// 查找群组内的所有用户
		GroupUserReq := pb.GroupUserReq{
			GroupId: in.ReceiverId,
		}
		GroupUserResp, err := api_ws.WsDaoChat.GetGroupUsers(&GroupUserReq, auth)
		if err != nil {
			log.Error(err)
			return SystemMessageResp, err
		}
		for _, GroupUser := range GroupUserResp.GroupUser {
			go func(in *pb.SendMessage, auth *pb.Auth, toUserId int64) {
				wg.Add(1)
				err = api_ws.WsDaoChat.Send(in, auth, in.ReceiverId)
				if err != nil {
					log.Error(err)
				}
				wg.Done()
			}(in, auth, GroupUser.UserId)
		}
	}
	SystemMessageResp = &pb.SystemMessageResp{
		ErrCode: 200,
		ErrMsg:  "ok",
	}
	return SystemMessageResp, nil
}
