package mutex

// var userSeedInUsage syncmap.Map

// func LockUserSeed(userID uint) {
// 	userSeedInUsage.Store(userID, true)
// }

// func ReleaseUserSeed(userID uint) {
// 	if _, prs := userSeedInUsage.Load(userID); prs {
// 		userSeedInUsage.Delete(userID)
// 	}
// }

// func CheckUserSeedInUsage(userID uint) (prs bool) {
// 	_, prs = userSeedInUsage.Load(userID)
// 	return
// }
