package seed

import (
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
)

func BorrowUserSeedPair(user db_aggregator.User) (*models.SeedPair, error) {
	return borrowUserSeedPair(user)
}

func ReturnUserSeedPair(user db_aggregator.User, seedPairID uint) error {
	return returnUserSeedPair(user, seedPairID)
}

func RotateUserSeedPair(
	user db_aggregator.User,
	serverSeedHash string,
	clientSeed string,
) (*models.SeedPair, error) {
	return rotateUserSeedPair(user, serverSeedHash, clientSeed)
}

func InitUserSeedPair(user db_aggregator.User) (*models.SeedPair, error) {
	return initUserSeedPair(user)
}

func GetActiveUserSeedPair(user db_aggregator.User) (*models.SeedPair, error) {
	return getActiveUserSeedPair(user)
}

func GetExpiredUserSeedPairs(
	user db_aggregator.User,
	skip uint,
	limit uint,
) ([]models.SeedPair, error) {
	return getExpiredUserSeedPairs(user, skip, limit)
}

func GetUnhashedServerSeed(hashedServerSeed string) (string, error) {
	return getUnhashedServerSeed(hashedServerSeed)
}
