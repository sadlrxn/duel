package prelude

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type InitialDuelUser struct {
	ID            uint
	Name          string
	WalletAddress string
}

func buildUserAccount(userInfo InitialDuelUser) models.User {
	return models.User{
		Model: gorm.Model{
			ID: userInfo.ID,
		},
		Name:          userInfo.Name,
		WalletAddress: userInfo.WalletAddress,
		Avatar:        "https://duelana-bucket.s3.us-east-2.amazonaws.com/avatar/default.png",
		Wallet: models.Wallet{
			Balance: models.Balance{
				ChipBalance: &models.ChipBalance{},
				NftBalance: &models.NftBalance{
					Balance: pq.StringArray{},
				},
			},
		},
	}
}

func getInitialUsers() []InitialDuelUser {
	return []InitialDuelUser{
		{
			ID:            config.JACKPOT_TEMP_ID,
			Name:          "JP_TEMP",
			WalletAddress: "5akt2sxYcVaDvrG89GNanTmSnu4yaQvX59uaBKYKBRu6",
		},
		{
			ID:            config.JACKPOT_FEE_ID,
			Name:          "JP_FEE",
			WalletAddress: "7BM7wQPaXe4pKkinqjqf3wu7awLaHnZHQQ8cLc2jKFsG",
		},
		{
			ID:            config.GRAND_JACKPOT_TEMP_ID,
			Name:          "GJ_TEMP",
			WalletAddress: "9YmQoL4k9m7EoXdLvHjNwyjDWL9AK26kHu1NiEoqgURi",
		},
		{
			ID:            config.GRAND_JACKPOT_FEE_ID,
			Name:          "GJ_FEE",
			WalletAddress: "3GqgyTFrFwuf6rRCo55RQhjVQhRU9yAePEa9msZ7GsYc",
		},
		{
			ID:            config.COINFLIP_TEMP_ID,
			Name:          "CF_TEMP",
			WalletAddress: "Fw5RE6fWw91yuqD4WwJRri3mia8GSjZvq6e1JqL8QPQJ",
		},
		{
			ID:            config.COINFLIP_FEE_ID,
			Name:          "CF_FEE",
			WalletAddress: "9Xgx8Nj4Sxd9MeqJCBv9u83QkdME4pVaJucvvtE7oAgw",
		},
		{
			ID:            config.COINFLIP_BOT_ID,
			Name:          "DuelBot",
			WalletAddress: "DqbW8Z7SckhDemc1Xows75JfU92HE1hBtsrhNj2NpE4g",
		},
		{
			ID:            config.DREAMTOWER_TEMP_ID,
			Name:          "DT_TEMP",
			WalletAddress: "9jUrSQuoNdya6SyrtuCPy9MqpdrA1j2dWaZdSoh9aTHV",
		},
		{
			ID:            config.DREAMTOWER_FEE_ID,
			Name:          "DT_FEE",
			WalletAddress: "2EusYSbaekKNv5DVxYmZ4tHqDaKDou6bzY3FYoyeM8aD",
		},
		{
			ID:            config.DUEL_BOT_STAKE_ID,
			Name:          "DB_STAKE",
			WalletAddress: "AXSmFvYoXCtWBs1vFim6stuKM6eMAqa6FRZ4WhVTg7Sy",
		},
		{
			ID:            config.COUPON_TEMP_ID,
			Name:          "CP_TEMP",
			WalletAddress: "HZq3sjwhj2aWatuyJgpTdF5idMYYJdfBnSmsh1SZL8q6",
		},
		{
			ID:            config.CRASH_TEMP_ID,
			Name:          "CH_TEMP",
			WalletAddress: "B1fuc1pmCdsLp2ceehs9NBMth9idrpmFRrWevG7wmq5x",
		},
		{
			ID:            config.CRASH_FEE_ID,
			Name:          "CH_FEE",
			WalletAddress: "475ALhTThzNsqeA46sD55181KxVgZsGRbqmxm3ERKN7B",
		},
		{
			ID:            config.DAILY_RACE_TEMP_ID,
			Name:          "DR_TEMP",
			WalletAddress: "EohHXvADJy3jFTWsNiTjmWhtCEGfkvgFdq6JsNjdZt96",
		},
		{
			ID:            config.WEEKLY_RAFFLE_TEMP_ID,
			Name:          "WR_TEMP",
			WalletAddress: "A34Rv49byu8ebEY6LtLDp9hzLT3iW3uNWHgQq1Hzsrb1",
		},
	}
}

