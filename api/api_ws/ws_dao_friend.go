package api_ws

import (
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
