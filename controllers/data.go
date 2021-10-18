package controllers

import (
	m "WeChatLottery/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strconv"

	//"strconv"
	//"strings"
	"time"

	"github.com/astaxie/beego"
)

//DataController data
type DataController struct {
	beego.Controller
}

//HandleLogin 处理登录
func (c *DataController) HandleLogin() {
	username, password := c.GetString("login"), c.GetString("pwd")
	if m.CheckRedis(username) {
		var user m.UserInfo
		js := m.GetRedis(username)
		json.Unmarshal([]byte(js), &user)
		if user.Password == password {
			user.LastTime = time.Now().Format("2006-01-02 15:04:05")
			if user.LoginCount > 1 {
				msg := fmt.Sprintf("登录成功！<br /><br />欢迎回来,%s", user.Nickname)
				c.Data["json"] = map[string]interface{}{"status": true, "msg": msg}
			} else {
				c.Data["json"] = map[string]interface{}{"status": true, "msg": "登录成功！<br /><br />欢迎来到hzy的个人网站"}
			}
			user.LoginCount++
			jsDat, _ := json.Marshal(user)
			m.SetRedis(username, string(jsDat))
			var person m.PersonalInfo
			perjs := m.GetRedis("personal" + username)
			json.Unmarshal([]byte(perjs), &person)
			md := md5.New()
			io.WriteString(md, username)
			sessionID := fmt.Sprintf("%x", md.Sum(nil))
			var session m.Session
			session.UserInfoData = user
			session.PersonalData = person
			sessionDat, _ := json.Marshal(session)
			c.SetSession(sessionID, string(sessionDat))
			c.Ctx.SetCookie("sessionID", sessionID)
		} else {
			c.Data["json"] = map[string]interface{}{"status": false, "msg": "密码错误"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"status": false, "msg": "用户名不存在"}
	}
	c.ServeJSON()
}

//HandleRegister 处理注册
func (c *DataController) HandleRegister() {
	username, password := c.GetString("username"), c.GetString("pwd")
	nickname := c.GetString("nickname")
	if m.CheckRedis(username) {
		c.Data["json"] = map[string]interface{}{"status": false, "msg": "用户名已存在"}
	} else {
		var user m.UserInfo
		var person m.PersonalInfo
		user.Username = username
		user.Nickname = nickname
		user.Password = password
		user.Icon = "/static/img/pic.jpg"
		user.LoginCount = 1
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		user.LastTime = nowTime
		js, _ := json.Marshal(user)
		m.SetRedis(username, string(js))
		person.Username = username
		person.Nickname = nickname
		person.Created = nowTime
		coinAdd := rand.Intn(10)
		person.Coin += coinAdd
		perjs, _ := json.Marshal(person)
		m.SetRedis("personal"+username, string(perjs))
		c.Data["json"] = map[string]interface{}{"status": true, "msg": "注册成功"}
	}
	c.ServeJSON()
}

//GetSessionUserInfo  获得session
func (c *DataController) GetSessionUserInfo() {
	sessionID := c.GetString("sessionID")
	sessionTemp := c.GetSession(sessionID)
	sessionDat, _ := sessionTemp.(string)
	var user m.UserInfo
	var session m.Session
	json.Unmarshal([]byte(sessionDat), &session)
	user = session.UserInfoData
	userjs, _ := json.Marshal(user)
	userData := make(map[string]interface{})
	json.Unmarshal(userjs, &userData)
	c.Data["json"] = userData
	c.ServeJSON()
}

func (c *DataController) GetSessionPersonal() {
	sessionID := c.GetString("sessionID")
	sessionTemp := c.GetSession(sessionID)
	sessionDat, _ := sessionTemp.(string)
	var session m.Session
	var person m.PersonalInfo
	json.Unmarshal([]byte(sessionDat), &session)
	person = session.PersonalData
	perjs, _ := json.Marshal(person)
	per := make(map[string]interface{})
	json.Unmarshal(perjs, &per)
	c.Data["json"] = per
	c.ServeJSON()
}

func (c *DataController) SavePersonal() {
	var person m.PersonalInfo
	sessionID := c.GetString("sessionID")
	sessionTemp := c.GetSession(sessionID)
	sessionDat, _ := sessionTemp.(string)
	var session m.Session
	json.Unmarshal([]byte(sessionDat), &session)
	person = session.PersonalData
	person.Nickname, person.Birth, person.Sex = c.GetString("NickName"), c.GetString("Birth"), c.GetString("Sex")
	if c.GetString("Province") != "" {
		person.City, person.Province, person.Area = c.GetString("City"), c.GetString("Province"), c.GetString("Area")
	}
	perjs, _ := json.Marshal(person)
	m.SetRedis("personal"+person.Username, string(perjs))
	session.UserInfoData = session.UserInfoData
	session.PersonalData = person
	sessionDats, _ := json.Marshal(session)
	c.SetSession(sessionID, string(sessionDats))
	c.Data["json"] = map[string]interface{}{"status": false, "msg": "修改成功！"}
	c.ServeJSON()
}

func (c *DataController) SaveIcon() {
	sessionID := c.GetString("sessionID")
	icon := c.GetString("Icon")
	sessionTemp := c.GetSession(sessionID)
	sessionDat, _ := sessionTemp.(string)
	var session m.Session
	json.Unmarshal([]byte(sessionDat), &session)
	session.UserInfoData.Icon = icon
	session.PersonalData.Icon = icon
	perjs, _ := json.Marshal(session.PersonalData)
	m.SetRedis("personal"+session.PersonalData.Username, string(perjs))
	loginjs, _ := json.Marshal(session.UserInfoData)
	m.SetRedis(session.UserInfoData.Username, string(loginjs))
	sessionDats, _ := json.Marshal(session)
	c.SetSession(sessionID, string(sessionDats))
	c.Data["json"] = map[string]interface{}{"status": true, "msg": "修改成功！", "path": icon}
	c.ServeJSON()
}

func (c *DataController) ChangeSecurity() {
	sessionID := c.GetString("sessionID")
	sessionTemp := c.GetSession(sessionID)
	sessionDat, _ := sessionTemp.(string)
	var session m.Session
	json.Unmarshal([]byte(sessionDat), &session)
	securityType := c.GetString("securityType")
	switch securityType {
	case "0":
		session.UserInfoData.Password = c.GetString("Password")
	case "1":
		session.UserInfoData.Email = c.GetString("Email")
	case "2":
		session.UserInfoData.Iphone = c.GetString("iphone")
	}
	loginjs, _ := json.Marshal(session.UserInfoData)
	m.SetRedis(session.UserInfoData.Username, string(loginjs))
	sessionDats, _ := json.Marshal(session)
	c.SetSession(sessionID, string(sessionDats))
	c.Data["json"] = map[string]interface{}{"status": true, "msg": "修改成功！"}
	c.ServeJSON()
}

//getClientIp 获取用户IP地址
func (p *DataController) getClientIp() string {
	var ip = p.Ctx.Input.IP()
	fmt.Printf(ip)
	// s := strings.Split(p.Ctx.Request.RemoteAddr, ":")
	return ip
}

func (c *DataController) CheckPassword() {
	sessionID := c.GetString("sessionID")
	sessionTemp := c.GetSession(sessionID)
	sessionDat, _ := sessionTemp.(string)
	var session m.Session
	json.Unmarshal([]byte(sessionDat), &session)
	password := c.GetString("pwd")
	if password != session.UserInfoData.Password {
		c.Data["json"] = map[string]interface{}{"status": false, "msg": "密码错误"}
	} else {
		c.Data["json"] = map[string]interface{}{"status": true, "msg": "密码正确"}
	}
	c.ServeJSON()
}

type LotteryStart struct {
	LotteryId string
}

func (lotteryStart LotteryStart) Run() {
	var lotteryInfo m.LotteryInfo
	var lottertResult m.WeChatLotteryResult
	var data = m.GetRedis("lott" + lotteryStart.LotteryId)
	json.Unmarshal([]byte(data), &lotteryInfo)
	var joinNumber = lotteryInfo.JoinNumber
	var winner []string
	for _, rewardInfo := range lotteryInfo.RewardInfo {
		winner, joinNumber = MicsSlice(joinNumber, rewardInfo.RewardCount)
		var weChatLotteryInfo m.WeChatLotteryInfo
		weChatLotteryInfo.Openid = winner
		weChatLotteryInfo.RewradName = rewardInfo.RewardName
		weChatLotteryInfo.RewardItemName = rewardInfo.RewordItems
		lottertResult.Winner = append(lottertResult.Winner, weChatLotteryInfo)
	}
	lottertResult.Loser = joinNumber
	lottertResult.LotteryTitle = lotteryInfo.LotteryTitle
	NoticeResult(lottertResult)
}

func MicsSlice(origin []string, count int) ([]string, []string) {
	tmpOrigin := make([]string, len(origin))
	copy(tmpOrigin, origin)
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(tmpOrigin), func(i int, j int) {
		tmpOrigin[i], tmpOrigin[j] = tmpOrigin[j], tmpOrigin[i]
	})

	result := make([]string, 0, count)
	for index, value := range tmpOrigin {
		if index == count {
			break
		}
		result = append(result, value)
	}
	tmpOrigin = tmpOrigin[count:]
	return result, tmpOrigin
}

