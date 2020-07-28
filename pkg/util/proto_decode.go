package util

import (
	"golang_im/pkg/pb"

	"github.com/golang/protobuf/proto"
)

func SendMessageItem(selfMessage *pb.MessageItem) ([]byte, error) {
	msgBytes, err := proto.Marshal(selfMessage)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}

func SyncResp(SyncResp *pb.SyncResp) ([]byte, error) {
	msgBytes, err := proto.Marshal(SyncResp)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}
