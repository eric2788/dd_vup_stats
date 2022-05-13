package blive

import "encoding/json"

const (
	DanmuMsg         = "DANMU_MSG"
	SendGift         = "SEND_GIFT"
	GuardBuy         = "GUARD_BUY"
	SuperChatMessage = "SUPER_CHAT_MESSAGE"
	Live             = "LIVE"
	InteractWord     = "INTERACT_WORD"
	Preparing        = "PREPARING"
)

type (
	MapParser interface {
		Parse(m map[string]interface{}) error
	}

	LiveData struct {
		Command  string                 `json:"command"`
		LiveInfo LiveInfo               `json:"live_info"`
		Content  map[string]interface{} `json:"content"`
	}

	LiveInfo struct {
		UID             int64   `json:"uid"`
		Title           string  `json:"title"`
		Name            string  `json:"name"`
		Cover           *string `json:"cover"`
		RoomId          int64   `json:"room_id"`
		UserFace        string  `json:"user_face"`
		UserDescription string  `json:"user_description"`
	}

	ListeningInfo struct {
		LiveInfo
		OfficialRole int `json:"official_role"`
	}

	SuperChatMessageData struct {
		UID       int64  `json:"uid"`
		Price     int    `json:"price"`
		Message   string `json:"message"`
		StartTime int64  `json:"start_time"`

		BackgroundColorStart string `json:"background_color_start"`
		BackgroundImage      string `json:"background_image"`
		BackgroundColor      string `json:"background_color"`

		UserInfo struct {
			Face      string `json:"face"`
			NameColor string `json:"name_color"`
			UName     string `json:"uname"`
		} `json:"user_info"`
	}

	SendGiftData struct {
		Action         string `json:"action"`
		BatchComboID   string `json:"batch_combo_id"`
		BatchComboSend struct {
			Action        string      `json:"action"`
			BatchComboID  string      `json:"batch_combo_id"`
			BatchComboNum int         `json:"batch_combo_num"`
			BlindGift     interface{} `json:"blind_gift"`
			GiftID        int64       `json:"gift_id"`
			GiftName      string      `json:"gift_name"`
			GiftNum       int         `json:"gift_num"`
			SendMaster    interface{} `json:"send_master"`
			Uid           int         `json:"uid"`
			Uname         string      `json:"uname"`
		} `json:"batch_combo_send"`
		BeatID           string      `json:"beatId"`
		BizSource        string      `json:"biz_source"`
		BlindGift        interface{} `json:"blind_gift"`
		BroadcastID      int64       `json:"broadcast_id"`
		CoinType         string      `json:"coin_type"`
		ComboResourcesID int64       `json:"combo_resources_id"`
		ComboSend        struct {
			Action     string      `json:"action"`
			ComboID    string      `json:"combo_id"`
			ComboNum   int         `json:"combo_num"`
			GiftID     int64       `json:"gift_id"`
			GiftName   string      `json:"gift_name"`
			GiftNum    int         `json:"gift_num"`
			SendMaster interface{} `json:"send_master"`
			UID        int64       `json:"uid"`
			Uname      string      `json:"uname"`
		} `json:"combo_send"`
		ComboStayTime     int64   `json:"combo_stay_time"`
		ComboTotalCoin    int     `json:"combo_total_coin"`
		CritProb          int     `json:"crit_prob"`
		Demarcation       int     `json:"demarcation"`
		DiscountPrice     int     `json:"discount_price"`
		Dmscore           int     `json:"dmscore"`
		Draw              int     `json:"draw"`
		Effect            int     `json:"effect"`
		EffectBlock       int     `json:"effect_block"`
		Face              string  `json:"face"`
		FloatScResourceID int64   `json:"float_sc_resource_id"`
		GiftID            int64   `json:"giftId"`
		GiftName          string  `json:"giftName"`
		GiftType          int     `json:"giftType"`
		Gold              int     `json:"gold"`
		GuardLevel        int     `json:"guard_level"`
		IsFirst           bool    `json:"is_first"`
		IsSpecialBatch    int     `json:"is_special_batch"`
		Magnification     float64 `json:"magnification"`
		MedalInfo         struct {
			AnchorRoomid     int    `json:"anchor_roomid"`
			AnchorUname      string `json:"anchor_uname"`
			GuardLevel       int    `json:"guard_level"`
			IconID           int64  `json:"icon_id"`
			IsLighted        int    `json:"is_lighted"`
			MedalColor       int    `json:"medal_color"`
			MedalColorBorder int64  `json:"medal_color_border"`
			MedalColorEnd    int64  `json:"medal_color_end"`
			MedalColorStart  int64  `json:"medal_color_start"`
			MedalLevel       int    `json:"medal_level"`
			MedalName        string `json:"medal_name"`
			Special          string `json:"special"`
			TargetID         int    `json:"target_id"`
		} `json:"medal_info"`
		NameColor         string      `json:"name_color"`
		Num               int         `json:"num"`
		OriginalGiftName  string      `json:"original_gift_name"`
		Price             int         `json:"price"`
		Rcost             int         `json:"rcost"`
		Remain            int         `json:"remain"`
		Rnd               string      `json:"rnd"`
		SendMaster        interface{} `json:"send_master"`
		Silver            int         `json:"silver"`
		Super             int         `json:"super"`
		SuperBatchGiftNum int         `json:"super_batch_gift_num"`
		SuperGiftNum      int         `json:"super_gift_num"`
		SvgaBlock         int         `json:"svga_block"`
		TagImage          string      `json:"tag_image"`
		TID               string      `json:"tid"`
		Timestamp         int64       `json:"timestamp"`
		TopList           interface{} `json:"top_list"`
		TotalCoin         int         `json:"total_coin"`
		UID               int64       `json:"uid"`
		Uname             string      `json:"uname"`
	}

	GuardBuyData struct {
		GuardLevel int    `json:"guard_level"`
		Price      int    `json:"price"`
		UID        int64  `json:"uid"`
		Num        int    `json:"num"`
		GiftID     int64  `json:"gift_id"`
		GiftName   string `json:"gift_name"`
		StartTime  int64  `json:"start_time"`
		EndTime    int64  `json:"end_time"`
		Username   string `json:"username"`
	}
)

func (d *SuperChatMessageData) Parse(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}

func (d *SendGiftData) Parse(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}

func (d *GuardBuyData) Parse(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}
