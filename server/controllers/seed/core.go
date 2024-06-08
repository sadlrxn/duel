package seed

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

// @External
// Borrow user's seed pair.
// 1. Lock seed pair record by user id.
// 2. Increase using count.
// 3. Increase nonce.
// 4. Commit session and return seed pair.
// *. Return error in case of
//   - Invalid parameter
//   - Banned or not found user record
//   - Multiple retrieved unexpired seed pairs
//   - No retrieved unexpired seed pairs
//   - Too many borrow
func borrowUserSeedPair(user db_aggregator.User) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if user == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"borrowUserSeedPair",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"borrowUserSeedPair",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Lock and retrieve seed pair.
	seedPair, err := lockAndRetrieveUserActiveSeedPair(user, sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"borrowUserSeedPair",
			"failed to lock and retrieve user seed pair",
			err,
		)
	}

	// 4. Check too many borrows.
	if checkTooManyBorrows(seedPair) {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"borrowUserSeedPair",
			"too many borrows",
			ErrCodeTooManyBorrows,
			fmt.Errorf("seedPair: %v", seedPair),
		)
	}

	// 5. Increase nonce and using count.
	seedPair, err = increaseNonceAndUsingCountUnchecked(
		seedPair,
		sessionId,
	)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"borrowUserSeedPair",
			"failed to increase and using count",
			err,
		)
	}

	// 6. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"borrowUserSeedPair",
			"failed to commit session",
			err,
		)
	}

	return seedPair, nil
}

// @External
// Return user's seed pair.
// 1. Lock seed pair record by user id.
// 2. Decrease using count.
// 3. Commit session and return.
// *. Return error in case of
//   - Invalid parameter
//   - Banned or not found user record
//   - Multiple retrieved unexpired seed pairs
//   - No retrieved unexpired seed pairs
//   - Negative using count detected
//   - Not used seed pair
func returnUserSeedPair(user db_aggregator.User, seedPairId uint) error {
	// 1. Validate parameter.
	if user == 0 {
		return utils.MakeErrorWithCode(
			"seed_db",
			"returnUserSeedPair",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user argument is invalid"),
		)
	}
	if seedPairId == 0 {
		return utils.MakeErrorWithCode(
			"seed_db",
			"returnUserSeedPair",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided seedPairId argument is invalid"),
		)
	}

	// 2. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return utils.MakeError(
			"seed_db",
			"returnUserSeedPair",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Lock and retrieve seed pair.
	seedPair, err := lockAndRetrieveUserActiveSeedPair(user, sessionId)
	if err != nil {
		return utils.MakeError(
			"seed_db",
			"returnUserSeedPair",
			"failed to lock and retrieve user seed pair",
			err,
		)
	}

	// 4. Check whether owned seed pair id.
	if seedPair.ID != seedPairId {
		return utils.MakeErrorWithCode(
			"seed_db",
			"returnUserSeedPair",
			"not matching owned seed pair id",
			ErrCodeNotOwnedSeedPairID,
			fmt.Errorf(
				"seedPair: %v, seedPairId: %v",
				seedPair,
				seedPairId,
			),
		)
	}

	// 5. Increase nonce and using count.
	_, err = decreaseNonceAndUsingCountUnchecked(
		seedPair,
		sessionId,
	)
	if err != nil {
		return utils.MakeError(
			"seed_db",
			"returnUserSeedPair",
			"failed to decrease and using count",
			err,
		)
	}

	// 6. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return utils.MakeError(
			"seed_db",
			"returnUserSeedPair",
			"failed to commit session",
			err,
		)
	}

	return nil
}

// @External
// Rotate user's seed pair.
// 1. Lock seed pair record by user id.
// 2. Check whether the using count is zero.
// 3. Expire the current seed.
// 4. Create new client/server seed and seed pair record.
// 5. Commit session and return newly created seed pair.
// *. Return error in case of
//   - Invalid parameter
//   - Banned or not found user record
//   - Multiple retrieved unexpired seed pairs
//   - No retrieved unexpired seed pairs
//   - Positive using count detected(checking whether in-using)
//   - Negative using count detected(checking whether in-using)
//   - Too many rotates
//   - Empty client seed string
func rotateUserSeedPair(
	user db_aggregator.User,
	serverSeedHash string,
	clientSeed string) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if user == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"rotateUserSeedPair",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user argument is invalid"),
		)
	}
	if clientSeed == "" {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"rotateUserSeedPair",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided seedPairId argument is invalid"),
		)
	}

	// 2. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"rotateUserSeedPair",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Check too many rotates.
	tooMany, err := checkTooManyRotates(user, sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"rotateUserSeedPair",
			"failed to check too many rotates",
			err,
		)
	}
	if tooMany {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"rotateUserSeedPair",
			"too many rotates",
			ErrCodeTooManyRotates,
			fmt.Errorf("userId: %v", user),
		)
	}

	// 4. Lock and retrieve seed pair.
	seedPair, err := lockAndRetrieveUserActiveSeedPair(user, sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"rotateUserSeedPair",
			"failed to lock and retrieve user seed pair",
			err,
		)
	}

	// 5. Check whether the server seed hash is matching.
	if seedPair.ServerSeed.Hash != serverSeedHash {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"rotateUserSeedPair",
			"current server seed hash not matching",
			ErrCodeCurrentServerSeedHashNotMatching,
			fmt.Errorf(
				"expected: %s, actual: %s",
				serverSeedHash,
				seedPair.ServerSeed.Hash,
			),
		)
	}

	// 6. Expire seed pair.
	seedPair, err = expireSeedPairUnchecked(
		seedPair,
		clientSeed,
		sessionId,
	)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"rotateUserSeedPair",
			"failed to expire seed pair unchecked",
			err,
		)
	}

	// 7. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"rotateUserSeedPair",
			"failed to commit session",
			err,
		)
	}

	return seedPair, nil
}

