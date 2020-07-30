package api_ws

import (
	"golang_im/internal/internal_ws"
	"golang_im/pkg/db"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/log"
	"golang_im/pkg/pb"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
)

type wsDaoGroup struct{}

var WsDaoGroup = new(wsDaoGroup)

// 创建群组
func (*wsDaoGroup) CreatGroup(in *pb.CreateGroupReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	// 创建群组，获得群id
	sql_str := `insert ignore into im_group(app_id,name,introduction,type,extra,privacy,avatar,user_id,coordinatex,coordinatey,commandword) value(?,?,?,?,?,?,?,?,?,?,?)`
	result, err := db.DBCli.Exec(sql_str, auth.AppId, in.Group.Name, in.Group.Introduction, in.Group.Type, in.Group.Extra, in.Group.Privacy, in.Group.Avatar,
		auth.UserId, in.Group.Coordinatex, in.Group.Coordinatey, in.Group.Commandword)
	if err != nil {
		log.Error(err)
		return conn, err
	}
	num, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return conn, err
	}
	if num == 0 {
		return conn, gerrors.ErrGroupAlreadyExist
	}
	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		log.Error(err)
		return conn, err
	}

	// 建个会话
	selfConversation := pb.ConversationItem{
		ReceiverType: int32(2),
		ReceiverId:   lastInsertID,
	}
	ConversationId, err := WsDaoChat.AddConversation(&selfConversation, auth, auth.UserId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	// 添加群主
	AddUserbyGroup := pb.AddUserbyGroup{
		GroupId: lastInsertID,
		UserId:  auth.UserId,
		Istype:  3,
		Examine: 1,
		Way:     6,
	}
	err = WsDaoGroup.AddUser(&AddUserbyGroup, auth)
	if err != nil {
		return conn, err
	}

	// 批量添加群员
	userIds := strings.Split(in.Group.UserIds, `,`)
	for _, str := range userIds {
		uid, err := strconv.ParseInt(str, 10, 64)
		// 看看有这个人吗
		errp := WsDaoGroup.People(auth.AppId, uid)
		if errp == nil {
			// 添加群员
			AddUserbyGroup := pb.AddUserbyGroup{
				GroupId: lastInsertID,
				UserId:  uid,
				Istype:  1,
				Examine: 1,
				Way:     7,
			}
			err = WsDaoGroup.AddUser(&AddUserbyGroup, auth)
			if err != nil {
				return conn, err
			}
		}
	}

	msgBytes, err := proto.Marshal(&pb.CreateGroupResp{GroupLastInsertID: lastInsertID, ConversationLastInsertID: ConversationId})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("creatgroup", nil, msgBytes)
	return conn, nil
}

// 将用户添加到群组
func (*wsDaoGroup) AddUser(in *pb.AddUserbyGroup, auth *pb.Auth) error {
	is_read := 0
	if in.Examine == 1 {
		// 未读数判断
		sql_str := `SELECT id from im_message where app_id = ? and receiver_type = 2 and receiver_id = ?`
		row := db.DBCli.QueryRow(sql_str, auth.AppId, in.GroupId)
		err := row.Scan(&is_read)
		if err != nil {
			is_read = 0
		}
	}
	sql_str := `insert into im_group_user(app_id, group_id, user_id, label, type, way, examine, examinetext, is_read) values(?,?,?,?,?,?,?,?,?) 
	ON DUPLICATE KEY UPDATE is_read=?, examine=?, update_time=now(), status = 0, create_time=now()`
	_, err := db.DBCli.Exec(sql_str, auth.AppId, in.GroupId, in.UserId, in.Label, in.Istype, in.Way, in.Examine, in.Examinetext, is_read, is_read, in.Examine)
	if err != nil {
		return err
	}

	return nil
}

// 看看有这个人吗
func (*wsDaoGroup) People(appId, userId int64) error {
	sql_str := `select id from im_user where app_id = ? and user_id = ? and status = 1`
	row := db.DBCli.QueryRow(sql_str, appId, userId)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
