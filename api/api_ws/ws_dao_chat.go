package api_ws

import (
	"database/sql"
	"errors"
	"golang_im/internal/internal_ws"
	"golang_im/internal/internal_ws/cache"
	"golang_im/pkg/db"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/log"
	"golang_im/pkg/models"
	"golang_im/pkg/pb"
	"golang_im/pkg/util"
	"strconv"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
)

type wsDaoChat struct{}

var WsDaoChat = new(wsDaoChat)

var wg sync.WaitGroup

// 获取会话列表
func (*wsDaoChat) GetConversationList(in *pb.ConversationReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	limit := in.Limit
	if in.Limit == 0 {
		limit = 15
	}
	sql_str := `select id, receiver_id, receiver_type, disturb, top, message, sender_id, UNIX_TIMESTAMP(update_time), sender_name, help
	from im_conversation
	where app_id = ? and user_id = ? limit ?,?`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, auth.UserId, in.Offset, limit)
	if err != nil {
		log.Error(err)
		return conn, err
	}
	var Conversations []*pb.ConversationItem
	for rows.Next() {
		var cs pb.ConversationItem
		err := rows.Scan(&cs.Id, &cs.ReceiverId, &cs.ReceiverType, &cs.Disturb, &cs.Top, &cs.Messagenewcontent, &cs.SenderId,
			&cs.UpdateTime, &cs.SenderName, &cs.Help)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		if cs.ReceiverType == 1 {
			// 单聊会话接收人昵称
			sql_str := `
			select nickname, avatar_url
			from im_user 
			where app_id = ? and user_id = ?
			`
			err := db.DBCli.QueryRow(sql_str, auth.AppId, cs.ReceiverId).
				Scan(&cs.Username, &cs.AvatarUrl)
			if err != nil {
				log.Error(err)
				return conn, err
			}

			// 未读数量
			sql_str = `
			select count(*)  
			from im_message 
			where app_id = ? and conversation_id = ? and is_read = 0 and status != 2 and status != 1 
			`
			err = db.DBCli.QueryRow(sql_str, auth.AppId, cs.Id).
				Scan(&cs.ReceiveNum)
			if err != nil {
				log.Error(err)
				return conn, err
			}
		} else if cs.ReceiverType == 2 {
			// 群聊会话接收群的群名
			sql_str := `
			select name, avatar
			from im_group 
			where app_id = ? and id = ?
			`
			err := db.DBCli.QueryRow(sql_str, auth.AppId, cs.ReceiverId).
				Scan(&cs.Username, &cs.AvatarUrl)
			if err != nil {
				log.Error(err)
				return conn, err
			}

			// 未读数量
			sql_str = `
			select count(*)
			from message as m
			where m.app_id = ? and m.receiver_type = 2 and m.receiver_id = ? and m.status != 2 and m.status != 1 
			and m.id > (select is_read from group_user as g where g.user_id = ? and g.app_id = m.app_id and g.group_id = m.receiver_id)
			`
			err = db.DBCli.QueryRow(sql_str, auth.AppId, cs.ReceiverId, auth.UserId).
				Scan(&cs.ReceiveNum)
			if err != nil {
				log.Error(err)
				return conn, err
			}

			// 群聊会话获取最新记录
			sql_str = `select conversation_message from im_message where app_id = ? and receiver_type = 2 and receiver_id = ? 
			order by id desc limit 1`
			err = db.DBCli.QueryRow(sql_str, auth.AppId, cs.ReceiverId).
				Scan(&cs.Messagenewcontent)
			if err != nil {
				log.Error(err)
				return conn, err
			}
		}

		Conversations = append(Conversations, &cs)
	}

	msgBytes, err := proto.Marshal(&pb.ConversationResp{Conversation: Conversations})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("getconversationlist", nil, msgBytes)
	return conn, nil
}

