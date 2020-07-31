package util

import (
	"golang_im/pkg/models"
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

func Group(groups []models.Group) ([]*pb.Group, error) {
	pbGroups := make([]*pb.Group, 0, len(groups))
	for i := range groups {
		pbGroups = append(pbGroups, &pb.Group{
			GroupId:      groups[i].GroupId,
			Name:         groups[i].Name,
			Introduction: groups[i].Introduction,
			UserMum:      groups[i].UserNum,
			Type:         groups[i].Type,
			Extra:        groups[i].Extra,
			Privacy:      groups[i].Privacy,
			Avatar:       groups[i].Avatar,
			CreateTime:   groups[i].CreateTime.Unix(),
			UpdateTime:   groups[i].UpdateTime.Unix(),
			UserType:     groups[i].UserType,
		})
	}
	return pbGroups, nil
}
