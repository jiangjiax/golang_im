package api_ws

import (
	"errors"
	"golang_im/internal/internal_ws"
	"golang_im/pkg/db"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/log"
	"golang_im/pkg/pb"
	"strings"

	"github.com/gogo/protobuf/proto"
)

type wsDaoCircle struct{}

var WsDaoCircle = new(wsDaoCircle)

// 发朋友圈
func (*wsDaoCircle) AddTrend(in *pb.AddTrend, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	// 可以限制图片数
	imgs := strings.Split(in.Imgs, ",")
	if len(imgs) > 9 {
		var Errlevel error = errors.New("图片不能超过9张")
		return conn, Errlevel
	}

	sql_str := `insert into im_trends(app_id, user_id, writing, imgs, videos, to_user_ids) values(?,?,?,?,?,?)`
	_, err := db.DBCli.Exec(sql_str, auth.AppId, auth.UserId, in.Writing, in.Imgs, in.Videos, in.ToUserIds)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	conn.WriteMSG("addtrend", nil, nil)
	return conn, nil
}

// 获取朋友圈列表（包括内容、点赞人列表、评论列表，下拉加载）
func (*wsDaoCircle) GetTrends(in *pb.GetTrendsReq, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	limit := in.Limit
	if in.Limit == 0 {
		limit = 6
	}
	sql_str := `select t.user_id, t.writing, t.imgs, t.videos, t.id, UNIX_TIMESTAMP(t.create_time), 
	UNIX_TIMESTAMP(t.update_time), u.nickname, u.avatar_url
	from im_trends as t
	left join im_user as u on u.user_id = t.user_id
	where t.app_id = ? and t.status = 1 and u.status = 1 limit ?,?`
	rows, err := db.DBCli.Query(sql_str, auth.AppId, in.Offset, limit)
	if err != nil {
		return conn, err
	}
	var Trends []*pb.Trends
	for rows.Next() {
		var trend pb.Trends
		err := rows.Scan(&trend.UserId, &trend.Writing, &trend.Imgs, &trend.Videos, &trend.Id, &trend.CreateTime,
			&trend.UpdateTime, &trend.Nickname, &trend.AvatarUrl)
		if err != nil {
			return conn, err
		}

		// 点赞数
		sql_str := `select count(*) from im_trends_handle where app_id = ? and trends_id = ? and status = 0 and istype = 1`
		err = db.DBCli.QueryRow(sql_str, auth.AppId, trend.Id).Scan(&trend.ThumbNum)
		if err != nil {
			return conn, err
		}

		// 评论数
		sql_str = `select count(*) from im_trends_comment where app_id = ? and trends_id = ? and status = 0`
		err = db.DBCli.QueryRow(sql_str, auth.AppId, trend.Id).Scan(&trend.CommentNum)
		if err != nil {
			return conn, err
		}

		// 朋友圈评论（显示10条）
		sql_str = `select t.id, t.trends_id, t.reply_id, t.comment_id, t.user_id, t.writing, t.istype, 
		UNIX_TIMESTAMP(t.create_time), UNIX_TIMESTAMP(t.update_time), u.nickname, u.avatar_url
		from im_trends_comment as t
		left join im_user as u on u.user_id = t.user_id
		where t.app_id = ? and t.status = 0 and u.status = 1 and t.trends_id = ? limit 10`
		rows, err := db.DBCli.Query(sql_str, auth.AppId, trend.Id)
		if err != nil {
			return conn, err
		}
		var TrendsComments []*pb.TrendsComment
		for rows.Next() {
			var tc pb.TrendsComment
			err := rows.Scan(&tc.Id, &tc.TrendsId, &tc.ReplyId, &tc.CommentId, &tc.UserId, &tc.Writing, &tc.Istype,
				&tc.CreateTime, &tc.UpdateTime, &tc.Nickname, &tc.AvatarUrl)
			if err != nil {
				return conn, err
			}

			// 被回复人的昵称
			if tc.Istype == 2 {
				sql_str := `select nickname, avatar_url from im_user where im_user = ?`
				err := db.DBCli.QueryRow(sql_str, tc.ReplyId).
					Scan(&tc.ReplyNickname, &tc.ReplyAvatarUrl)
				if err != nil {
					return conn, err
				}
			}

			TrendsComments = append(TrendsComments, &tc)
		}

		// 朋友圈点赞（显示15条）
		sql_str = `
		select t.id, t.trends_id, t.reply_id, t.user_id, UNIX_TIMESTAMP(t.create_time), UNIX_TIMESTAMP(t.update_time), 
		u.nickname, u.avatar_url
		from im_trends_handle as t
		left join im_user as u on u.user_id = t.user_id
		where t.app_id = ? and t.status = 0 and u.status = 1 and t.trends_id = ? and istype = 1 limit 15`
		rows, err = db.DBCli.Query(sql_str, auth.AppId, trend.Id)
		if err != nil {
			return conn, err
		}
		var TrendThumbs []*pb.TrendThumb
		for rows.Next() {
			var tt pb.TrendThumb
			err := rows.Scan(&tt.Id, &tt.TrendsId, &tt.ReplyId, &tt.UserId, &tt.CreateTime, &tt.UpdateTime,
				&tt.Nickname, &tt.AvatarUrl)
			if err != nil {
				return conn, err
			}

			TrendThumbs = append(TrendThumbs, &tt)
		}

		trend.TrendsComment = TrendsComments
		trend.TrendsThumb = TrendThumbs
		Trends = append(Trends, &trend)
	}

	msgBytes, err := proto.Marshal(&pb.GetTrendsResp{Trends: Trends})
	if err != nil {
		return conn, err
	}
	conn.WriteMSG("gettrends", nil, msgBytes)
	return conn, nil
}

