package duel_bots

import "github.com/Duelana-Team/duelana-v1/models"

type DuelBotMeta struct {
	Status        models.DuelBotStatus `json:"status"`
	TotalEarned   int64                `json:"totalEarned"`
	StakingReward int64                `json:"stakingReward"`
	Name          string               `json:"name"`
	MintAddress   string               `json:"mintAddress"`
	Image         string               `json:"image"`
}