// 根据会话id获取消息
func (*wsDaoChat) Sync(in *pb.SyncReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	limit := in.Limit
	if in.Limit == 0 {
		limit = 15
	}
	var istype int
	var content string
	var SyncResp []*pb.MessageItem
	if in.ReceiverType == 1 {
		// 单聊
		sql_str := `
		select m.id, m.sender_type, m.sender_id, m.receiver_type, m.receiver_id, m.type, m.content, 
		UNIX_TIMESTAMP(m.send_time), m.status, m.conversation_id, m.help, u.nickname, u.avatar_url 
		from im_message as m 
		left join im_user as u on u.user_id = m.sender_id
		where m.app_id = ? and m.conversation_id = ? and m.status != 2 and u.app_id = ? and u.status = 1 order by id DESC limit ?,?
		`
		rows, err := db.DBCli.Query(sql_str, auth.AppId, in.ConversationId, auth.AppId, in.Offset, limit)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		for rows.Next() {
			var mi pb.MessageItem
			err := rows.Scan(&mi.Id, &mi.SenderType, &mi.SenderId, &mi.ReceiverType, &mi.ReceiverId, &istype, &content,
				&mi.SendTime, &mi.Status, &mi.ConversationId, &mi.Help, &mi.SenderName, &mi.SenderAvatar)
			if err != nil {
				log.Error(err)
				return conn, err
			}

			mi.MessageBody = internal_ws.NewMessageBody(istype, content)
			SyncResp = append(SyncResp, &mi)
		}

		// 已读
		sql_str = `update im_message set is_read = 1 where conversation_id = ?`
		_, err = db.DBCli.Exec(sql_str, in.ConversationId)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	} else if in.ReceiverType == 2 {
		// 群聊
		sql_str := `
		select m.id, m.sender_type, m.sender_id, m.receiver_type, m.receiver_id, m.to_user_ids, m.type, m.content, 
		UNIX_TIMESTAMP(m.send_time), m.status, m.conversation_id, m.help, u.nickname, u.avatar_url 
		from im_message as m 
		left join im_user as u on u.user_id = m.sender_id
		where receiver_id in (select group_id as g from im_group_user where g.app_id = m.app_id and g.user_id = ? and g.status = 0 and g.examine = 1) 
		and receiver_type = 2 and m.app_id = ? and m.status != 2 and u.app_id = ? and u.status = 1 order by id DESC limit ?,?
		`
		rows, err := db.DBCli.Query(sql_str, auth.UserId, auth.AppId, auth.AppId, in.Offset, limit)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		var mi pb.MessageItem
		for rows.Next() {
			err := rows.Scan(&mi.Id, &mi.SenderType, &mi.SenderId, &mi.ReceiverType, &mi.ReceiverId, &mi.ToUserIds,
				&istype, &content, &mi.SendTime, &mi.Status, &mi.ConversationId, &mi.Help, &mi.SenderName, &mi.SenderAvatar)
			if err != nil {
				log.Error(err)
				return conn, err
			}

			mi.MessageBody = internal_ws.NewMessageBody(istype, content)
			SyncResp = append(SyncResp, &mi)
		}

		// 已读
		sql_str = `update im_group_user set is_read = ? where user_id = ? and group_id = ?`
		_, err = db.DBCli.Exec(sql_str, SyncResp[len(SyncResp)-1].Id, auth.UserId, mi.ReceiverId)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	}

	msgBytes, err := util.SyncResp(&pb.SyncResp{Messages: SyncResp})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("sync", nil, msgBytes)
	return conn, nil
}

