package service

import (
	"MedalHelper/dto"
	"MedalHelper/manager"
	"MedalHelper/util"
	"MedalHelper/util/log"
	"encoding/json"
	"fmt"
)

func LoginVerify(accessKey string, roomId int) (interface{}, error) {
	rawUrl := "http://app.bilibili.com/x/v2/account/mine"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	body, err := manager.Get(rawUrl, util.Map2Params(data))
	if err != nil {
		log.Error("LoginVerify error: %v, data: %v", err, data)
		return nil, err
	}
	var resp dto.BiliDataResp
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Error("Unmarshal BiliDataResp error: %v, raw data: %v", err, body)
		return nil, err
	}
	return resp.Data, nil
}

func SignIn(accessKey string, roomId int) (interface{}, error) {
	rawUrl := "http://api.live.bilibili.com/rc/v1/Sign/doSign"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	body, err := manager.Get(rawUrl, util.Map2Params(data))
	if err != nil {
		log.Error("SignIn error: %v, data: %v", err, data)
		return nil, err
	}
	var resp dto.BiliDataResp
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Error("Unmarshal BiliDataResp error: %v, raw data: %v", err, body)
		return nil, err
	}
	return resp.Data, nil
}

func GetUserInfo(accessKey string, roomId int) (interface{}, error) {
	rawUrl := "http://api.live.bilibili.com/xlive/app-ucenter/v1/user/get_user_info"
	data := map[string]string{
		"access_key": accessKey,
		"actionKey":  "appkey",
		"appkey":     util.AppKey,
		"ts":         util.GetTimestamp(),
	}
	util.Signature(&data)
	body, err := manager.Get(rawUrl, util.Map2Params(data))
	if err != nil {
		log.Error("SignIn error: %v, data: %v", err, data)
		return nil, err
	}
	var resp dto.BiliDataResp
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Error("Unmarshal BiliDataResp error: %v, raw data: %v", err, body)
		return nil, err
	}
	return resp.Data, nil
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
		body, err := manager.Get(rawUrl, util.Map2Params(data))
		if err != nil {
			log.Error("GetFansMedalAndRoomID error: %v, data: %v", err, data)
			return medals
		}
		var resp dto.BiliMedalResp
		if err = json.Unmarshal(body, &resp); err != nil {
			log.Error("Unmarshal BiliMedalResp error: %v, raw data: %v", err, body)
			return medals
		}
		medals = append(medals, resp.Data.SpecialList...)
		medals = append(medals, resp.Data.List...)
		if len(resp.Data.List) == 0 {
			break
		}
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
	body, err := manager.Post(rawUrl, util.Map2Params(data))
	if err != nil {
		log.Error("LikeInteract error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
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
	body, err := manager.Post(rawUrl, util.Map2Params(data))
	if err != nil {
		log.Error("ShareRoom error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}

// TODO: add detail later
func SendDanmaku(accessKey string, roomId int) bool {
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
	body, err := manager.Post(rawUrl, util.Map2Params(data))
	if err != nil {
		log.Error("GetFansMedalAndRoomID error: %v, data: %v", err, data)
	}
	var resp dto.BiliBaseResp
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Error("Unmarshal BiliBaseResp error: %v, raw data: %v", err, body)
	}
	return resp.Code == 0
}
