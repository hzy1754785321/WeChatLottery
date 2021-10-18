package models

type WeChatModels struct {
}

//微信数据请求结构体
type WeChatUserRequest struct {
}

type WeChatUserInfo struct {
	Subscribe       int
	Openid          string
	Nickname        string
	Sex             int
	Language        string
	City            int
	Province        string
	Country         string
	Headimgurl      string
	Subscribe_time  int
	Unionid         string
	Remark          string
	Groupid         int
	Tagid_list      []int
	Subscribe_scene string
	QrScene         int
	QrSceneStr      string
}

//微信数据响应结构体
type WeChatUserResponse struct {
	WeChatUserInfos []WeChatUserInfo
}

type WeChatLotteryInfo struct {
	Openid         []string
	RewardItemName string
	RewradName     string //动名
}

type WeChatRewardInfo struct {
}

type WeChatQRCodeResponse struct {
	QrCodeUrl string
}

type WeChatLotteryResult struct {
	Winner       []WeChatLotteryInfo
	Loser        []string
	LotteryTitle string
}

type WeChatQRCodeRequest struct {
	SceneStr string
}


type WeChatLotteryResponse struct {
}

