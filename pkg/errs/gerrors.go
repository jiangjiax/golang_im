package gerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthorized      = newError(1, "未登录")
	ErrUnDeviceid        = newError(2, "无设备")
	ErrGroupAlreadyExist = newError(3, "群组已经存在")
)

func newError(code int, message string) error {
	return status.New(codes.Code(code), message).Err()
}
