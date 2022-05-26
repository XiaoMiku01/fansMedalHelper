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

// BiliAccountResp represent account mine response
type BiliAccountResp struct {
	BiliBaseResp
	Data struct {
		Mid               int     `json:"mid"`
		Name              string  `json:"name"`
		ShowNameGuide     bool    `json:"show_name_guide"`
		Face              string  `json:"face"`
		ShowFaceGuide     bool    `json:"show_face_guide"`
		Coin              float64 `json:"coin"`
		Bcoin             int     `json:"bcoin"`
		Sex               int     `json:"sex"`
		Rank              int     `json:"rank"`
		Silence           int     `json:"silence"`
		ShowVideoup       int     `json:"show_videoup"`
		ShowCreative      int     `json:"show_creative"`
		Level             int     `json:"level"`
		VipType           int     `json:"vip_type"`
		AudioType         int     `json:"audio_type"`
		Dynamic           int     `json:"dynamic"`
		Following         int     `json:"following"`
		Follower          int     `json:"follower"`
		NewFollowers      int     `json:"new_followers"`
		NewFollowersRtime int     `json:"new_followers_rtime"`
		OfficialVerify    struct {
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"official_verify"`
		Vip struct {
			Type       int   `json:"type"`
			Status     int   `json:"status"`
			DueDate    int64 `json:"due_date"`
			VipPayType int   `json:"vip_pay_type"`
			ThemeType  int   `json:"theme_type"`
			Label      struct {
				Path        string `json:"path"`
				Text        string `json:"text"`
				LabelTheme  string `json:"label_theme"`
				TextColor   string `json:"text_color"`
				BgStyle     int    `json:"bg_style"`
				BgColor     string `json:"bg_color"`
				BorderColor string `json:"border_color"`
			} `json:"label"`
			AvatarSubscript    int    `json:"avatar_subscript"`
			NicknameColor      string `json:"nickname_color"`
			Role               int    `json:"role"`
			AvatarSubscriptURL string `json:"avatar_subscript_url"`
		} `json:"vip"`
		InRegAudit       int  `json:"in_reg_audit"`
		FirstLiveTime    int  `json:"first_live_time"`
		FaceNftNew       int  `json:"face_nft_new"`
		ShowNftFaceGuide bool `json:"show_nft_face_guide"`
		SeniorGate       struct {
			Identity   int    `json:"identity"`
			MemberText string `json:"member_text"`
		} `json:"senior_gate"`
	} `json:"data"`
}

// BiliUserInfo represent user live info
type BiliLiveUserInfo struct {
	BiliBaseResp
	Data struct {
		UID int `json:"uid"`
		Silver int `json:"silver"`
		Gold int `json:"gold"`
		Medal struct {
			MedalName string `json:"medal_name"`
			Level int `json:"level"`
			Color int `json:"color"`
			MedalIconURL string `json:"medal_icon_url"`
			TargetID int `json:"target_id"`
			MedalColorStart int `json:"medal_color_start"`
			MedalColorEnd int `json:"medal_color_end"`
			MedalColorBorder int `json:"medal_color_border"`
			IsLighted int `json:"is_lighted"`
			GuardLevel int `json:"guard_level"`
			GuardIcon string `json:"guard_icon"`
			HonorIcon string `json:"honor_icon"`
		} `json:"medal"`
		Vip struct {
			Vip int `json:"vip"`
			Svip int `json:"svip"`
			VipTime string `json:"vip_time"`
			SvipTime string `json:"svip_time"`
		} `json:"vip"`
		WearTitle struct {
			ID string `json:"id"`
			Img string `json:"img"`
		} `json:"wear_title"`
		Exp struct {
			Color int `json:"color"`
			UserLevel int `json:"user_level"`
			Cost int `json:"cost"`
			Unext int `json:"unext"`
			UserLevelCost int `json:"user_level_cost"`
		} `json:"exp"`
		RoomID int `json:"room_id"`
		VipViewStatus bool `json:"vip_view_status"`
		GuardCount int `json:"guard_count"`
	} `json:"data"`
}
