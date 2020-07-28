package main

import (
	"context"
	"fmt"
	"golang_im/config"
	"golang_im/pkg/pb"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// 连接服务器
	conn, err := grpc.Dial("127.0.0.1"+config.GRPCConf.GRPCListenAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	// 调用服务端的SayHello
	SendMessage := &pb.SendMessage{
		SenderType:   1,
		ReceiverType: 1,
		ReceiverId:   1,
		MessageBody: &pb.MessageBody{
			MessageType: pb.MessageType_MT_TEXT,
			MessageContent: &pb.MessageContent{
				Content: &pb.MessageContent_Text{
					Text: &pb.Text{
						Text: "这是一条系统发送的消息",
					},
				},
			},
		},
	}
	r, err := c.SystemMsgByConversation(getCtx(), SendMessage)
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	fmt.Println(r)
}

func getCtx() context.Context {
	// token, _ := util.GetToken(1, 2, 3, time.Now().Add(1*time.Hour).Unix(), util.PublicKey)
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"app_id", "1",
		"user_id", "1",
		"device_id", "1",
		"token", `hSXW0YdriAXR2LRn6P6veGBVMU0TqtX8Fy/MtZsDzgYEAkxBBW9o6MZobuT1BOekNkUT6YiSRd3bgVleInf65uAfSqxVVLdHtIliu4fSm4x6qWQJ6GHYUCDezSuTe64ziBeL3bws1/N+WRYejdYa7rUdJ6c7Zp7/2FzcDSPaDts=`,
		"request_id", strconv.FormatInt(time.Now().UnixNano(), 10)))
}