// @External
// Init seed pairs for newly signed up user.
// 1. Lock seed pair record by user id.
// 2. Try retrieving seed pair for the user.
// 3. Check whether user doesn't have seed pair.
// 4. Initiate new seed pair for the user.
// 5. Commit session and return newly created seed pair.
// *. Return error in case of
//   - Invalid parameter
//   - Banned or not found user record
//   - Retrieved seed pair for the user
func initUserSeedPair(user db_aggregator.User) (*models.SeedPair, error) {
	// 1. Validate parameter.
	if user == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"initUserSeedPair",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user argument is invalid"),
		)
	}

	// 2. Start session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initUserSeedPair",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 3. Try retrieving user active seed pair.
	if seedPair, err := lockAndRetrieveUserActiveSeedPair(
		user, sessionId,
	); err == nil || !utils.IsErrorCode(err, ErrCodeNotFoundUnexpiredPair) {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"initUserSeedPair",
			"already existing unexpired pair",
			ErrCodeAlreadyExistingPair,
			fmt.Errorf("errorFromLocking: %v, seedPair: %v", err, seedPair),
		)
	}

	// 4. Init new seed pair.
	seedPair, err := initNewSeedPairUnchecked(user, sessionId)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initUserSeedPair",
			"failed to init new seed pair",
			err,
		)
	}

	// 5. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"initUserSeedPair",
			"failed to commit session",
			err,
		)
	}

	return seedPair, nil
}

// @External
// Retrieve user's currently active seed pair.
// *. Return error in case of
//   - Invalid parameter
//   - Banned or not found user record
//   - Multiple retrieved unexpired seed pairs
//   - No retrieved unexpired seed pair
func getActiveUserSeedPair(user db_aggregator.User) (*models.SeedPair, error) {
	// 1. Validate parameters.
	if user == 0 {
		return nil, utils.MakeErrorWithCode(
			"seed_db",
			"getActiveUserSeedPair",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided user parameter is invalid"),
		)
	}

	// 2. Retrieve active user seed pair.
	seedPair, err := lockAndRetrieveUserActiveSeedPair(
		user,
		db_aggregator.MainSessionId(),
	)
	if err != nil {
		return nil, utils.MakeError(
			"seed_db",
			"getActiveUserSeedPair",
			"failed to retrieve user active seed pair",
			err,
		)
	}

	return seedPair, nil
}

// @Internal
// Check too many borrow.
// Returns whether the seed pair is being used too many.
func checkTooManyBorrows(seedPair *models.SeedPair) bool {
	if seedPair == nil {
		return false
	}

	if seedPair.UsingCount >= SeedPairBorrowLimit {
		return true
	}

	return false
}

// @Internal
// Check too many rotates.
// Returns whether the seed pair is rotated too many.
func checkTooManyRotates(user db_aggregator.User, sessionId db_aggregator.UUID) (bool, error) {
	// 1. Validate parameter
	if user == 0 {
		return false, utils.MakeError(
			"seed_db",
			"checkTooManyRotates",
			"invalid parameter",
			errors.New("provided user is nil"),
		)
	}

	// 2. Retrieve session
	session, err := db_aggregator.GetSession(sessionId)
	if err != nil {
		return false, utils.MakeError(
			"seed_db",
			"checkTooManyRotates",
			"invalid parameter",
			err,
		)
	}

	// 3. Retrieve expired count in the last mins
	oneMinAgo := time.UnixMilli(time.Now().UnixMilli() - time.Minute.Milliseconds())
	rotatedCnt := int64(0)
	if result := session.Model(
		&models.SeedPair{},
	).Where(
		"updated_at > ?",
		oneMinAgo,
	).Where(
		"is_expired = ?",
		true,
	).Where(
		"user_id = ?",
		user,
	).Count(&rotatedCnt); result.Error != nil {
		return false, utils.MakeError(
			"seed_db",
			"checkTooManyRotates",
			"failed to retrieve expired count",
			result.Error,
		)
	}
	if rotatedCnt >= int64(SeedPairRotateLimitPerMinute) {
		return true, nil
	}

	return false, nil
}

// @Internal
// Check whether the seed pair can be expired or not.
// func canBeExpired(seedPair *models.SeedPair) bool {
// 	if seedPair == nil ||
// 		seedPair.ID == 0 ||
// 		seedPair.IsExpired ||
// 		seedPair.UsingCount != 0 ||
// 		seedPair.UserID == 0 {
// 		return false
// 	}

// 	return true
// }
