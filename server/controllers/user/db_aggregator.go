package user

import (
	"errors"
	"strings"

	"github.com/Duelana-Team/duelana-v1/controllers/seed"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @Internal
// Get user info by ID
func GetUserInfoByID(userID uint) *models.User {
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	userInfo := models.User{}
	if result := session.Preload(
		"Statistics",
	).Select(
		"id",
		"name",
		"wallet_address",
		"role",
		"avatar",
		"banned",
		"private_profile",
	).First(
		&userInfo,
		userID,
	); result.Error != nil {
		return nil
	}

	return &userInfo
}

// @Internal
// Get total statistic values
func getServerStatistics() (*ServerStatisticsResult, error) {
	// 1. Get main session from DB.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, utils.MakeError(
			"user_db_aggregator",
			"getServerStatistics",
			"failed to get main session",
			err,
		)
	}

	var totalBets, totalWagered, totalProfit int64
	var coinflipBetCount, jackpotBetCount, dreamtowerBetCount, crashBetCount int64

	// 2. Get total wagered amount.
	if err := session.Model(
		&models.Statistics{},
	).Select(
		"SUM(total_wagered)",
	).Row().Scan(
		&totalWagered,
	); err != nil {
		return nil, utils.MakeError(
			"user_db_aggregator",
			"getServerStatistics",
			"failed to get total wagered amount",
			err,
		)
	}

	// 3. Get total winned amount.
	if err := session.Model(
		&models.Statistics{},
	).Select(
		"SUM(total_win)",
	).Row().Scan(
		&totalProfit,
	); err != nil {
		return nil, utils.MakeError(
			"user_db_aggregator",
			"getServerStatistics",
			"failed to get total win amount",
			err,
		)
	}

	// 4. Get coinflip bet count.
	if result := session.Model(
		&models.CoinflipRound{},
	).Where(
		"winner_id is not null",
	).Count(
		&coinflipBetCount,
	); result.Error != nil {
		return nil, utils.MakeError(
			"user_db_aggregator",
			"getServerStatistics",
			"failed to get coinflip rounds",
			result.Error,
		)
	}
	coinflipBetCount *= 2

	// 5. Get jackpot bet count.
	if result := session.Model(
		&models.JackpotBet{},
	).Count(
		&jackpotBetCount,
	); result.Error != nil {
		return nil, utils.MakeError(
			"user_db_aggregator",
			"getServerStatistics",
			"failed to get jackpot bet count",
			result.Error,
		)
	}

	// 6. Get dreamtower bet count.
	if result := session.Model(
		&models.DreamTowerRound{},
	).Count(
		&dreamtowerBetCount,
	); result.Error != nil {
		return nil, utils.MakeError(
			"user_db_aggregator",
			"getServerStatistics",
			"failed to get dream tower bet count",
			result.Error,
		)
	}

	// 7. Get crash bet count.
	if result := session.Model(
		&models.CrashBet{},
	).Count(
		&crashBetCount,
	); result.Error != nil {
		return nil, utils.MakeError(
			"user_db_aggregator",
			"getServerStatistics",
			"failed to get crash bet count",
			result.Error,
		)
	}

	totalBets = coinflipBetCount + jackpotBetCount + dreamtowerBetCount + crashBetCount

	return &ServerStatisticsResult{
		TotalBets:    totalBets,
		TotalWagered: totalWagered,
		TotalProfit:  totalProfit,
	}, nil

}

// @Internal
// Check whether a user with provided wallet address is exists
// and get user info.
func getUserInfoByWalletAddress(walletAddress string) (*models.User, bool, error) {
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil, false, utils.MakeError(
			"user_db_aggregator",
			"getUserInfoByWalletAddress",
			"failed to get main session",
			err,
		)
	}

	var user models.User

	if result := session.Where(
		"wallet_address = ?",
		walletAddress,
	).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &user, false, nil
		}
		return &user, false, utils.MakeError(
			"user_db_aggregator",
			"getUserInfoByWalletAddress",
			"failed to get user info by wallet",
			result.Error,
		)
	}

	return &user, true, nil
}

// @Internal
// Get user info by name
func getUserInfoByName(name string) *models.User {
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	userInfo := models.User{}
	if result := session.Preload(
		"Statistics",
	).Select(
		"id",
		"name",
		"wallet_address",
		"role",
		"avatar",
		"banned",
		"private_profile",
	).Where(
		"LOWER(name) = ?",
		strings.ToLower(name),
	).First(
		&userInfo,
	); result.Error != nil {
		return nil
	}

	return &userInfo
}

