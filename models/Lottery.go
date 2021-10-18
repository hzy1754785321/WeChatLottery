package models

import "github.com/robfig/cron"

var CronControl *cron.Cron

type LotteryInfo struct {
	LotteryID     string       `json:"lotteryID"`
	LotteryTitle  string       `json:"lotteryTitle"`
	LotteryCreate string       `json:"lotteryCreate"`
	LotteryEnd    string       `json:"lotteryEnd"`
	JoinNumber    []string     `json:"joinNumber"`
	JoinPeopleNum int          `json:"JoinPeopleNum"`
	RewardInfo    []RewardInfo `json:"rewardInfo"`
}

type RewardInfo struct {
	RewardName   string   `json:"rewardName"`
	RewardCount  int      `json:"rewardCount"`
	RewardPeople []string `json:"rewardPeople"`
	RewordItems  string   `json:"rewordItems"`
}

type LotteryInfoReturn struct {
	//	Status       int           `json:"status"`
	//	Message      string        `json:"message"`
	//	Total        int           `json:"total"`
	LotteryInfos []LotteryInfo //`json:"-"`
}
