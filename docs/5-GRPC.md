## 编写GRPC服务器代码

### GRPC服务器

在目录 internal/internal_grpc 下新建 server.go 文件：

``` Go
package internal_grpc

import (
	"golang_im/api/api_grpc"
	"net"

	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
		return
	}

	s := grpc.NewServer() // 创建gRPC服务器
	var GrpcService api_grpc.GrpcService
	pb.RegisterGreeterServer(s, &GrpcService) // 在gRPC服务端注册服务
	reflection.Register(s)                    //在给定的gRPC服务器上注册服务器反射服务
	err = s.Serve(lis)
	if err != nil {
		log.Panic(err)
		return
	}

}
```

GrpcService 结构体中定义了 GRPC 接口，通过 pb.RegisterGreeterServer(s, &GrpcService) 注册了它。
在目录 api/api_grpc 下新建 grpc_service_logic.go 文件，编写一个可以向会话发送系统消息的接口：

``` Go
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
```

### 通过拦截器验证权限

我们没必要为每个接口添加验证权限的代码，只需要在 GRPC 拦截器中写一次就行了，拦截器将会在每个 GRPC 接口执行前执行。更改 server.go 文件：

``` Go
package internal_grpc

import (
	"context"
	"golang_im/api/api_grpc"
	"net"

	"golang_im/pkg/log"
	"golang_im/pkg/pb"
	"golang_im/pkg/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
		return
	}
	// 拦截器
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			log.Panic(err)
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor)) // 创建gRPC服务器
	var GrpcService api_grpc.GrpcService
	pb.RegisterGreeterServer(s, &GrpcService) // 在gRPC服务端注册服务
	reflection.Register(s)                    //在给定的gRPC服务器上注册服务器反射服务
	err = s.Serve(lis)
	if err != nil {
		log.Panic(err)
		return
	}

}

// auth 验证 Token
func auth(ctx context.Context) error {
	appId, userId, deviceId, err := util.GetCtxData(ctx)
	if err != nil {
		return err
	}
	token, err := util.GetCtxToken(ctx)
	if err != nil {
		return err
	}

	// 验证
	err = util.VerifyToken(appId, userId, deviceId, token)
	if err != nil {
		return err
	}

	return nil
}
```

## 编写GRPC客户端代码进行测试

### 编写GRPC客户端代码

在目录 test/grpc 下新建 main.go 文件：

``` Go
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
	return metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"app_id", "1",
		"user_id", "1",
		"device_id", "1",
		"token", `hSXW0YdriAXR2LRn6P6veGBVMU0TqtX8Fy/MtZsDzgYEAkxBBW9o6MZobuT1BOekNkUT6YiSRd3bgVleInf65uAfSqxVVLdHtIliu4fSm4x6qWQJ6GHYUCDezSuTe64ziBeL3bws1/N+WRYejdYa7rUdJ6c7Zp7/2FzcDSPaDts=`,
		"request_id", strconv.FormatInt(time.Now().UnixNano(), 10)))
}
```