/* Initializes main wallets on fresh DB.
* 1. JP_TEMP,	1001	5akt2sxYcVaDvrG89GNanTmSnu4yaQvX59uaBKYKBRu6
* 2. JP_FEE, 	1002	7BM7wQPaXe4pKkinqjqf3wu7awLaHnZHQQ8cLc2jKFsG
* 3. GJ_TEMP, 	1003	9YmQoL4k9m7EoXdLvHjNwyjDWL9AK26kHu1NiEoqgURi
* 4. GJ_FEE, 	1004	3GqgyTFrFwuf6rRCo55RQhjVQhRU9yAePEa9msZ7GsYc
* 5. CF_TEMP, 	1005	Fw5RE6fWw91yuqD4WwJRri3mia8GSjZvq6e1JqL8QPQJ
* 6. CF_FEE, 	1006	9Xgx8Nj4Sxd9MeqJCBv9u83QkdME4pVaJucvvtE7oAgw
* 7. DuelBot,	1007	DqbW8Z7SckhDemc1Xows75JfU92HE1hBtsrhNj2NpE4g
* 8. DT_TEMP, 	1008	9jUrSQuoNdya6SyrtuCPy9MqpdrA1j2dWaZdSoh9aTHV
* 9. DT_FEE, 	1009	2EusYSbaekKNv5DVxYmZ4tHqDaKDou6bzY3FYoyeM8aD
* 10.DB_STAKE, 	10001	AXSmFvYoXCtWBs1vFim6stuKM6eMAqa6FRZ4WhVTg7Sy
* 11.CP_TEMP,   100001  HZq3sjwhj2aWatuyJgpTdF5idMYYJdfBnSmsh1SZL8q6
* 12.CH_TEMP,	100002 	B1fuc1pmCdsLp2ceehs9NBMth9idrpmFRrWevG7wmq5x
* 13.CH_FEE,	100003	475ALhTThzNsqeA46sD55181KxVgZsGRbqmxm3ERKN7B
* 14.DR_TEMP,	100004	EohHXvADJy3jFTWsNiTjmWhtCEGfkvgFdq6JsNjdZt96
* 14.WR_TEMP,	100005	A34Rv49byu8ebEY6LtLDp9hzLT3iW3uNWHgQq1Hzsrb1
 */
func InitDuelMainUsers(db *gorm.DB) error {
	initialUsers := getInitialUsers()

	initialUserIDs := []uint{}
	for _, user := range initialUsers {
		initialUserIDs = append(initialUserIDs, user.ID)
	}

	var existUsers []models.User
	if result := db.Where(
		"id in ?",
		initialUserIDs,
	).Order(
		"id",
	).Find(
		&existUsers,
	); result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	addingUsers := []models.User{}
	j := 0
	for _, user := range existUsers {
		for ; j < len(initialUsers) && initialUsers[j].ID < user.ID; j++ {
			addingUsers = append(addingUsers, buildUserAccount(initialUsers[j]))
		}
		if j == len(initialUsers) {
			break
		}
		j++
	}
	for ; j < len(initialUsers); j++ {
		addingUsers = append(addingUsers, buildUserAccount(initialUsers[j]))
	}

	if len(addingUsers) == 0 {
		return nil
	}

	result := db.Create(&addingUsers)
	return result.Error
}
