package utils

import (
	"errors"

	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/models"
	"gorm.io/gorm"
)

func SetWinnerStatistics(userID uint, wagered int64, profit int64, game models.GameType) {
	db := db.GetDB()
	var statistics models.Statistics
	if result := db.Where("user_id = ?", userID).First(&statistics); result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		statistics.UserID = userID
		db.Create(&statistics)
	}
	statistics.TotalWin += profit
	statistics.TotalProfit += profit
	statistics.TotalWagered += wagered
	if profit > statistics.MaxProfit {
		statistics.MaxProfit = profit
	}
	statistics.WinStreaks++
	statistics.LoseStreaks = 0
	if statistics.WinStreaks > statistics.BestStreaks {
		statistics.BestStreaks = statistics.WinStreaks
	}
	switch game {
	case models.Jackpot:
		statistics.JackpotStats.TotalRounds++
		statistics.JackpotStats.WinnedRounds++
		statistics.JackpotStats.Wagered += wagered
		statistics.JackpotStats.Profit += profit
	case models.Coinflip:
		statistics.CoinflipStats.TotalRounds++
		statistics.CoinflipStats.WinnedRounds++
		statistics.CoinflipStats.Wagered += wagered
		statistics.CoinflipStats.Profit += profit
	case models.Dreamtower:
		statistics.DreamtowerStats.TotalRounds++
		statistics.DreamtowerStats.WinnedRounds++
		statistics.DreamtowerStats.Wagered += wagered
		statistics.DreamtowerStats.Profit += profit
	case models.Crash:
		statistics.CrashStats.TotalRounds++
		statistics.CrashStats.WinnedRounds++
		statistics.CrashStats.Wagered += wagered
		statistics.CrashStats.Profit += profit
	}
	db.Save(&statistics)
}

func SetLoserStatistics(userID uint, wagered int64, game models.GameType) {
	db := db.GetDB()
	var statistics models.Statistics
	if result := db.Where("user_id = ?", userID).First(&statistics); result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		statistics.UserID = userID
		db.Create(&statistics)
	}
	statistics.TotalLoss += wagered
	statistics.TotalProfit -= wagered
	statistics.TotalWagered += wagered
	statistics.WinStreaks = 0
	statistics.LoseStreaks++
	if statistics.LoseStreaks > statistics.WorstStreaks {
		statistics.WorstStreaks = statistics.LoseStreaks
	}
	switch game {
	case models.Jackpot:
		statistics.JackpotStats.TotalRounds++
		statistics.JackpotStats.LostRounds++
		statistics.JackpotStats.Wagered += wagered
		statistics.JackpotStats.Loss += wagered
	case models.Coinflip:
		statistics.CoinflipStats.TotalRounds++
		statistics.CoinflipStats.LostRounds++
		statistics.CoinflipStats.Wagered += wagered
		statistics.CoinflipStats.Loss += wagered
	case models.Dreamtower:
		statistics.DreamtowerStats.TotalRounds++
		statistics.DreamtowerStats.LostRounds++
		statistics.DreamtowerStats.Wagered += wagered
		statistics.DreamtowerStats.Loss += wagered
	case models.Crash:
		statistics.CrashStats.TotalRounds++
		statistics.CrashStats.LostRounds++
		statistics.CrashStats.Wagered += wagered
		statistics.CrashStats.Loss += wagered
	}
	db.Save(&statistics)
}
