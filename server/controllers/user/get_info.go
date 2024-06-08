package user

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

// @Internal
// Get user info
func (c *Controller) getInfo(request userInfoRequest, viewerID uint) (*userInfoResponse, error) {
	var targetUser *models.User
	if request.UserID != nil {
		targetUser = GetUserInfoByID(*request.UserID)
		if targetUser == nil {
			return nil, utils.MakeErrorWithCode(
				"user_get_info",
				"getInfo",
				"failed to get target user info by ID",
				ErrCodeFailedToGetTargetUserInfo,
				nil,
			)
		}
	} else if request.UserName != nil {
		targetUser = getUserInfoByName(*request.UserName)
		if targetUser == nil {
			return nil, utils.MakeErrorWithCode(
				"user_get_info",
				"getInfo",
				"failed to get target user info by name",
				ErrCodeFailedToGetTargetUserInfo,
				nil,
			)
		}
	} else {
		return nil, utils.MakeErrorWithCode(
			"user_get_info",
			"getInfo",
			"invalid parameter",
			ErrCodeInvalidParameter,
			nil,
		)
	}

	if checkHiddenUser(targetUser.Name) {
		return nil, utils.MakeErrorWithCode(
			"user_get_info",
			"getInfo",
			"target user is duel hidden user",
			ErrCodeIsDuelHiddenUser,
			nil,
		)
	}

	var viewer *models.User
	if viewerID != 0 {
		viewer = GetUserInfoByID(viewerID)
		if viewer == nil {
			return nil, utils.MakeErrorWithCode(
				"user_get_info",
				"getInfo",
				"failed to get viewer user info",
				ErrCodeFailedToGetViewerUserInfo,
				nil,
			)
		}
	}

	targetUserInfo := utils.GetUserDataWithPermissions(
		*targetUser,
		viewer,
		0,
		c.Chat.IsMuted(targetUser.ID),
	)

	var statistics *userStatistics
	if !shouldHideStatistics(*targetUser, viewer) {
		statistics = &userStatistics{
			TotalRounds: targetUser.Statistics.JackpotStats.TotalRounds +
				targetUser.Statistics.CoinflipStats.TotalRounds +
				targetUser.Statistics.DreamtowerStats.TotalRounds +
				targetUser.Statistics.CrashStats.TotalRounds,
			WinnedRounds: targetUser.Statistics.JackpotStats.WinnedRounds +
				targetUser.Statistics.CoinflipStats.WinnedRounds +
				targetUser.Statistics.DreamtowerStats.WinnedRounds +
				targetUser.Statistics.CrashStats.WinnedRounds,
			LostRounds: targetUser.Statistics.JackpotStats.LostRounds +
				targetUser.Statistics.CoinflipStats.LostRounds +
				targetUser.Statistics.DreamtowerStats.LostRounds +
				targetUser.Statistics.CrashStats.LostRounds,
			BestStreaks:    targetUser.Statistics.BestStreaks,
			WorstStreaks:   targetUser.Statistics.WorstStreaks,
			TotalWagered:   targetUser.Statistics.TotalWagered,
			PrivateProfile: targetUser.PrivateProfile,
		}
	}

	return &userInfoResponse{
		Info:       targetUserInfo,
		Statistics: statistics,
	}, nil
}

// @Internal
// Checks whether is Duel hidden users
func checkHiddenUser(name string) bool {
	for _, hidden_user := range config.HIDDEN_USERS {
		if hidden_user == name {
			return true
		}
	}
	return false
}

// @Internal
// Check should hide statistics.
func shouldHideStatistics(targetUser models.User, viewer *models.User) bool {
	return (viewer == nil ||
		(targetUser.ID != viewer.ID &&
			viewer.Role != models.AdminRole &&
			viewer.Role != models.ModeratorRole)) &&
		targetUser.PrivateProfile
}