func (c *DataController) AddNewLottery() {

	lotteryTitleStr := c.GetString("lotteryTitle")
	lotteryEndStr := c.GetString("lotteryEnd")
	rewardInfoStr := c.GetString("rewardInfo")

	// 	lottType, lottEndTime := c.GetString("lottTtpe"), c.GetString("lottEndTime")

	// 	rewardCount, _ := strconv.Atoi(rewardCountStr)

	// var lottery m.LotteryInfo

	// lottery.LotteryTitle = lotteryTitleStr
	// lottery.LotteryCreate = time.Now().Format("2006-01-02 15:04:05")
	// lottery.LotteryEnd = lotteryEndStr
	// lottery.JoinNumber = nil
	// lottery.RewardInfo = rewardInfoStr

	var rewardInfo []m.RewardInfo
	json.Unmarshal([]byte(rewardInfoStr), &rewardInfo)

	// 	var rewardInfos = make([]m.RewardInfo, rewardCount)

	// 	rewardNameSplit := strings.Split(rewardName, ";")
	// 	rewardPeopleNumSplit := strings.Split(rewardPeopleNum, ";")
	// 	rewardItemNameSplit := strings.Split(rewardItemName, "&")
	// 	for i := 0; i < rewardCount; i++ {
	// 		rewardInfos[i].RewardName = rewardNameSplit[i]
	// 		rewardInfos[i].RewardPeopleNum = rewardPeopleNumSplit[i]
	// 		ItemNameSplit := strings.Split(rewardItemNameSplit[i], ";")
	// 		rewardItemInfos := make([]m.RewordItemInfo, len(ItemNameSplit))
	// 		for j := 0; j < len(ItemNameSplit); j++ {
	// 			rewardItemInfos[j].RewardItemName = ItemNameSplit[i]
	// 			rewardItemInfos[j].RewordState = "false"
	// 		}
	// 		rewardInfos[i].RewordItems = rewardItemInfos
	// 	}

	var lotteryRedisKey = "LotteryIDList"
	var lotteryRedisCountKey = "LotteryListCount"
	var lotteryRedisCount = 1
	var lotteryValue = 0001
	if m.CheckRedis(lotteryRedisKey) {
		lotteryRedisCount, _ = strconv.Atoi(m.GetRedis(lotteryRedisCountKey))
		lotteryRedisCount++
		lotteryValue = lotteryRedisCount
		m.AddRedisSet(lotteryRedisKey, fmt.Sprintf("%04d", lotteryRedisCount))
		m.SetRedis(lotteryRedisCountKey, strconv.Itoa(lotteryRedisCount))
	} else {
		m.AddRedisSet(lotteryRedisKey, fmt.Sprintf("%04d", lotteryRedisCount))
		m.SetRedis(lotteryRedisCountKey, strconv.Itoa(lotteryRedisCount))
	}
	// rewardInfos := []m.RewardInfo{}
	// rewardInfos = rewardInfo

	var lotteryInfo m.LotteryInfo
	lotteryInfo.LotteryID = fmt.Sprintf("%04d", lotteryValue)
	lotteryInfo.LotteryTitle = lotteryTitleStr
	lotteryInfo.LotteryCreate = time.Now().Format("2006-01-02 15:04")
	lotteryInfo.LotteryEnd = lotteryEndStr
	lotteryInfo.JoinPeopleNum = 0
	lotteryInfo.JoinNumber = nil
	lotteryInfo.RewardInfo = rewardInfo
	lottInfojs, _ := json.Marshal(lotteryInfo)
	m.SetRedis("lott"+lotteryInfo.LotteryID, string(lottInfojs))

	lottertTime, _ := time.Parse("2006-01-02 15:04:05", lotteryEndStr)
	spec := fmt.Sprintf("%d %d %d %d %d *", lottertTime.Second(), lottertTime.Minute(), lottertTime.Hour(), lottertTime.Day(), lottertTime.Month())
	lotteryStart := LotteryStart{lotteryInfo.LotteryID}
	m.CronControl.AddJob(spec, lotteryStart)

	c.Data["json"] = map[string]interface{}{"status": true, "lotteryID": lotteryValue}
	c.ServeJSON()
}

