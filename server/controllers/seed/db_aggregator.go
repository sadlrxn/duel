package seed

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm/clause"
)

// @Internal
// Set lock and retrieve user seed pair.
// 1. Lock seed pair record by user id.
// 2. Return seed pair.
// *. Return error in case of
//   - Invalid parameter
//   - Multiple retrieved unexpired seed pairs
//   - No retrieved unexpired seed pairs
func lockAndRetrieveUserActiveSeedPair(
	user db_aggregator.User,
	sessionId db_aggregator.UUID,
) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if user == 0 {
		return nil, utils.MakeError(
			"seed_db",
			"lockAndRetrieveUserSeedPair",
			"invalid parameter",
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"lockAndRetrieveUserSeedPair",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Lock and retrieve user seed pair.
	seedPairs := []models.SeedPair{}
	if result := session.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Preload(
		"ClientSeed",
	).Preload(
		"ServerSeed",
	).Preload(
		"NextServerSeed",
	).Preload(
		"User",
	).Where(
		"user_id = ?",
		user,
	).Where(
		"is_expired = ?",
		false,
	).Find(&seedPairs); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"lockAndRetrieveUserSeedPair",
			"failed to retrieve seed pair",
			result.Error,
		)
	}

	// 4. Check count of retrieved seed pairs.
	if len(seedPairs) == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"lockAndRetrieveUserSeedPair",
			"not found unexpired seed pair",
			ErrCodeNotFoundUnexpiredPair,
			fmt.Errorf("user: %d", user),
		)
	}

	if len(seedPairs) > 1 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"lockAndRetrieveUserSeedPair",
			"multiple retrieved unexpired seed pairs",
			ErrCodeMultipleUnexpiredPairs,
			fmt.Errorf("user: %d", user),
		)
	}

	// 5. Check whether the user is banned
	if seedPairs[0].User.Banned {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"lockAndRetrieveUserSeedPair",
			"banned user",
			ErrCodeBannedUser,
			fmt.Errorf("user: %d", user),
		)
	}

	return &seedPairs[0], nil
}

// @Internal
// Retrieve user info.
// To make seed management system to be safe,
// we utilize retrieving user info specifically.
// *. Retrun error in case of
//   - Invalid parameter
//   - Banned or not found user record
func getUserInfoChecked(user db_aggregator.User) (*models.User, error) {
	// 1. Validate parameter.
	if user == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"getUserInfoChecked",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"getUserInfoChecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve user info.
	userInfo := models.User{}
	if result := session.First(&userInfo, user); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"getUserInfoChecked",
			"failed to retrieve user info",
			result.Error,
		)
	}

	// 4. Check whether the user is banned.
	if userInfo.Banned {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"getUserInfoChecked",
			"banned user",
			ErrCodeBannedUser,
			fmt.Errorf("user: %d", user),
		)
	}

	return &userInfo, nil
}

// @Internal
// Takes the current seed pair as argument.
// Update the current seed pair as expired and create new seed pair.
// *. Returns error in case of
//   - Invalid parameter(Should check ID presence)
//   - Empty client seed string.
//   - Seed pair already expired.
//   - Seed pair still in use.
func expireSeedPairUnchecked(
	seedPair *models.SeedPair,
	clientSeed string,
	sessionId db_aggregator.UUID,
) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if seedPair == nil ||
		seedPair.ID == 0 ||
		seedPair.UserID == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"expireSeedPairUnchecked",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided seedPair argument is nil"),
		)
	}

	if clientSeed == "" {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"expireSeedPairUnchecked",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided clientSeed argument is empty"),
		)
	}

	// 2. Check whether the seed pair is in use.
	if seedPair.UsingCount != 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"expireSeedPairUnchecked",
			"pair is still in use",
			ErrCodePairsStillInUse,
			fmt.Errorf("seedPair: %v", seedPair),
		)
	}

	// 3. Check whether the seed pair is already expired.
	if seedPair.IsExpired {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"expireSeedPairUnchecked",
			"seed pair is already expired",
			ErrCodeExpiredPair,
			fmt.Errorf("seedPair: %v", seedPair),
		)
	}

	// 4. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"expireSeedPairUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 5. Update current seed to be expired one.
	seedPair.IsExpired = true
	if result := session.Save(&seedPair); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"expireSeedPairUnchecked",
			"failed to update expired seed pair",
			result.Error,
		)
	}

	// 6. Generate new server seed and save seed pairs.
	serverSeed, serverHash, err := utils.GenerateServerSeed(ServerSeedDefaultLength)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"expireSeedPairUnchecked",
			"failed to generate server seed",
			err,
		)
	}
	newSeedPair := models.SeedPair{
		UserID: seedPair.UserID,
		ClientSeed: models.ClientSeed{
			Seed: clientSeed,
		},
		ServerSeedID: seedPair.NextServerSeedID,
		NextServerSeed: models.ServerSeed{
			Seed: serverSeed,
			Hash: serverHash,
		},
	}
	if result := session.Create(&newSeedPair); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"expireSeedPairUnchecked",
			"failed to create new seed pair",
			result.Error,
		)
	}
	if result := session.Preload(
		"ServerSeed",
	).Preload(
		"ClientSeed",
	).Preload(
		"NextServerSeed",
	).Preload(
		"User",
	).First(&newSeedPair, newSeedPair.ID); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"expiredSeedPairUnchecked",
			"failed to get just saved seed pair",
			result.Error,
		)
	}

	return &newSeedPair, nil
}

