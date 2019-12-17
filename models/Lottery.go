package models


type LotteryInfo struct {
	LotteryID      string        `json:"lotteryID"`
	LotteryType    string 	     `json:"lotteryType"`
	LotteryCreate  string        `json:"lotteryCreate"`
	LotteryStart   string 	     `json:"lotteryStart"`
	LotteryEnd     string        `json:"lotteryEnd"`
	JoinNumber     int 		     `json:"joinNumber"`
	RewardInfo     []RewardInfo  `json:"rewardInfo"`
 }

 
type RewardInfo struct {
	RewardName  	 string              `json:"rewardName"`
	RewardPeopleNum  string 		     `json:"rewardPeopleNum"`
	RewordItems      []RewordItemInfo    `json:"rewordItems"`
 }

 type RewordItemInfo struct {
	RewardItemName   string         `json:"rewardItemName"`
	RewordState      string         `json:"rewordState"`
 }
