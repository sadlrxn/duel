package crash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func calculateOutCome(
	serverSeed string,
	clientSeed string,
	houseEdge int64,
) float64 {
	// 1. Calculate random Uint64 value from seeds.
	randomUint := calculateRandomUint(serverSeed, clientSeed)

	// 2. This part results 5% instant crash.
	// `houseEdge` == 500, `hs` == 20
	var hs int64 = 10000 / houseEdge
	if randomUint%uint64(hs) == 0 {
		return 1
	}

	// 3. Calculate and return outcome.
	var h float64 = float64(randomUint)
	var e float64 = math.Pow(2, 52)

	return math.Floor((100*e-h)/(e-h)) / 100
}

func calculateRandomUint(
	serverSeed string,
	clientSeed string,
) uint64 {
	// 1. Generate SHA256 hash of `serverSeed` + `clientSeed`.
	sum := sha256.Sum256([]byte(serverSeed + clientSeed))

	// 2. Convert first 8 bytes of hash array as Uint64.
	randomUint := uint64(0)
	for i := 0; i < (52 / 8); i++ {
		randomUint = randomUint<<8 + uint64(sum[i])
	}
	randomUint = randomUint<<4 + (uint64(sum[52/8]&0b11110000) >> 4)

	// randomUint := binary.BigEndian.Uint64(sum[:7])
	return randomUint
}

/*
/* @Internal
/* Assumes that we already have server seeds in crash_rounds table.
*/
func determineClientSeed(
	clientSeed string,
	houseEdge int64,
	startIndex int,
) error {
	// 1. Fetch total round count from DB.
	count := getTotalCrashRoundCount()
	if count <= 0 {
		return utils.MakeError(
			"crash seed",
			"determineClientSeed",
			"no rounds exist",
			nil,
		)
	}

	log.LogMessage(
		"determineClientSeed",
		fmt.Sprintf(
			"*********** Starting to calculate %d outcomes. from %d ***********",
			count-int64(startIndex), startIndex,
		),
		"info",
		logrus.Fields{},
	)

	// 2. Record operation starting time.
	var startTime = time.Now()

	// 3. Calculate outcomes for each rounds.
	for i := startIndex; i <= int(count); i++ {

		// 3.1. Lock and retrieve `i` round.
		round, err := lockAndRetrieveCrashRound(
			uint(i),
			db_aggregator.MainSessionId(),
		)
		if err != nil {
			return utils.MakeError(
				"crash seed",
				"determineClientSeed",
				fmt.Sprintf("failed to lock and retrieve round: %d", i),
				err,
			)
		}

		// 3.2. Calculate `outcome` for the round `i`.
		outcome := calculateOutCome(round.Seed, clientSeed, houseEdge)

		// 3.3. Save round with calcualted `outcome`.
		if err := updateCrashRoundOutcome(
			round,
			outcome,
		); err != nil {
			return utils.MakeError(
				"crash seed",
				"determineClientSeed",
				fmt.Sprintf("failed to update round outcome: %d", i),
				err,
			)
		}

		log.LogMessage(
			"determineClientSeed",
			fmt.Sprintf(
				"%d / %d : %f",
				i,
				count,
				outcome,
			),
			"info",
			logrus.Fields{},
		)
	}

	// 4. Calculate total time taken for this operation.
	var endTime = time.Now()
	var timeTaken = endTime.Sub(startTime).String()

	log.LogMessage(
		"determineClientSeed",
		fmt.Sprintf(
			"*********** Finished to outcomes of %d rounds, %s taken. ***********",
			count,
			timeTaken,
		),
		"info",
		logrus.Fields{},
	)
	return nil
}

/*
/* @Internal
/* Generate whole seed chain for the crash.
/*  - Migrates `CrashRounds` table in DB.
/*  - Generate `count` of seeds one-by-one for each rounds.
/*  - Save rounds with generated seed & initial values.
*/
func generateSeedChainWithSalt(
	salt string,
	length int,
) (string, error) {
	// 1. Start operation with provided salt string.
	var temp = salt

	// 2. Auto migrate `crash_rounds` table before creating rounds.
	if err := autoMigrateCrashRound(); err != nil {
		return "", utils.MakeError(
			"crash seed",
			"generateSeedChainWithSalt",
			"failed to migrate table",
			err,
		)
	}

	// 3. Check total round count and return error if `count` != 0.
	if count := getTotalCrashRoundCount(); count != 0 {
		return "", utils.MakeError(
			"crash seed",
			"generateSeedChainWithSalt",
			"records already exists in crash round table",
			nil,
		)
	}

	log.LogMessage(
		"generateSeedChainWithSalt",
		fmt.Sprintf(
			"*********** Starting to generate %d hashes from '%s'. ***********",
			length,
			salt,
		),
		"info",
		logrus.Fields{},
	)

	// 4. Record operation starting time.
	var startTime = time.Now()

	// 5. Create `length` sized hash chain.
	for i := 0; i < length; i++ {
		// 5.1. Generate SHA256 hash of previous value.
		temp = generateDerivedHash(temp)

		// 5.2. Create crash round with generated seed.
		var round = models.CrashRound{
			Model: gorm.Model{
				ID: uint(length - i),
			},
			Seed: temp,
		}

		if err := createCrashRound(&round); err != nil {
			log.LogMessage(
				"generateSeedChainWithSalt",
				"Failed to save hash to DB.",
				"error",
				logrus.Fields{
					"ID":    round.ID,
					"Seed":  round.Seed,
					"error": err.Error(),
				},
			)
			return "", utils.MakeError(
				"crash seed",
				"generateSeedChainWithSalt",
				fmt.Sprintf("failed to create round: %d", round.ID),
				err,
			)
		}
		log.LogMessage(
			"generateSeedChainWithSalt",
			fmt.Sprintf(
				"%d / %d : %s",
				i+1,
				length,
				temp,
			),
			"info",
			logrus.Fields{},
		)
	}

	// 6. Calculate total time taken for this operation.
	var endTime = time.Now()
	var timeTaken = endTime.Sub(startTime).String()

	log.LogMessage(
		"generateSeedChainWithSalt",
		fmt.Sprintf(
			"*********** Finished to generate the whole hash chain, %s taken. ***********",
			timeTaken,
		),
		"info",
		logrus.Fields{},
	)

	return generateDerivedHash(temp), nil
}

/*
/* @Internal
/* Returns SHA256 hash of input string.
*/
func generateDerivedHash(str string) string {
	hashBytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hashBytes[:])
}