// @External
// Get user's unhashed server seed.
// *. Return error in case of
//   - Invalid parameter
//   - Banned or not found user record
//   - Multiple retrieved server seed
//   - No retrieved server seed
//   - Not expired seed pair
func getUnhashedServerSeed(hashedServerSeed string) (string, error) {
	// 1. Validate parameter.
	if hashedServerSeed == "" {
		return "", utils.MakeErrorWithCode(
			"seed_db",
			"getUnhashedServerSeed",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided hashedServerSeed argument is empty"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return "", utils.MakeError(
			"seed_db",
			"getUnhashedServerSeed",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Retrieve server seed
	serverSeeds := []models.ServerSeed{}
	if result := session.Preload(
		"SeedPair.User",
	).Where(
		"hash = ?",
		hashedServerSeed,
	).Find(&serverSeeds); result.Error != nil {
		return "", utils.MakeError(
			"seed_db",
			"getUnhashedServerSeed",
			"failed to retrieve server seed",
			result.Error,
		)
	}

	// 4. Check count of retrieved seeds.
	if len(serverSeeds) == 0 {
		return "", utils.MakeErrorWithCode(
			"seed_db",
			"getUnhashedServerSeed",
			"not found server seed",
			ErrCodeNotFoundServerSeed,
			fmt.Errorf("hashedServerSeed: %s", hashedServerSeed),
		)
	}
	if len(serverSeeds) > 1 {
		return "", utils.MakeErrorWithCode(
			"seed_db",
			"getUnhashedServerSeed",
			"multiple server seeds",
			ErrCodeMultipleUnexpiredPairs,
			fmt.Errorf(
				"hashedServerSeed: %s, serverSeeds: %v",
				hashedServerSeed,
				serverSeeds,
			),
		)
	}

	// 5. Check whether the seed pair is not expired yet.
	serverSeed := &serverSeeds[0]
	if !serverSeed.SeedPair.IsExpired ||
		serverSeed.SeedPair.UsingCount != 0 {
		return "", utils.MakeErrorWithCode(
			"seed_db",
			"getUnhashedServerSeed",
			"not expired seed pair",
			ErrCodeUnexpiredPair,
			fmt.Errorf("serverSeed: %v", serverSeed),
		)
	}

	// 6. Check whether the creator is banned or not.
	if serverSeeds[0].SeedPair.User.Banned {
		return "", utils.MakeErrorWithCode(
			"seed_db",
			"getUnhashedServerSeed",
			"banned user",
			ErrCodeBannedUser,
			fmt.Errorf("serverSeed: %v", serverSeed),
		)
	}

	return serverSeeds[0].Seed, nil
}

// @Internal
// Increases nonce and using count.
func increaseNonceAndUsingCountUnchecked(
	seedPair *models.SeedPair,
	sessionId db_aggregator.UUID,
) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if seedPair == nil ||
		seedPair.ID == 0 {
		return nil, utils.MakeError(
			"seed_db",
			"increaseNonceAndUsingCountUnchecked",
			"invalid parameter",
			errors.New("provided seedPair is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"increaseNonceAndUsingCountUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Increase nonce and using count, and save.
	seedPair.Nonce += 1
	seedPair.UsingCount += 1
	if result := session.Save(&seedPair); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"increaseNonceAndUsingCountUnchecked",
			"failed to update seed pair's nonce and using count",
			result.Error,
		)
	}

	return seedPair, nil
}

// @Internal
// Decreases nonce and using count.
func decreaseNonceAndUsingCountUnchecked(
	seedPair *models.SeedPair,
	sessionId db_aggregator.UUID,
) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if seedPair == nil ||
		seedPair.ID == 0 {
		return nil, utils.MakeError(
			"seed_db",
			"increaseNonceAndUsingCountUnchecked",
			"invalid parameter",
			errors.New("provided seedPair is invalid"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"increaseNonceAndUsingCountUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Increase nonce and using count, and save.
	if seedPair.UsingCount == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"increaseNonceAndUsingCountUnchecked",
			"invalid seed pair's nonce and using count",
			ErrCodePairsNotInUse,
			fmt.Errorf(
				"nonce: %d, usingCount: %d",
				seedPair.Nonce,
				seedPair.UsingCount,
			),
		)
	}
	seedPair.UsingCount -= 1
	if result := session.Save(&seedPair); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"increaseNonceAndUsingCountUnchecked",
			"failed to update seed pair's nonce and using count",
			result.Error,
		)
	}

	return seedPair, nil
}