// @Internal
// Get user count by IP address
func getUserCountByIPAddress(ipAddress string) (int64, error) {
	var count int64
	session, err := db_aggregator.GetSession()
	if err != nil {
		return count, utils.MakeError(
			"user_db_aggregator",
			"getUserCountByIPAddress",
			"failed to get main session",
			err,
		)
	}

	if result := session.Model(
		&models.User{},
	).Where(
		"ip_address = ?",
		ipAddress,
	).Count(
		&count,
	); result.Error != nil {
		return count, utils.MakeErrorWithCode(
			"user_db_aggregator",
			"getUserCountByIPAddress",
			"failed to get user count by IP address",
			ErrCodeFailedToCountUserByIP,
			result.Error,
		)
	}
	return count, nil
}

// @Internal
// Create a new user.
func createUser(request CreateUserRequest) error {
	// 1. Get main session from DB.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"user_db_aggregator",
			"createUser",
			"failed to get main session",
			err,
		)
	}

	// 2. Create a new record to `users` table.
	var newUser models.User = models.User(request)
	if result := session.Create(
		&newUser,
	); result.Error != nil {
		return utils.MakeErrorWithCode(
			"user_db_aggregator",
			"createUser",
			"failed to create a new user",
			ErrCodeFailedToCreateNewUser,
			result.Error,
		)
	}

	// 3. Generate rakeback for created user.
	if _, err := db_aggregator.GenerateRakebackInfo(
		db_aggregator.User(
			newUser.ID,
		),
	); err != nil {
		return utils.MakeErrorWithCode(
			"user_db_aggregator",
			"createUser",
			"failed to generate rakeback info for new user",
			ErrCodeFailedToGenerateRakeback,
			err,
		)
	}

	// 4. Create statistics record for the new user.
	var statistics = models.Statistics{
		UserID: newUser.ID,
	}
	if result := session.Create(
		&statistics,
	); result.Error != nil {
		return utils.MakeErrorWithCode(
			"user_db_aggregator",
			"createUser",
			"failed to create statistcs for new user",
			ErrCodeFailedToCreateStatistics,
			result.Error,
		)
	}

	// 5. Init seed pair for the new user.
	if _, err := seed.InitUserSeedPair(
		db_aggregator.User(
			newUser.ID,
		),
	); err != nil {
		return utils.MakeErrorWithCode(
			"user_db_aggregator",
			"createUser",
			"failed to initialize seed pair for user",
			ErrCodeFailedToInitSeedPair,
			err,
		)
	}
	return nil
}

// @Internal
// Save user info
func saveUser(
	user *models.User,
	updates map[string]interface{},
) error {
	// 1. Get main session from DB.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return utils.MakeError(
			"user_db_aggregator",
			"saveUser",
			"failed to get main session",
			err,
		)
	}

	// 2. Save user info.
	if result := session.Model(
		user,
	).Clauses(
		clause.Returning{},
	).Updates(
		updates,
	); result.Error != nil {
		return utils.MakeErrorWithCode(
			"user_db_aggregator",
			"saveUser",
			"failed to save user info",
			ErrCodeFailedToSaveUser,
			result.Error,
		)
	}

	return nil
}

// @Internal
// Get user info with balances
func getUserInfoWithBalances(userID uint) *models.User {
	session, err := db_aggregator.GetSession()
	if err != nil {
		return nil
	}

	userInfo := models.User{}
	if result := session.Preload(
		"Wallet.Balance.ChipBalance",
	).Preload(
		"Wallet.Balance.NftBalance",
	).First(
		&userInfo,
		userID,
	); result.Error != nil {
		return nil
	}

	return &userInfo
}

// @Internal
// Get nft details from mint addresses
func getNftDetailsFromMintAddresses(
	mintAddresses []string,
) (int64, []types.NftDetails) {
	// 1. Get main session from DB.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return 0, nil
	}

	// 2. Get deposited nfts with mint addresses
	var totalPrice int64
	depositedNfts := []models.DepositedNft{}
	session.Where(
		"mint_address IN ?",
		mintAddresses,
	).Find(
		&depositedNfts,
	)

	// 3. Convert deposited nfts to nft details.
	nftDetails := []types.NftDetails{}
	for _, nft := range depositedNfts {
		var collection models.NftCollection
		if result := session.First(
			&collection,
			nft.CollectionID,
		); result.Error != nil {
			return 0, nil
		}

		nftDetails = append(nftDetails, types.NftDetails{
			Name:            nft.Name,
			MintAddress:     nft.MintAddress,
			Image:           nft.Image,
			CollectionName:  collection.Name,
			CollectionImage: collection.Image,
			Price:           collection.FloorPrice,
		})
		totalPrice += collection.FloorPrice
	}
	return totalPrice, nftDetails
}

/*
* @Internal
* Returns count of users have similar names.
 */
func countSimilarNames(nameLike string) int {
	// 1. Validate parameter
	if nameLike == "" {
		return -1
	}

	// 2. Get main session from DB.
	session, err := db_aggregator.GetSession()
	if err != nil {
		return -1
	}

	// 3. Count users with similar names.
	var count int64
	session.Where(
		"name like ?",
		nameLike+"%",
	).Count(&count)

	return int(count)
}
