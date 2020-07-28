package internal_ws

import (
	"golang_im/pkg/log"
	"golang_im/pkg/pb"
	"golang_im/pkg/util"

	jsoniter "github.com/json-iterator/go"
)

// 创建一个消息体类型
func NewMessageBody(msgType int, msgContent string) *pb.MessageBody {
	content := pb.MessageContent{}
	switch pb.MessageType(msgType) {
	case pb.MessageType_MT_TEXT:
		var text pb.Text
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &text)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Text{Text: &text}
	case pb.MessageType_MT_FACE:
		var face pb.Face
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &face)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Face{Face: &face}
	case pb.MessageType_MT_VOICE:
		var voice pb.Voice
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &voice)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Voice{Voice: &voice}
	case pb.MessageType_MT_IMAGE:
		var image pb.Image
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &image)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Image{Image: &image}
	case pb.MessageType_MT_FILE:
		var file pb.File
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &file)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_File{File: &file}
	case pb.MessageType_MT_LOCATION:
		var location pb.Location
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &location)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Location{Location: &location}
	case pb.MessageType_MT_COMMAND:
		var command pb.Command
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &command)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Command{Command: &command}
	case pb.MessageType_MT_CUSTOM:
		var custom pb.Custom
		err := jsoniter.Unmarshal(util.Str2bytes(msgContent), &custom)
		if err != nil {
			log.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Custom{Custom: &custom}
	}

	return &pb.MessageBody{
		MessageType:    pb.MessageType(msgType),
		MessageContent: &content,
	}
}
