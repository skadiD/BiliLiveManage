package room

// WS_CMD_Danmaku
type danmuMsg struct {
	Contribution struct {
		Grade int `json:"grade"`
	} `json:"contribution"`
	CoreUserType int `json:"core_user_type"`
	Dmscore      int `json:"dmscore"`
	FansMedal    struct {
		AnchorRoomid     int    `json:"anchor_roomid"`
		GuardLevel       int    `json:"guard_level"`
		IconId           int    `json:"icon_id"`
		IsLighted        int    `json:"is_lighted"`
		MedalColor       int    `json:"medal_color"`
		MedalColorBorder int    `json:"medal_color_border"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		Score            int    `json:"score"`
		Special          string `json:"special"`
		TargetId         int    `json:"target_id"`
	} `json:"fans_medal"` // edit
	Identities    []int  `json:"identities"`
	IsSpread      int    `json:"is_spread"`
	MsgType       int    `json:"msg_type"`
	PrivilegeType int    `json:"privilege_type"`
	Roomid        int    `json:"roomid"`
	Score         int64  `json:"score"`
	SpreadDesc    string `json:"spread_desc"`
	SpreadInfo    string `json:"spread_info"`
	TailIcon      int    `json:"tail_icon"`
	Timestamp     int    `json:"timestamp"`
	TriggerTime   int64  `json:"trigger_time"`
	Uid           int    `json:"uid"`
	Uname         string `json:"uname"`
	UnameColor    string `json:"uname_color"`
}

type watchedChange struct {
	Num       int    `json:"num"`
	TextSmall string `json:"text_small"`
	TextLarge string `json:"text_large"`
}
type likeClick struct {
	ShowArea   int    `json:"show_area"`
	MsgType    int    `json:"msg_type"`
	LikeIcon   string `json:"like_icon"`
	Uid        int    `json:"uid"`
	LikeText   string `json:"like_text"`
	Uname      string `json:"uname"`
	UnameColor string `json:"uname_color"`
	Identities []int  `json:"identities"`
	FansMedal  struct {
		TargetId         int    `json:"target_id"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		MedalColor       int    `json:"medal_color"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorBorder int    `json:"medal_color_border"`
		IsLighted        int    `json:"is_lighted"`
		GuardLevel       int    `json:"guard_level"`
		Special          string `json:"special"`
		IconId           int    `json:"icon_id"`
		AnchorRoomid     int    `json:"anchor_roomid"`
		Score            int    `json:"score"`
	} `json:"fans_medal"`
	ContributionInfo struct {
		Grade int `json:"grade"`
	} `json:"contribution_info"`
	Dmscore int `json:"dmscore"`
}
type onlineRank struct {
	Uid        int    `json:"uid"`
	Face       string `json:"face"`
	Score      string `json:"score"`
	Uname      string `json:"uname"`
	Rank       int    `json:"rank"`
	GuardLevel int    `json:"guard_level"`
}

// 榜三通报
type onlineRankTop3 struct {
	Dmscore int `json:"dmscore"`
	List    []struct {
		Msg  string `json:"msg"`
		Rank int    `json:"rank"`
	} `json:"list"`
}

type interactRoom struct {
	Contribution struct {
		Grade int `json:"grade"`
	} `json:"contribution"`
	CoreUserType int `json:"core_user_type"`
	Dmscore      int `json:"dmscore"`
	FansMedal    struct {
		AnchorRoomid     int    `json:"anchor_roomid"`
		GuardLevel       int    `json:"guard_level"`
		IconId           int    `json:"icon_id"`
		IsLighted        int    `json:"is_lighted"`
		MedalColor       int    `json:"medal_color"`
		MedalColorBorder int    `json:"medal_color_border"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		Score            int    `json:"score"`
		Special          string `json:"special"`
		TargetId         int    `json:"target_id"`
	} `json:"fans_medal"`
	Identities    []int  `json:"identities"`
	IsSpread      int    `json:"is_spread"`
	MsgType       int    `json:"msg_type"`
	PrivilegeType int    `json:"privilege_type"`
	Roomid        int    `json:"roomid"`
	Score         int64  `json:"score"`
	SpreadDesc    string `json:"spread_desc"`
	SpreadInfo    string `json:"spread_info"`
	TailIcon      int    `json:"tail_icon"`
	Timestamp     int    `json:"timestamp"`
	TriggerTime   int64  `json:"trigger_time"`
	Uid           int    `json:"uid"`
	Uname         string `json:"uname"`
	UnameColor    string `json:"uname_color"`
}

// 禁言事件
type muteEvent struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Dmscore  int    `json:"dmscore"`
		Operator int    `json:"operator"`
		Uid      int    `json:"uid"`
		Uname    string `json:"uname"`
	} `json:"data"`
	Uid   string `json:"uid"`
	Uname string `json:"uname"`
}

