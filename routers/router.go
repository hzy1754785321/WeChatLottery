package routers

import (
	"WeChatLottery/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/home", &controllers.MainController{}, "*:Home")
	beego.Router("/handleLogin", &controllers.DataController{}, "post:HandleLogin")
	beego.Router("/HandleRegister", &controllers.DataController{}, "post:HandleRegister")
	beego.Router("/GetSessionUserInfo", &controllers.DataController{}, "post:GetSessionUserInfo")
	beego.Router("/GetSessionPersonal", &controllers.DataController{}, "post:GetSessionPersonal")
	beego.Router("/SavePersonal", &controllers.DataController{}, "post:SavePersonal")
	beego.Router("/SaveIcon", &controllers.DataController{}, "post:SaveIcon")
	beego.Router("/CheckPassword", &controllers.DataController{}, "post:CheckPassword")
	beego.Router("/ChangeSecurity", &controllers.DataController{}, "post:ChangeSecurity")
	beego.Router("/personal", &controllers.MainController{}, "*:Personal")
	beego.Router("/personal/userInfo", &controllers.MainController{}, "*:UserInfo")
	beego.Router("/personal/security", &controllers.MainController{}, "*:Security")
	beego.Router("/personal/icon", &controllers.MainController{}, "*:Icon")
	beego.Router("/UploadFiles", &controllers.SendController{}, "post:UploadFiles")
	beego.Router("/SendVerifyMail", &controllers.SendController{}, "post:SendVerifyMail")
	beego.Router("/VerifyMailCode", &controllers.SendController{}, "post:VerifyMailCode")
	beego.Router("/lotteryManage", &controllers.MainController{}, "*:LotteryManage")
	beego.Router("/AddNewLottery", &controllers.DataController{}, "post:AddNewLottery")
	beego.Router("/GetLotteryInfos", &controllers.DataController{}, "post:GetLotteryInfos")
	beego.Router("/GetWeChatUserInfo", &controllers.WeChatController{}, "post:GetWeChatUserInfo")
	beego.Router("/GetQRCode", &controllers.WeChatController{}, "post:GetQRCode")
}
