package controllers

import (
	m "WeChatLottery/models"
	"fmt"
	"net/rpc"

	"github.com/astaxie/beego"
)

//DataController data
type WeChatController struct {
	beego.Controller
}

func (c *WeChatController) GetWeChatUserInfo() []m.WeChatUserInfo {
	conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8095")
	if err != nil {
		fmt.Printf("Start RPC Error: %+v\n", err)
	}
	defer conn.Close()
	var req m.WeChatUserRequest
	var resp m.WeChatUserResponse
	err = conn.Call("WeChatModels.RequestUserInfo", req, &resp)
	if err != nil {
		fmt.Printf("RPC Call Error: %+v\n", err)
	}
	var userInfos []m.WeChatUserInfo
	for index, userInfo := range resp.WeChatUserInfos {
		userInfos[index] = userInfo
	}
	return userInfos
}

func (c *WeChatController) GetQRCode() {

	lottertId := c.GetString("lottertId")
	if m.CheckRedis("lott" + lottertId) {
		conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8095")
		if err != nil {
			fmt.Printf("Start RPC Error: %+v\n", err)
		}
		defer conn.Close()
		var req m.WeChatQRCodeRequest
		var resp m.WeChatQRCodeResponse
		req.SceneStr = lottertId
		err = conn.Call("WeChatModels.GetQRCodeUrl", req, &resp)
		if err != nil {
			fmt.Printf("RPC Call Error: %+v\n", err)
		}

		c.Data["json"] = map[string]interface{}{"qrCodeUrl": resp.QrCodeUrl, "status": true}
		c.ServeJSON()
	} else {
		c.Data["json"] = map[string]interface{}{"qrCodeUrl": nil, "status": false}
		c.ServeJSON()
	}

}

func NoticeResult(req m.WeChatLotteryResult) {
	conn, err := rpc.DialHTTP("tcp", "127.0.0.1:8095")
	if err != nil {
		fmt.Printf("Start RPC Error: %+v\n", err)
	}
	defer conn.Close()

	var resp m.WeChatLotteryResponse
	err = conn.Call("WeChatModels.RequestNotice", req, resp)
	if err != nil {
		fmt.Printf("RPC Call Error: %+v\n", err)
	}

	//	var resp m.WeChatQRCodeResponse
}
