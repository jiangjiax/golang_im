package api_ws

import (
	"fmt"
	"golang_im/internal/internal_ws"
)

func Controller_init() {
	// 基础
	logic_Controller()
	// 单聊、群聊
	chat_Controller()
	// 群组
	group_Controller()
	// 好友
	friend_Controller()
	// 朋友圈
	circle_Controller()

	fmt.Println("controller start")
}

// 基础
func logic_Controller() {
	// 获取用户信息 [已测试]
	internal_ws.Controllers["getuserinfo"] = WsServiceLogic.GetUserInfo
	// 心跳 [已测试]
	internal_ws.Controllers["heartbeat"] = WsServiceLogic.Heartbeat
}

// 单聊，群聊，好友/加群请求
func chat_Controller() {
	// 获取消息会话列表（及未读数，分页） [已测试]
	internal_ws.Controllers["getconversationlist"] = WsServiceChat.GetConversationList
	// 根据会话id获取消息 [已测试]
	internal_ws.Controllers["sync"] = WsServiceChat.Sync
	// 发送消息 [已测试]
	internal_ws.Controllers["sendmessage"] = WsServiceChat.SendMessage
	// 消息回执 [已测试]
	internal_ws.Controllers["messageread"] = WsServiceChat.MessageRead
	// 发送好友/加群请求 [已测试]
	internal_ws.Controllers["addexamine"] = WsServiceChat.AddExamine
	// 获取好友/加群请求 [已测试]
	internal_ws.Controllers["getexamine"] = WsServiceChat.GetExamine
	// 处理好友/加群请求 [已测试]
	internal_ws.Controllers["upexamine"] = WsServiceChat.UpExamine
	// 好友/加群请求未读数 [已测试]
	internal_ws.Controllers["examinereadnum"] = WsServiceChat.ExamineReadNum
	// 聊天未读数 [已测试]
	internal_ws.Controllers["chatreadnum"] = WsServiceChat.ChatReadNum
	// 发起聊天 [已测试]
	internal_ws.Controllers["addconversation"] = WsServiceChat.AddConversation
	// 会话置顶
	internal_ws.Controllers["upconversation"] = WsServiceChat.UpConversation
	// 会话免打扰
	internal_ws.Controllers["disturbconversation"] = WsServiceChat.DisturbConversation

	// [以下接口下次补充]
	// 聊天信息删除
	// 聊天信息撤回
	// 获取客服
	// 处理好友/加群请求增加新用户加群提示
}

// 群组
func group_Controller() {

	// [以下接口下次补充]
	// 创建群组
	// 更新群组（群头像，群名）
	// 删除群组
	// 获取用户加入的所有群组
	// 转让群主
	// 更新群组成员信息（改名）
	// 更新群组成员信息（改权限）
	// 获取群组中所有成员的信息
	// 获取群组中所有管理员和群主的信息
	// 获取群组中所有除了管理员和群主的信息
	// 删除群组成员
	// 批量删除群组成员
	// 批量添加群组成员
	// 批量更新群组成员信息（权限）
}

// 好友
func friend_Controller() {

	// [以下接口下次补充]
	// 获取好友列表
	// 搜索好友
	// 更新好友备注
	// 删除好友
	// 拉黑和取消拉黑
}

// 朋友圈
func circle_Controller() {
	// 发朋友圈 [已测试]
	internal_ws.Controllers["addtrend"] = WsServiceChat.AddTrend
	// 获取朋友圈列表（包括内容、点赞人列表、评论列表，下拉加载） [已测试]
	internal_ws.Controllers["gettrends"] = WsServiceChat.GetTrends
	// 点赞与取消点赞 [已测试]
	internal_ws.Controllers["thumb"] = WsServiceChat.Thumb
	// 评论与回复 [已测试]
	internal_ws.Controllers["addtrendscomment"] = WsServiceChat.AddTrendsComment

	// [以下接口下次补充]
	// 获取一条朋友圈详情
	// 朋友圈未读数
	// 获取未读的评论与点赞列表
	// 编辑朋友圈
	// 删除朋友圈
	// 删除评论
	// 获取某人的动态
}