func (c *DataController) GetLotteryInfos() {
	var lotteryRedisKey = "LotteryIDList"
	var lotteryInfojs []string
	var lotteryInfo m.LotteryInfo
	var lottertInfoResult m.LotteryInfoReturn
	lottCount := 0

	var lotteryInfos []m.LotteryInfo
	if m.CheckRedis(lotteryRedisKey) {
		lotteryInfojs = m.GetRedisSet(lotteryRedisKey)
		for _, v := range lotteryInfojs {
			var data = m.GetRedis("lott" + v)
			json.Unmarshal([]byte(data), &lotteryInfo)
			lottCount += 1
			lottertInfoResult.LotteryInfos = append(lottertInfoResult.LotteryInfos, lotteryInfo)
			lotteryInfos = append(lotteryInfos, lotteryInfo)
		}
	}

	// lottertInfoResult.Status = 200
	// lottertInfoResult.Message = ""
	// lottertInfoResult.Total = 1

	// mapResult["status"] = 200
	// mapResult["message"] = ""
	// mapResult["total"] = 1

	// var ret = make(map[string]interface{})
	// ret["status"] = true
	// ret["lotteryIDs"] = lottIDs
	// ret["lotteryTypes"] = lottTitle
	// ret["lotteryCreates"] = lottCreates
	// ret["lotteryStarts"] = lottStarts
	// ret["lotteryEnds"] = lottEnds
	// ret["joinNumbers"] = joinNumbers
	// ret["rewardNames"] = rewardNames
	// ret["rewardPeopleNums"] = rewardPeopleNums
	// ret["rewardItemName"] = rewardItemName
	// ret["rewordState"] = rewordState
	// ret["rewardCount"] = rewardCount
	// ret["rewardItemCount"] = rewardItemCount
	// ret["lottertCount"] = strconv.Itoa(lottCount)
	//	lottertInfoJson, _ := json.Marshal(lottertInfoResult)
	lottertInfoJson, _ := json.Marshal(lotteryInfos)
	//var items = map[string]interface{}{"item": string(lottertInfoJson)}
	//c.Data["json"] = lotteryInfos
	c.Data["json"] = map[string]interface{}{"code": 200, "msg": "", "count": lottCount, "rows": lotteryInfos}
	fmt.Println(string(lottertInfoJson))

	c.ServeJSON()
}

// Convert map json string
func MapToJson(m map[string]interface{}) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", nil
	}

	return string(jsonByte), nil
}

//切片转换为Map
func slice_To_Map(s_key, s_value []string) map[string]string {
	mapObj := map[string]string{}
	for s_key_index := range s_key {
		mapObj[s_key[s_key_index]] = s_value[s_key_index]
	}
	return mapObj
}
