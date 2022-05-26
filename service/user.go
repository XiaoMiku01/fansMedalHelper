package service

import (
	"MedalHelper/dto"
	"MedalHelper/manager"
	"MedalHelper/util"
	"time"

	"github.com/sethvargo/go-retry"
)

type User struct {
	// 用户ID
	Uid string
	// 用户名称
	Name string
	// 是否登录
	isLogin bool

	// 登录凭证
	accessKey string
	// 被禁止的房间ID
	bannedUIDs []int
	// 用户所有勋章
	medals []dto.MedalList
	// 用户等级小于20的勋章
	medalsLow []dto.MedalList
	// 最大重试次数
	retryTimes int32
}

func NewUser(accessKey string, uids []int) User {
	return User{
		accessKey:  accessKey,
		bannedUIDs: uids,
		medals:     make([]dto.MedalList, 0, 10),
		medalsLow:  make([]dto.MedalList, 0, 10),
		retryTimes: 10,
	}
}

func (user *User) loginVerify() bool {
	data, err := manager.LoginVerify(user.accessKey)
	if err != nil || data.(map[string]string)["mid"] != "0" {
		user.isLogin = false
		return false
	}
	user.Uid = data.(map[string]string)["mid"]
	user.Name = data.(map[string]string)["name"]
	user.isLogin = true
	return true
}

func (user *User) signIn() error {
	signInfo, err := manager.SignIn(user.accessKey)
	if err != nil {
		return nil
	}
	signed := signInfo.(map[string]string)["hadSignDays"]
	all := signInfo.(map[string]string)["allDays"]
	util.Info("签到成功,本月签到次数: %s/%s", signed, all)

	userInfo, err := manager.GetUserInfo(user.accessKey)
	if err != nil {
		return nil
	}
	level := userInfo.(map[string]map[string]string)["exp"]["user_level"]
	unext := userInfo.(map[string]map[string]string)["exp"]["unext"]
	util.Info("当前用户UL等级: %s ,还差 %s 经验升级", level, unext)
	return nil
}

func (user *User) setMedals() {
	medals := manager.GetFansMedalAndRoomID(user.accessKey)
	for _, medal := range medals {
		if util.IntContain(user.bannedUIDs, medal.Medal.TargetID) != -1 {
			continue
		}
		if medal.RoomInfo.RoomID == 0 {
			continue
		}
		user.medals = append(user.medals, medal)
		if medal.Medal.Level < 20 {
			user.medalsLow = append(user.medalsLow, medal)
		}
	}
}

func (user *User) Init() {
	if user.loginVerify() {
		user.signIn()
		user.setMedals()
	} else {
		util.Error("用户登录失败, accessKey: %s", user.accessKey)
	}
}

func (user User) Start() {
	if user.isLogin {
		task := NewTask(user, []IAction{
			&Like{},
			&Share{},
			&Danmaku{},
		})
		task.Start()
	} else {
		util.Error("用户未登录, accessKey: %s", user.accessKey)
	}
}
