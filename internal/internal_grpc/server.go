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
