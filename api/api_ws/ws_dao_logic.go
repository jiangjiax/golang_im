package api_ws

import (
	"golang_im/internal/internal_ws"
	"golang_im/pkg/db"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/pb"

	"github.com/gogo/protobuf/proto"
)

type wsDaoLogic struct{}

var WsDaoLogic = new(wsDaoLogic)

// 获取用户信息
func (*wsDaoLogic) GetUserInfo(in *pb.GetUserInfoReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	var uid = in.UserId
	if uid == 0 {
		uid = auth.UserId
	}
	var GetUserInfo pb.GetUserInfoResp
	sql_str := `
	select user_id, nickname, sex, avatar_url, sign, account 
	from im_user 
	where app_id = ? and user_id = ?
	`
	err := db.DBCli.QueryRow(sql_str, auth.AppId, uid).
		Scan(&GetUserInfo.UserId, &GetUserInfo.Nickname, &GetUserInfo.Sex, &GetUserInfo.AvatarUrl, &GetUserInfo.Sign, &GetUserInfo.Account)
	if err != nil {
		return conn, err
	}

	msgBytes, err := proto.Marshal(&GetUserInfo)
	if err != nil {
		return conn, err
	}

	conn.WriteMSG("getuserinfo", nil, msgBytes)
	return conn, nil
}

// 心跳
func (*wsDaoLogic) Heartbeat(auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	conn.WriteMSG("heartbeat", nil, nil)
	return conn, nil
}