type captainEntry struct {
	Id                   int    `json:"id"`
	Uid                  int    `json:"uid"`
	TargetId             int    `json:"target_id"`
	MockEffect           int    `json:"mock_effect"`
	Face                 string `json:"face"`
	PrivilegeType        int    `json:"privilege_type"`
	CopyWriting          string `json:"copy_writing"`
	CopyColor            string `json:"copy_color"`
	HighlightColor       string `json:"highlight_color"`
	Priority             int    `json:"priority"`
	BasemapUrl           string `json:"basemap_url"`
	ShowAvatar           int    `json:"show_avatar"`
	EffectiveTime        int    `json:"effective_time"`
	WebBasemapUrl        string `json:"web_basemap_url"`
	WebEffectiveTime     int    `json:"web_effective_time"`
	WeBEffectClose       int    `json:"web_effect_close"`
	WebCloseTime         int    `json:"web_close_time"`
	Business             int    `json:"business"`
	CopyWritingV2        string `json:"copy_writing_v2"`
	IconList             []int  `json:"icon_list"`
	MaxDelayTime         int    `json:"max_delay_time"`
	TriggerTime          int64  `json:"trigger_time"`
	Identities           int    `json:"identities"`
	EffectSilentTime     int    `json:"effect_silent_time"`
	EffectiveTimeNew     int    `json:"effective_time_new"`
	WebDynamicUrlWebp    string `json:"web_dynamic_url_webp"`
	WebDynamicUrlApng    string `json:"web_dynamic_url_apng"`
	MobileDynamicUrlWebp string `json:"mobile_dynamic_url_webp"`
}
type sendGift struct {
	Action            string      `json:"action"`
	BatchComboId      string      `json:"batch_combo_id"`
	BatchComboSend    interface{} `json:"batch_combo_send"`
	BeatId            string      `json:"beatId"`
	BizSource         string      `json:"biz_source"`
	BlindGift         interface{} `json:"blind_gift"`
	BroadcastId       int         `json:"broadcast_id"`
	CoinType          string      `json:"coin_type"`
	ComboResourcesId  int         `json:"combo_resources_id"`
	ComboSend         interface{} `json:"combo_send"`
	ComboStayTime     int         `json:"combo_stay_time"`
	ComboTotalCoin    int         `json:"combo_total_coin"`
	CritProb          int         `json:"crit_prob"`
	Demarcation       int         `json:"demarcation"`
	DiscountPrice     int         `json:"discount_price"`
	Dmscore           int         `json:"dmscore"`
	Draw              int         `json:"draw"`
	Effect            int         `json:"effect"`
	EffectBlock       int         `json:"effect_block"`
	Face              string      `json:"face"`
	FaceEffectId      int         `json:"face_effect_id"`
	FaceEffectType    int         `json:"face_effect_type"`
	FloatScResourceId int         `json:"float_sc_resource_id"`
	GiftId            int         `json:"giftId"`
	GiftName          string      `json:"giftName"`
	GiftType          int         `json:"giftType"`
	Gold              int         `json:"gold"`
	GuardLevel        int         `json:"guard_level"`
	IsFirst           bool        `json:"is_first"`
	IsJoinReceiver    bool        `json:"is_join_receiver"`
	IsNaming          bool        `json:"is_naming"`
	IsSpecialBatch    int         `json:"is_special_batch"`
	Magnification     int         `json:"magnification"`
	MedalInfo         struct {
		AnchorRoomid     int    `json:"anchor_roomid"`
		AnchorUname      string `json:"anchor_uname"`
		GuardLevel       int    `json:"guard_level"`
		IconId           int    `json:"icon_id"`
		IsLighted        int    `json:"is_lighted"`
		MedalColor       int    `json:"medal_color"`
		MedalColorBorder int    `json:"medal_color_border"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		Special          string `json:"special"`
		TargetId         int    `json:"target_id"`
	} `json:"medal_info"`
	NameColor        string `json:"name_color"`
	Num              int    `json:"num"`
	OriginalGiftName string `json:"original_gift_name"`
	Price            int    `json:"price"`
	Rcost            int    `json:"rcost"`
	ReceiveUserInfo  struct {
		Uid   int    `json:"uid"`
		Uname string `json:"uname"`
	} `json:"receive_user_info"`
	Remain            int         `json:"remain"`
	Rnd               string      `json:"rnd"`
	SendMaster        interface{} `json:"send_master"`
	Silver            int         `json:"silver"`
	Super             int         `json:"super"`
	SuperBatchGiftNum int         `json:"super_batch_gift_num"`
	SuperGiftNum      int         `json:"super_gift_num"`
	SvgaBlock         int         `json:"svga_block"`
	Switch            bool        `json:"switch"`
	TagImage          string      `json:"tag_image"`
	Tid               string      `json:"tid"`
	Timestamp         int         `json:"timestamp"`
	TopList           interface{} `json:"top_list"`
	TotalCoin         int         `json:"total_coin"`
	Uid               int         `json:"uid"`
	Uname             string      `json:"uname"`
}
