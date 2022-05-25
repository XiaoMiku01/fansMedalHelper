package dto

// BiliBaseResp is basic response body of all bilibili API
type BiliBaseResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

// BiliDataResp only check status and data
type BiliDataResp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	TTL     int         `json:"ttl,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data"`
}

type MedalList struct {
	Medal struct {
		UID              int    `json:"uid"`
		TargetID         int    `json:"target_id"`
		TargetName       string `json:"target_name"`
		MedalID          int    `json:"medal_id"`
		Level            int    `json:"level"`
		MedalName        string `json:"medal_name"`
		MedalColor       int    `json:"medal_color"`
		Intimacy         int    `json:"intimacy"`
		NextIntimacy     int    `json:"next_intimacy"`
		DayLimit         int    `json:"day_limit"`
		TodayFeed        int    `json:"today_feed"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorBorder int    `json:"medal_color_border"`
		IsLighted        int    `json:"is_lighted"`
		GuardLevel       int    `json:"guard_level"`
		WearingStatus    int    `json:"wearing_status"`
		MedalIconID      int    `json:"medal_icon_id"`
		MedalIconURL     string `json:"medal_icon_url"`
		GuardIcon        string `json:"guard_icon"`
		HonorIcon        string `json:"honor_icon"`
		CanDelete        bool   `json:"can_delete"`
	} `json:"medal"`
	AnchorInfo struct {
		NickName string `json:"nick_name"`
		Avatar   string `json:"avatar"`
		Verify   int    `json:"verify"`
	} `json:"anchor_info"`
	Superscript interface{} `json:"superscript"`
	RoomInfo    struct {
		RoomID       int    `json:"room_id"`
		LivingStatus int    `json:"living_status"`
		URL          string `json:"url"`
	} `json:"room_info"`
}

// BiliMedalResp obtain the response with all medal info
type BiliMedalResp struct {
	BiliBaseResp
	Data struct {
		List        []MedalList `json:"list"`
		SpecialList []MedalList `json:"special_list"`
		BottomBar   interface{} `json:"bottom_bar"`
		PageInfo    struct {
			Number          int  `json:"number"`
			CurrentPage     int  `json:"current_page"`
			HasMore         bool `json:"has_more"`
			NextPage        int  `json:"next_page"`
			NextLightStatus int  `json:"next_light_status"`
			TotalPage       int  `json:"total_page"`
		} `json:"page_info"`
		TotalNumber int `json:"total_number"`
		HasMedal    int `json:"has_medal"`
	} `json:"data"`
}
