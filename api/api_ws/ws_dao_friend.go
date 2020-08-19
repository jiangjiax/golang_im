package api_ws

import (
	"errors"
	"golang_im/internal/internal_ws"
	"golang_im/pkg/db"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/log"
	"golang_im/pkg/pb"

	"github.com/gogo/protobuf/proto"
)

type wsDaoFriend struct{}

var WsDaoFriend = new(wsDaoFriend)

// 获取好友列表
func (*wsDaoFriend) GetFriendList(in *pb.GetFriendListReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	sql_str := `
	select f.fid, f.remark, u.avatar_url, u.nickname 
	from im_friend as f 
	left join im_user as u on u.user_id = f.fid and u.app_id = f.app_id 
	where f.app_id = ? and f.user_id = ? and f.status = 1 and f.examine = 1 and u.user_id = f.fid
	`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, auth.UserId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var FriendItems []*pb.FriendItem
	for rows.Next() {
		var fi pb.FriendItem
		err := rows.Scan(&fi.FId, &fi.Remark, &fi.AvatarUrl, &fi.Nickname)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		FriendItems = append(FriendItems, &fi)
	}

	msgBytes, err := proto.Marshal(&pb.GetFriendListResp{Friends: FriendItems})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("getfriends", nil, msgBytes)
	return conn, nil
}

// 搜索好友
func (*wsDaoFriend) SearchFriendList(in *pb.GetFriendListReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	sql_str := `
	select f.fid, f.remark, u.avatar_url, u.nickname
	from im_friend as f 
	left join im_user as u on u.user_id = f.fid and u.app_id = f.app_id 
	where f.app_id = ? and f.userid = ? and f.status = 1 and f.examine = 1 and f.remark = ? and u.nickname = ? and u.user_id = f.fid
	`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, auth.UserId, in.Keyword, in.Keyword)
	if err != nil {
		log.Error(err)
		return conn, err
	}
	var FriendItems []*pb.FriendItem
	for rows.Next() {
		var fi pb.FriendItem
		err := rows.Scan(&fi.FId, &fi.Remark, &fi.AvatarUrl, &fi.Nickname)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		FriendItems = append(FriendItems, &fi)
	}

	msgBytes, err := proto.Marshal(&pb.GetFriendListResp{Friends: FriendItems})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("searchfriends", nil, msgBytes)
	return conn, nil
}

// 更新好友备注
func (*wsDaoFriend) UpRemarkFriend(in *pb.AddFriendReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	if in.Friends.FId == auth.UserId {
		var Errlevel error = errors.New("不能备注自己")
		return conn, Errlevel
	}

	sql_str := `
	INSERT INTO friend (userid, fid, remark, way, app_id, status) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE remark=?, update_time=now()`
	_, err := db.DBCli.Exec(sql_str, auth.UserId, in.Friends.FId, in.Friends.Remark, in.Friends.Way, auth.AppId, 0, in.Friends.Remark)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	conn.WriteMSG("upremarkfriend", nil, nil)
	return conn, nil
}

// 删除好友
func (*wsDaoFriend) DeleteFriend(in *pb.DeleteFriendReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	sql_str := `
	update im_friend set status = 0 where app_id = ? and fid = ? and userid = ?`
	_, err := db.DBCli.Exec(sql_str, auth.AppId, in.FId, in.UserId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	// 删除好友会话
	err = WsDaoChat.DeleteConversationFriend(auth.AppId, in.FId, in.UserId)
	if err != nil {
		log.Error(err)
		return conn, err
	}
	err = WsDaoChat.DeleteConversationFriend(auth.AppId, in.UserId, in.FId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	conn.WriteMSG("deletefriend", nil, nil)
	return conn, nil
}