// SendMessage 发送消息
func (*wsDaoChat) SendMessage(in *pb.SendMessage, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	// 获取发送者昵称和发送者头像
	in.SenderId = auth.UserId
	if in.SenderType == 2 {
		sql_str := `
		select g.label, u.avatar_url from im_user as u left join group_user as g on g.user_id = u.user_id 
		where u.user_id = ? and g.group_id = ` + strconv.FormatInt(in.ReceiverId, 10)
		if in.ReceiverType == 1 {
			sql_str = `
			select nickname, avatar_url from im_user where user_id = ?
			`
		}
		err := db.DBCli.QueryRow(sql_str, in.SenderId).
			Scan(&in.SenderName, &in.SenderAvatar)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	}

	// 发送消息
	if in.ReceiverType == 1 {
		// 单聊消息发送
		// 发送者
		err := WsDaoChat.Send(in, auth, in.SenderId)
		if err != nil {
			log.Error(err)
			return conn, err
		}
		// 接收者
		err = WsDaoChat.Send(in, auth, in.ReceiverId)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	} else if in.ReceiverType == 2 {
		// 群聊消息发送
		// 查找群组内的所有用户
		GroupUserReq := pb.GroupUserReq{
			GroupId: in.ReceiverId,
		}
		GroupUserResp, err := WsDaoChat.GetGroupUsers(&GroupUserReq, auth)
		if err != nil {
			log.Error(err)
			return conn, err
		}
		for _, GroupUser := range GroupUserResp.GroupUser {
			go func(in *pb.SendMessage, auth *pb.Auth, toUserId int64) {
				wg.Add(1)
				err = WsDaoChat.Send(in, auth, in.ReceiverId)
				if err != nil {
					log.Error(err)
				}
				wg.Done()
			}(in, auth, GroupUser.UserId)
		}
	}

	conn.WriteMSG("sendmessage", nil, nil)
	return conn, nil
}

