package manager

import (
	"MedalHelper/dto"
	"MedalHelper/util"
	"encoding/json"
	"fmt"
	"math/rand"
)

func LoginVerify(accessKey string) (dto.BiliAccountResp, error) {
	rawUrl := "http://app.bilibili.com/x/v2/account/mine"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	var resp dto.BiliAccountResp
	body, err := Get(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("LoginVerify error: %v, data: %v", err, data)
		return resp, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliAccountResp error: %v, raw data: %v", err, body)
		return resp, err
	}
	return resp, nil
}

func SignIn(accessKey string) (dto.BiliDataResp, error) {
	rawUrl := "http://api.live.bilibili.com/rc/v1/Sign/doSign"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	var resp dto.BiliDataResp
	body, err := Get(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("SignIn error: %v, data: %v", err, data)
		return resp, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliDataResp error: %v, raw data: %v", err, body)
		return resp, err
	}
	return resp, nil
}

func GetUserInfo(accessKey string) (dto.BiliLiveUserInfo, error) {
	rawUrl := "http://api.live.bilibili.com/xlive/app-ucenter/v1/user/get_user_info"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	body, err := Get(rawUrl, util.Map2Params(data))
	var resp dto.BiliLiveUserInfo
	if err != nil {
		util.Error("GetUserInfo error: %v, data: %v", err, data)
		return resp, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliLiveUserInfo error: %v, raw data: %v", err, body)
		return resp, err
	}
	return resp, nil
}

func GetFansMedalAndRoomID(accessKey string) []dto.MedalList {
	medals := make([]dto.MedalList, 0, 20)
	page := 1
	for {
		rawUrl := "http://api.live.bilibili.com/xlive/app-ucenter/v1/fansMedal/panel"
		data := map[string]string{
			"access_key": accessKey,
			"actionKey":  "appkey",
			"appkey":     util.AppKey,
			"ts":         util.GetTimestamp(),
			"page":       fmt.Sprint(page),
			"page_size":  "100",
		}
		util.Signature(&data)
		body, err := Get(rawUrl, util.Map2Params(data))
		if err != nil {
			util.Error("GetFansMedalAndRoomID error: %v, data: %v", err, data)
			return medals
		}
		var resp dto.BiliMedalResp
		if err = json.Unmarshal(body, &resp); err != nil {
			util.Error("Unmarshal BiliMedalResp error: %v, raw data: %v", err, body)
			return medals
		}
		medals = append(medals, resp.Data.SpecialList...)
		medals = append(medals, resp.Data.List...)
		if len(resp.Data.List) == 0 {
			break
		}
		page++
	}
	return medals
}

func LikeInteract(accessKey string, roomId int) bool {
	rawUrl := "http://api.live.bilibili.com/xlive/web-ucenter/v1/interact/likeInteract"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
		"roomid":     fmt.Sprint(roomId),
	}
	util.Signature(&data)
	body, err := Post(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("LikeInteract error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

func ShareRoom(accessKey string, roomId int) bool {
	rawUrl := "http://api.live.bilibili.com/xlive/app-room/v1/index/TrigerInteract"
	data := map[string]string{
		"access_key":    accessKey,
		"actionKey":     "appkey",
		"appkey":        util.AppKey,
		"interact_type": "3",
		"ts":            util.GetTimestamp(),
		"roomid":        fmt.Sprint(roomId),
	}
	util.Signature(&data)
	body, err := Post(rawUrl, util.Map2Params(data))
	if err != nil {
		util.Error("ShareRoom error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

func SendDanmaku(accessKey string, roomId int) bool {
	rawUrl := "http://api.live.bilibili.com/xlive/app-room/v1/dM/sendmsg"
	params := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	data := map[string]string{
		"cid":      fmt.Sprint(roomId),
		"msg":      util.GlobalConfig.Danmaku[rand.Intn(len(util.GlobalConfig.Danmaku))],
		"rnd":      util.GetTimestamp(),
		"color":    "16777215",
		"fontsize": "25",
	}
	util.Signature(&params)
	body, err := PostWithParam(rawUrl, util.Map2Params(params), util.Map2Params(data))
	if err != nil {
		util.Error("GetFansMedalAndRoomID error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		util.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}