// @Internal
// Initiate new seed pair.
func initNewSeedPairUnchecked(
	user db_aggregator.User,
	sessionId db_aggregator.UUID,
) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if user == 0 {
		return nil, utils.MakeError(
			"seed_db",
			"initNewSeedPairUnchecked",
			"invalid parameter",
			errors.New("provided user param is nil"),
		)
	}

	// 2. Retrieve session.
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initNewSeedPairUnchecked",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Generate client, and server seeds.
	clientSeed, err := utils.GenerateClientSeed(ClientSeedDefaultLength)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initNewSeedPairUnchecked",
			"failed to generate client seed",
			err,
		)
	}
	serverSeed, serverHash, err := utils.GenerateServerSeed(ServerSeedDefaultLength)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initNewSeedPairUnchecked",
			"failed to generate server seed",
			err,
		)
	}
	nextServerSeed, nextServerHash, err := utils.GenerateServerSeed(ServerSeedDefaultLength)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initNewSeedPairUnchecked",
			"failed to generate next server seed",
			err,
		)
	}
	seedPair := models.SeedPair{
		UserID: uint(user),
		ClientSeed: models.ClientSeed{
			Seed: clientSeed,
		},
		ServerSeed: models.ServerSeed{
			Seed: serverSeed,
			Hash: serverHash,
		},
		NextServerSeed: models.ServerSeed{
			Seed: nextServerSeed,
			Hash: nextServerHash,
		},
	}

	// 4. Create seed pair.
	if result := session.Create(&seedPair); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initNewSeedPairUnchecked",
			"failed to create seed pair",
			result.Error,
		)
	}

	return &seedPair, nil
}

// @External
// Get user's expired seed pairs.
// *. Return error in case of
//   - Invalid parameter
//   - Banned or not found user record
func getExpiredUserSeedPairs(
	user db_aggregator.User,
	skip uint,
	limit uint,
) ([]models.SeedPair, error) {
	// 1. Validate parameters.
	if user == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"getExpiredUserSeedPairs",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user parameter is invalid"),
		)
	}
	if limit == 0 {
		return []models.SeedPair{}, nil
	}

	// 2. Get main session.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"getExpiredUserSeedPairs",
			"failed to get main session",
			err,
		)
	}

	// 3. Retrieve expired user seed pairs.
	seedPairs := []models.SeedPair{}
	if result := session.Preload(
		"ClientSeed",
	).Preload(
		"ServerSeed",
	).Preload(
		"NextServerSeed",
	).Where(
		"is_expired = ?", true,
	).Where(
		"user_id = ?", user,
	).Order(
		"id desc",
	).Offset(
		int(skip),
	).Limit(
		int(limit),
	).Find(&seedPairs); result.Error != nil {
		return nil, utils.MakeError(
			"seed_db",
			"getExpiredUserSeedPairs",
			"failed retrieve expired user seed pairs",
			result.Error,
		)
	}

	return seedPairs, nil
}
