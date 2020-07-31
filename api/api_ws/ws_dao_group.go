package api_ws

import (
	"errors"
	"golang_im/internal/internal_ws"
	"golang_im/pkg/db"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/log"
	"golang_im/pkg/models"
	"golang_im/pkg/pb"
	"golang_im/pkg/util"
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

// 删除群组
func (*wsDaoGroup) Delete(appId int64, groupId int64) error {
	sql_str := `update im_group set status = 4, update_time = now() where app_id = ? and id = ?`
	_, err := db.DBCli.Exec(sql_str, appId, groupId)
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

// 更新群组
func (*wsDaoGroup) UpGroup(in *pb.CreateGroupReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	sets := ""
	if in.Group.Name != "" {
		sets += " name = '" + in.Group.Name + "', "
	}
	if in.Group.Introduction != "" {
		sets += " introduction = '" + in.Group.Introduction + "', "
	}
	if in.Group.Avatar != "" {
		sets += " avatar = '" + in.Group.Avatar + "', "
	}
	sql_str := `update im_group set ` + sets + ` where app_id = ? and id = ?`
	_, err := db.DBCli.Exec(sql_str, auth.AppId, in.Group.GroupId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	conn.WriteMSG("upgroup", nil, nil)
	return conn, nil
}

// 删除群组
func (*wsDaoGroup) DelGroup(in *pb.DeleteGroupReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	// 只能删除自己是群主的群
	istype, err := WsDaoChat.TypeByUser(auth, in.GroupId)
	if err != nil {
		log.Error(err)
		return conn, err
	}
	if istype != 3 {
		var Errlevel error = errors.New("权限不足")
		return conn, Errlevel
	}

	// 删除群
	err = WsDaoGroup.Delete(auth.AppId, in.GroupId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	// 删除会话
	err = WsDaoChat.DeleteConversationGroup(auth.AppId, in.GroupId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	conn.WriteMSG("delgroup", nil, nil)
	return conn, nil
}

// 删除群组
func (*wsDaoGroup) GroupByUser(in *pb.GetUserGroupsReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	intype := in.Type

	if intype == 1 {
		// 我管理的群
		sql_str := `
		select g.id, g.name, g.introduction, g.user_num, g.type, g.extra, g.create_time, g.update_time, g.privacy, g.avatar, u.type as usertype 
		from im_group_user as u 
		left join im_group as g on u.app_id = g.app_id and u.group_id = g.id and u.app_id = g.app_id 
		where u.app_id = ? and u.user_id = ? and g.status = 1 and (u.type = 2 or u.type = 3) and u.examine = 1 and u.status = 0 and g.status != 4 and g.status != 3
		`
		rows, err := db.DBCli.Query(sql_str, auth.AppId, auth.UserId)
		if err != nil {
			return conn, err
		}
		var groups []models.Group
		var group models.Group
		for rows.Next() {
			err := rows.Scan(&group.GroupId, &group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime, &group.Privacy, &group.Avatar, &group.UserType)
			if err != nil {
				return conn, err
			}
			groups = append(groups, group)
		}

		pbGroups, err := util.Group(groups)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		msgBytes, err := proto.Marshal(&pb.GetUserGroupsResp{Groups: pbGroups})
		if err != nil {
			log.Error(err)
			return conn, err
		}
		conn.WriteMSG("groupbyuser", nil, msgBytes)
		return conn, nil
	}
	if intype == 2 {
		sql_str := `
		select g.id, g.name, g.introduction, g.user_num, g.type, g.extra, g.create_time, g.update_time, g.privacy, g.avatar, u.type as usertype
		from im_group_user as u
		left join im_group as g on u.app_id = g.app_id and u.group_id = g.id and u.app_id = g.app_id
		where u.app_id = ? and u.user_id = ? and g.status = 1 and u.type = 1 and u.examine = 1 and u.status = 0 and g.status != 4 and g.status != 3,
		`
		// 我做为群成员的群
		rows, err := db.DBCli.Query(sql_str, auth.AppId, auth.UserId)
		if err != nil {
			return conn, err
		}
		var groups []models.Group
		var group models.Group
		for rows.Next() {
			err := rows.Scan(&group.GroupId, &group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime, &group.Privacy, &group.Avatar, &group.UserType)
			if err != nil {
				return conn, err
			}
			groups = append(groups, group)
		}

		pbGroups, err := util.Group(groups)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		msgBytes, err := proto.Marshal(&pb.GetUserGroupsResp{Groups: pbGroups})
		if err != nil {
			log.Error(err)
			return conn, err
		}
		conn.WriteMSG("groupbyuser", nil, msgBytes)
		return conn, nil
	}
	// 两者都获取
	sql_str := `
	select g.id, g.name, g.introduction, g.user_num, g.type, g.extra, g.create_time, g.update_time, g.privacy, g.avatar, u.type as usertype
	from im_group_user as u
	left join im_group as g on u.app_id = g.app_id and u.group_id = g.id and u.app_id = g.app_id
	where u.app_id = ? and u.user_id = ? and g.status = 1 and u.examine = 1 and u.status = 0 and g.status != 4 and g.status != 3,
	`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, auth.UserId)
	if err != nil {
		return conn, err
	}
	var groups []models.Group
	var group models.Group
	for rows.Next() {
		err := rows.Scan(&group.GroupId, &group.Name, &group.Introduction, &group.UserNum, &group.Type, &group.Extra, &group.CreateTime, &group.UpdateTime, &group.Privacy, &group.Avatar, &group.UserType)
		if err != nil {
			return conn, err
		}
		groups = append(groups, group)
	}

	pbGroups, err := util.Group(groups)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	msgBytes, err := proto.Marshal(&pb.GetUserGroupsResp{Groups: pbGroups})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("groupbyuser", nil, msgBytes)
	return conn, nil
}