// Send 发送消息
func (*wsDaoChat) Send(in *pb.SendMessage, auth *pb.Auth, toUserId int64) error {
	// 消息体解析
	var content interface{}
	switch in.MessageBody.MessageType {
	case pb.MessageType_MT_TEXT:
		content = in.MessageBody.MessageContent.GetText()
	case pb.MessageType_MT_FACE:
		content = in.MessageBody.MessageContent.GetFace()
	case pb.MessageType_MT_VOICE:
		content = in.MessageBody.MessageContent.GetVoice()
	case pb.MessageType_MT_IMAGE:
		content = in.MessageBody.MessageContent.GetImage()
	case pb.MessageType_MT_FILE:
		content = in.MessageBody.MessageContent.GetFile()
	case pb.MessageType_MT_LOCATION:
		content = in.MessageBody.MessageContent.GetLocation()
	case pb.MessageType_MT_COMMAND:
		content = in.MessageBody.MessageContent.GetCommand()
	case pb.MessageType_MT_CUSTOM:
		content = in.MessageBody.MessageContent.GetCustom()
	case pb.MessageType_MT_VIDEO:
		content = in.MessageBody.MessageContent.GetVideo()
	}
	bytes, err := jsoniter.Marshal(content)
	if err != nil {
		log.Error(err)
		return err
	}

	MessageType := int32(in.MessageBody.MessageType)
	MessageContent := util.Bytes2str(bytes)
	// 最新记录
	ConversationMessage := ""
	switch in.MessageBody.MessageType {
	case pb.MessageType_MT_TEXT:
		var text pb.Text
		err = jsoniter.Unmarshal(util.Str2bytes(MessageContent), &text)
		if err != nil {
			log.Error(err)
			return err
		}
		decodetext := util.UnicodeEmojiCode(text.Text)
		ConversationMessage = decodetext
	case pb.MessageType_MT_FACE:
		ConversationMessage = "[表情]"
	case pb.MessageType_MT_VOICE:
		ConversationMessage = "[录音]"
	case pb.MessageType_MT_IMAGE:
		ConversationMessage = "[图像]"
	case pb.MessageType_MT_FILE:
		ConversationMessage = "[文件]"
	case pb.MessageType_MT_LOCATION:
		ConversationMessage = "[地址]"
	case pb.MessageType_MT_COMMAND:
		ConversationMessage = "[指令推送]"
	case pb.MessageType_MT_CUSTOM:
		ConversationMessage = "[自定义]"
	case pb.MessageType_MT_VIDEO:
		ConversationMessage = "[视频]"
	}
	// 持久化消息
	var mid int64
	selfMessage := pb.MessageItem{
		SenderType:          in.SenderType,
		SenderId:            in.SenderId,
		ReceiverType:        in.ReceiverType,
		ReceiverId:          in.ReceiverId,
		ToUserIds:           in.ToUserIds,
		SenderName:          in.SenderName,
		SenderAvatar:        in.SenderAvatar,
		Help:                in.Help,
		ConversationMessage: ConversationMessage,
		MessageBody: &pb.MessageBody{
			MessageType:    in.MessageBody.MessageType,
			MessageContent: in.MessageBody.MessageContent,
		},
	}
	selfConversation := pb.ConversationItem{
		SenderId:          in.SenderId,
		ReceiverType:      MessageType,
		ReceiverId:        in.ReceiverId,
		Messagenewcontent: ConversationMessage,
		Help:              in.Help,
		SenderName:        in.SenderName,
	}
	if in.ReceiverType == 1 {
		// 单聊持久化发送者和接收者的消息
		if selfConversation.ReceiverId == toUserId {
			selfConversation.ReceiverId = selfConversation.SenderId
		}
		ConversationId, err := WsDaoChat.AddConversation(&selfConversation, auth, toUserId)
		if err != nil {
			log.Error(err)
			return err
		}
		selfMessage.ConversationId = ConversationId
		// 将发送者设置为已读
		var read int32
		if toUserId == in.SenderId {
			// 已读
			read = 1
		} else {
			// 未读
			read = 0
		}
		mid, err = WsDaoChat.Add(&selfMessage, auth, toUserId, MessageType, read, MessageContent)
		if err != nil {
			return err
		}
	} else if in.ReceiverType == 2 && toUserId == in.SenderId {
		// 群聊只持久化发送者的消息
		ConversationId, err := WsDaoChat.AddConversation(&selfConversation, auth, toUserId)
		if err != nil {
			return err
		}
		selfMessage.ConversationId = ConversationId
		mid, err = WsDaoChat.Add(&selfMessage, auth, toUserId, MessageType, 1, MessageContent)
		if err != nil {
			return err
		}
		// 更新发送者用户消息索引
		err = WsDaoChat.Upindex(auth.AppId, mid, in.ReceiverId, toUserId)
		if err != nil {
			return err
		}
	}
	selfMessage.Id = mid

	// 将消息发送给需要接收的用户
	// 查询用户在线设备
	devices, err := WsDaoChat.ListOnlineByUserId(auth.AppId, toUserId)
	if err != nil {
		return err
	}
	for i := range devices {
		// 消息不需要投递给发送消息的设备
		if auth.DeviceId == devices[i].Id {
			continue
		}

		conn := internal_ws.Load(devices[i].Id)
		if conn == nil {
			continue
		}

		msgBytes, err := util.SendMessageItem(&selfMessage)
		if err != nil {
			log.Error(err)
			return err
		}
		conn.WriteMSG("message", nil, msgBytes)
	}

	return nil
}

// 持久化消息
func (*wsDaoChat) Add(in *pb.MessageItem, auth *pb.Auth, object_id int64, message_type, is_read int32, content string) (int64, error) {
	sql_str := `insert into im_message(app_id, object_id, sender_id, receiver_type, receiver_id, type,
		content, is_read, conversation_id, help, conversation_message) 
		values(?,?,?,?,?,?,?,?,?,?,?)`
	result, err := db.DBCli.Exec(sql_str, auth.AppId, object_id, in.SenderId, in.ReceiverType, in.ReceiverId,
		message_type, content, is_read, in.ConversationId, in.Help, in.ConversationMessage)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return lastInsertID, nil
}

