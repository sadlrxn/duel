package utils

import "github.com/Duelana-Team/duelana-v1/config"

func RevShareFromFee(fee int64) int64 {
	return fee * config.DUEL_BOT_TOTAL_SHARE / 100
}
