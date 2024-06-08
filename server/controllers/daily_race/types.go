package daily_race

import "time"

type WinnerInDailyRacePrizingResult struct {
	UserID uint  `json:"userId"`
	Rank   uint  `json:"rank"`
	Prize  int64 `json:"prize"`
}

type DailyRacePrizingResult struct {
	Date    time.Time                        `json:"date"`
	Index   int                              `json:"index"`
	Winners []WinnerInDailyRacePrizingResult `json:"winners"`
}

type UserInDailyRaceStatus struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Rank    int    `json:"rank"`
	Wagered int64  `json:"wagered"`
}

type DailyRaceRunningStatus string

const (
	DailyRaceStatusRunning DailyRaceRunningStatus = "running"
	DailyRaceStatusPending DailyRaceRunningStatus = "pending"
)

type DailyRaceStatus struct {
	Me        UserInDailyRaceStatus   `json:"me"`
	Players   []UserInDailyRaceStatus `json:"players"`
	Prizes    []int64                 `json:"prizes"`
	Remaining uint                    `json:"remaining"`
	Status    DailyRaceRunningStatus  `json:"status"`
}

type DailyRaceRewardsStatus struct {
	ID    uint      `json:"id"`
	Date  time.Time `json:"date"`
	Prize int64     `json:"prize"`
	Rank  uint      `json:"rank"`
}

type DetailInDailyRaceRewardResult struct {
	UserID   uint      `json:"userId"`
	Name     string    `json:"name"`
	Prize    int64     `json:"prize"`
	Rank     uint      `json:"rank"`
	Date     time.Time `json:"date"`
	RewardID uint      `json:"rewardId"`
}

type DailyRaceRewardResult struct {
	Count        uint                            `json:"count"`
	TotalPrize   int64                           `json:"totalPrize"`
	PrizeDetails []DetailInDailyRaceRewardResult `json:"prizeDetails"`
	RewardIDs    []uint                          `json:"rewardIds"`
}