// 点赞与取消点赞
func (*wsDaoCircle) Thumb(in *pb.Thumb, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	sql_str := `insert into im_trends_handle(trends_id, user_id, app_id, istype, status) values(?,?,?,?,?)
	ON DUPLICATE KEY UPDATE isread = 0, status = ?`
	_, err := db.DBCli.Exec(sql_str, in.TrendsId, auth.UserId, auth.AppId, 1, in.Type, in.Type)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	conn.WriteMSG("thumb", nil, nil)
	return conn, nil
}

// 评论与回复
func (*wsDaoCircle) AddTrendsComment(in *pb.AddTrendsComment, auth *pb.Auth) (*internal_ws.WSConn, error) {
	conn := internal_ws.Load(auth.DeviceId)
	if conn == nil {
		return conn, gerrors.ErrUnDeviceid
	}

	var reply_id int64
	// 查询回复人id（看istype，如果是评论动态就存这条动态的发送者id，如果是回复就存被回复的人的id）
	if in.Istype == 1 {
		// 评论动态
		sql_str := `select user_id from im_trends where app_id = ? and id = ? and status = 1`
		err := db.DBCli.QueryRow(sql_str, auth.AppId, in.TrendsId).Scan(&reply_id)
		if err != nil {
			return conn, err
		}
	} else if in.Istype == 2 {
		// 回复评论
		sql_str := `select user_id from im_trends_comment where app_id = ? and trends_id = ? and status = 1`
		err := db.DBCli.QueryRow(sql_str, auth.AppId, in.CommentId).Scan(&reply_id)
		if err != nil {
			return conn, err
		}
	}

	sql_str := `insert into im_trends_comment(trends_id, reply_id, comment_id, app_id, user_id, writing, istype) values(?,?,?,?,?,?,?)`
	_, err := db.DBCli.Exec(sql_str, in.TrendsId, reply_id, in.CommentId, auth.AppId, auth.UserId, in.Writing, in.Istype)
	if err != nil {
		log.Error(err)
		return conn, err
	}

	conn.WriteMSG("addtrendscomment", nil, nil)
	return conn, nil
}
