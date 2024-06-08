package utils

import (
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
)

func GetUserDataWithPermissions(userInfo models.User, viewer *models.User, count uint, muted ...bool) types.User {
	var user types.User
	user.ID = userInfo.ID
	user.Avatar = userInfo.Avatar
	user.Role = userInfo.Role
	user.Banned = userInfo.Banned
	if len(muted) > 0 {
		user.Muted = muted[0]
	}
	if (viewer == nil || (viewer.Role != models.AdminRole && viewer.Role != models.ModeratorRole && viewer.ID != userInfo.ID)) &&
		userInfo.PrivateProfile {
		user.Name = "HIDDEN"
	} else {
		user.Name = userInfo.Name
		user.WalletAddress = userInfo.WalletAddress
	}
	if count != 0 {
		user.Count = count
	}
	return user
}