// 添加会话
func (*wsDaoChat) AddConversation(in *pb.ConversationItem, auth *pb.Auth, userid int64) (int64, error) {
	helpwhere := ""
	if in.Help == 1 {
		helpwhere = ", help=0"
	}
	sql_str := `insert into im_conversation(app_id, user_id, receiver_id, receiver_type, sender_id, sender_name, message) 
		values(?,?,?,?,?,?,?)
		ON DUPLICATE KEY UPDATE disturb=1, top=1, sender_id=?, sender_name=?, status=0, message=?, update_time=now()` + helpwhere
	result, err := db.DBCli.Exec(sql_str, auth.AppId, userid, in.ReceiverId, in.ReceiverType, in.SenderId, in.SenderName,
		in.Messagenewcontent, in.SenderId, in.SenderName, in.Messagenewcontent)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return lastInsertID, nil
}

// 添加会话
func (*wsDaoChat) AddConversationDao(in *pb.AddConversationReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	selfConversation := pb.ConversationItem{
		ReceiverType: in.Type,
		ReceiverId:   in.UserId,
	}
	id, err := WsDaoChat.AddConversation(&selfConversation, auth, auth.UserId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	msgBytes, err := proto.Marshal(&pb.AddConversationResp{ConversationId: id})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("addconversation", nil, msgBytes)
	return conn, nil
}

// 更新群组用户消息索引
func (*wsDaoChat) Upindex(appId, is_read, group_id, user_id int64) error {
	sql_str := `update group_user set is_read = ? where app_id = ? and group_id = ? and user_id = ? and examine = 1 and status = 0`
	_, err := db.DBCli.Exec(sql_str, is_read, appId, group_id, user_id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// 查询用户所有的在线设备
func (*wsDaoChat) ListOnlineByUserId(appId, userId int64) ([]models.Device, error) {
	result, err := cache.DeviceCache.Get(appId, userId)
	if err == nil {
		return result, nil
	}

	sql_str := `select id, user_id, type, brand, model, system_version, status, create_time, update_time, identification 
	from im_device where app_id = ? and user_id = ? and status = 1`
	rows, err := db.DBCli.Query(sql_str, appId, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var devices []models.Device
	for rows.Next() {
		device := new(models.Device)
		err = rows.Scan(&device.Id, &device.UserId, &device.Type, &device.Brand, &device.Model, &device.SystemVersion,
			&device.Status, &device.CreateTime, &device.UpdateTime, &device.Identification)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		devices = append(devices, *device)
	}

	err = cache.DeviceCache.Set(appId, userId, devices)
	if err != nil {
		log.Error(err)
	}

	return devices, nil
}

// 查找群组内的所有群成员
func (*wsDaoChat) GetGroupUsers(in *pb.GroupUserReq, auth *pb.Auth) (*pb.GroupUserResp, error) {
	result, err := cache.GroupUserCache.Get(auth.AppId, in.GroupId)
	if err == nil {
		return result, nil
	}

	sql_str := `select user_id from im_group_user where app_id = ? and group_id = ? and status = 0 and examine = 1`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, in.GroupId)
	if err != nil {
		return nil, err
	}
	var GroupUser []*pb.GroupUser
	for rows.Next() {
		var gu pb.GroupUser
		err := rows.Scan(&gu.UserId)
		if err != nil {
			return nil, err
		}
		GroupUser = append(GroupUser, &gu)
	}

	err = cache.GroupUserCache.Set(auth.AppId, in.GroupId, &pb.GroupUserResp{GroupUser: GroupUser})
	if err != nil {
		log.Error(err)
	}

	return &pb.GroupUserResp{GroupUser: GroupUser}, nil
}

// 查找群组内的管理员和群主
func (*wsDaoChat) GetGroupUsersAdmin(in *pb.GroupUserReq, auth *pb.Auth) (*pb.GroupUserResp, error) {
	sql_str := `select user_id from im_group_user where app_id = ? and group_id = ? and status = 0 and examine = 1 and (type = 2 or type = 3)`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, in.GroupId)
	if err != nil {
		return nil, err
	}
	var GroupUser []*pb.GroupUser
	for rows.Next() {
		var gu pb.GroupUser
		err := rows.Scan(&gu.UserId)
		if err != nil {
			return nil, err
		}
		GroupUser = append(GroupUser, &gu)
	}

	return &pb.GroupUserResp{GroupUser: GroupUser}, nil
}

// 回执
func (*wsDaoChat) MessageRead(in *pb.MessageRead, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	sql_str := `update im_message set is_read = 1 where app_id = ? and id = ?`
	_, err := db.DBCli.Exec(sql_str, auth.AppId, in.MessageId)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	return conn, nil
}

// 发送好友/加群请求
func (*wsDaoChat) AddExamine(in *pb.AddExamine, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	if in.Type == 1 {
		if auth.UserId == in.Fid {
			var Errlevel error = errors.New("不能加自己")
			return conn, Errlevel
		}
		// 好友请求
		sql_str := `insert into im_friend(user_id, fid, remark, way, examinetext, examine_time, app_id) 
			values(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE 
			is_read=0, state=1, examine=0, create_time=now(), way=?, examinetext=?, examine_time=now(), remark=?`
		_, err := db.DBCli.Exec(sql_str, auth.UserId, in.Fid, in.Remark, in.Way, in.Examinetext, time.Now(),
			auth.AppId, in.Way, in.Examinetext, in.Remark)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		// 提醒接收者
		devices, err := WsDaoChat.ListOnlineByUserId(auth.AppId, in.Fid)
		if err != nil {
			log.Error(err)
			return conn, err
		}
		for i := range devices {
			conn := internal_ws.Load(devices[i].Id)
			if conn == nil {
				log.Error("设备错误：", devices[i].Id)
				continue
			}

			conn.WriteMSG("examine", nil, nil)
		}
	} else if in.Type == 2 {
		// 加群请求
		sql_str := `insert into im_group_user(app_id, user_id, group_id, label, way, examinetext, examine_time, examine) 
			values(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE 
			is_read=0, status=0, examine=?, create_time=now(), way=?, examinetext=?, examine_time=now(), label=?`
		_, err := db.DBCli.Exec(sql_str, auth.UserId, in.Fid, in.Remark, in.Way, in.Examinetext, time.Now(), in.Examine,
			in.Examine, in.Way, in.Examinetext, in.Remark)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		if in.Examine == 1 {
			cache.GroupUserCache.Del(auth.AppId, in.Fid)
		}

		// 提醒接收者
		GroupUserReq := pb.GroupUserReq{
			GroupId: in.Fid,
		}
		GroupUserResp, err := WsDaoChat.GetGroupUsersAdmin(&GroupUserReq, auth)
		if err != nil {
			log.Error(err)
			return conn, err
		}
		for _, GroupUser := range GroupUserResp.GroupUser {
			devices, err := WsDaoChat.ListOnlineByUserId(auth.AppId, GroupUser.UserId)
			if err != nil {
				log.Error(err)
				return conn, err
			}
			for i := range devices {
				conn := internal_ws.Load(devices[i].Id)
				if conn == nil {
					log.Error("设备错误：", devices[i].Id)
					continue
				}

				conn.WriteMSG("examine", nil, nil)
			}
		}
	}

	conn.WriteMSG("addexamine", nil, nil)
	return conn, nil
}

// 获取好友/加群请求
func (*wsDaoChat) GetExamine(in *pb.GetExamineReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	limit := in.Limit
	if in.Limit == 0 {
		limit = 15
	}
	where_friend := "("
	where_groupuser := "("

	sql_str := `
	select id, fid, user_id, remark, examinetext, UNIX_TIMESTAMP(examine_time) as examine_time, examine, 1 as type
	from im_friend 
	where app_id = ? and (fid = ? or user_id = ?) and state = 1
	UNION
	select id, group_id as fid, user_id, '' as remark, examinetext, UNIX_TIMESTAMP(examine_time) as examine, examine, 2 as type
	from im_group_user 
	where app_id = ? and status = 0 and 
	group_id in (select group_id from im_group_user where user_id = ? and status = 0 and app_id = ? and (type = 2 or type = 3)) 
	order by examine_time desc limit ?,?
	`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, auth.UserId, auth.UserId, auth.AppId, auth.UserId, auth.AppId,
		in.Offset, limit)
	if err != nil {
		log.Error(err)
		return conn, err
	}
	var GetExamine []*pb.GetExamine
	for rows.Next() {
		var gex pb.GetExamine
		err := rows.Scan(&gex.Id, &gex.Fid, &gex.UserId, &gex.Remark, &gex.Examinetext, &gex.ExamineTime, &gex.Examine, &gex.Type)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		// 获取昵称，头像，备注
		sql_str := `
		select nickname, avatar_url
		from im_user 
		where app_id = ? and user_id = ?
		`
		err = db.DBCli.QueryRow(sql_str, auth.AppId, gex.UserId).
			Scan(&gex.Nickname, &gex.AvatarUrl)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		if gex.Type == 1 {
			idstr := strconv.FormatInt(gex.Id, 10) + ","
			where_friend += idstr
		} else if gex.Type == 2 {
			idstr := strconv.FormatInt(gex.Id, 10) + ","
			where_groupuser += idstr
		}

		GetExamine = append(GetExamine, &gex)
	}

	// 设置已读
	if len(where_friend) > 2 {
		where_friend = where_friend[0 : len(where_friend)-1]
	}
	if len(where_groupuser) > 2 {
		where_groupuser = where_groupuser[0 : len(where_groupuser)-1]
	}
	where_friend += ")"
	where_groupuser += ")"
	if where_friend != "()" {
		_, err := db.DBCli.Exec("update im_friend set is_read = 1 where fid = ? and id in "+where_friend,
			auth.UserId)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	}
	if where_groupuser != "()" {
		_, err := db.DBCli.Exec("update im_group_user set friend_read = 1 where user_id != ? and id in "+where_groupuser, auth.UserId)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	}

	msgBytes, err := proto.Marshal(&pb.GetExamines{Examine: GetExamine})
	if err != nil {
		return conn, err
	}
	conn.WriteMSG("getexamine", nil, msgBytes)
	return conn, nil
}

// 处理好友/加群请求
func (*wsDaoChat) UpExamine(in *pb.UpExamineReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	if in.Type != 1 && in.Type != 2 {
		var Errlevel error = errors.New("请求类型错误")
		return conn, Errlevel
	}

	if in.Type == 1 {
		// 好友请求
		// 权限验证
		if in.Fid != auth.UserId {
			var Errlevel error = errors.New("权限不足")
			return conn, Errlevel
		}
		sql_str := `update im_friend set examine = ? where app_id = ? and user_id = ? and fid = ? and state = 1`
		_, err := db.DBCli.Exec(sql_str, in.Examine, auth.AppId, in.UserId, in.Fid)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	} else if in.Type == 2 {
		// 加群请求
		// 只有管理员和群主可以处理加群请求
		istype, err := WsDaoChat.TypeByUser(auth, in.Fid)
		if err != nil {
			log.Error(err)
			return conn, err
		}
		if istype != 2 && istype != 3 {
			var Errlevel error = errors.New("权限不足")
			return conn, Errlevel
		}

		sql_str := `update im_group_user set examine = ? where app_id = ? and user_id = ? and group_id = ? and status = 0`
		_, err = db.DBCli.Exec(sql_str, in.Examine, auth.AppId, in.UserId, in.Fid)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		if in.Examine == 1 {
			cache.GroupUserCache.Del(auth.AppId, in.Fid)
		}
	}

	// 创建会话
	selfConversation := pb.ConversationItem{
		ReceiverType: in.Type,
		ReceiverId:   in.Fid,
	}
	// 加好友成功创建会话
	if in.Examine == 1 && in.Type == 1 {
		// 创建发送者会话
		_, err := WsDaoChat.AddConversation(&selfConversation, auth, in.UserId)
		if err != nil {
			log.Error(err)
			return conn, err
		}

		// 创建接收者会话
		selfConversation.ReceiverId = in.UserId
		_, err = WsDaoChat.AddConversation(&selfConversation, auth, in.Fid)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	}
	// 加群成功创建会话
	if in.Examine == 1 && in.Type == 2 {
		_, err := WsDaoChat.AddConversation(&selfConversation, auth, in.UserId)
		if err != nil {
			log.Error(err)
			return conn, err
		}
	}

	conn.WriteMSG("upexamine", nil, nil)
	return conn, nil
}

// 查询群成员权限
func (*wsDaoChat) TypeByUser(auth *pb.Auth, group_id int64) (int32, error) {
	var istype int32
	sql_str := `
	select type 
	from im_group_user 
	where app_id = ? and user_id = ? and group_id = ? and status = 0
	`
	err := db.DBCli.QueryRow(sql_str, auth.AppId, auth.UserId, group_id).
		Scan(&istype)
	if err != nil {
		log.Error(err)
		return istype, err
	}

	return istype, nil
}

// 好友/加群请求未读数
func (*wsDaoChat) ExamineReadNum(auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	var one_num int64
	var all_num int64

	// 好友请求未读数
	sql_str := `select count(*) from im_friend as f 
	left join im_group_user as g on g.user_id = f.user_id and g.app_id = ? and g.status = 0 and g.friend_read = 0 and g.examine = 0
	where f.app_id = ? and f.fid = ? and f.state = 1 and f.examine = 0 and f.is_read = 0 and f.state = 1`
	err := db.DBCli.QueryRow(sql_str, auth.AppId, auth.AppId, auth.UserId).Scan(&one_num)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	// 群聊请求未读数
	sql_str = `select count(*) from im_group_user 
	where app_id = ? and group_id in (select group_id from im_group_user where user_id = ? and status = 0 and app_id = ? 
	and (type = 2 or type = 3)) and status = 0 and friend_read = 0 and user_id != ?`
	err = db.DBCli.QueryRow(sql_str, auth.AppId, auth.UserId, auth.AppId, auth.UserId).Scan(&all_num)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	num := one_num + all_num

	msgBytes, err := proto.Marshal(&pb.ReadNumResp{Num: num})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("examinereadnum", nil, msgBytes)
	return conn, nil
}

// 聊天未读数
func (*wsDaoChat) ChatReadNum(auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	var one_num int64
	var all_num int64

	// 单聊未读数
	sql_str := `
	select count(*) 
	from im_message as m 
	left join im_conversation as c on c.id = m.conversation_id
	where m.app_id = ? and m.receiver_type = 1 and m.object_id = ? and m.is_read = 0 
	and c.disturb = 1 and m.status != 2 and m.status != 1`
	err := db.DBCli.QueryRow(sql_str, auth.AppId, auth.UserId).
		Scan(&one_num)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return conn, err
	}

	// 群聊未读数
	sql_str = `
	select count(*)
	from im_message as m 
	left join im_group_user as gu on gu.user_id = ? and gu.group_id = m.receiver_id 
	left join im_conversation as c on c.user_id = ? and c.receiver_type = 2 
	where m.app_id = ? and m.receiver_type = 2 and c.status = 0 and gu.status = 0 and 
	m.receiver_id in (select group_id from im_group_user as g where g.user_id = ? and g.app_id = m.app_id) 
	and m.status != 2 and m.status != 1 and m.id > gu.is_read and gu.examine = 1 and gu.status = 0 group by m.id`
	err = db.DBCli.QueryRow(sql_str, auth.UserId, auth.UserId, auth.AppId, auth.UserId).
		Scan(&all_num)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return conn, err
	}

	num := one_num + all_num

	msgBytes, err := proto.Marshal(&pb.ReadNumResp{Num: num})
	if err != nil {
		log.Error(err)
		return conn, err
	}
	conn.WriteMSG("chatreadnum", nil, msgBytes)
	return conn, nil
}
