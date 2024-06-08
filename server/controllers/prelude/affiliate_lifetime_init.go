package prelude

import (
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

func shouldInitAffiliateLifetime() bool {
	const createTableQuery = `create table a_affiliate_lifetime_flag (
	is_initialized boolean default TRUE
);`
	const retrieveTableQuery = `select count(1) from pg_tables where tablename = 'a_affiliate_lifetime_flag';`

	session, err := db_aggregator.GetSession()
	if err != nil {
		log.LogMessage(
			"prelude_shouldInitAffiliateLifetime",
			"failed to retrieve main session",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return false
	}

	existence := int64(0)
	if err := session.Raw(retrieveTableQuery).Row().Scan(&existence); err != nil {
		log.LogMessage(
			"prelude_shouldInitAffiliateLifetime",
			"failed to execute retrieve table query",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		return false
	}
	if existence > 0 {
		return false
	}

	if result := session.Exec(createTableQuery); result.Error != nil {
		log.LogMessage(
			"prelude_shouldInitAffiliateLifetime",
			"failed to execute create table query",
			"error",
			logrus.Fields{
				"error": result.Error.Error(),
			},
		)
	}
	return true
}

func InitAffiliateLifetime() error {
	log.LogMessage(
		"prelude_affiliate_lifetime_InitAffiliateLifetime",
		"Initializing...",
		"info",
		logrus.Fields{},
	)
	if shouldInitAffiliateLifetime() {
		session, err := db_aggregator.GetSession()
		if err != nil {
			return utils.MakeError(
				"prelude",
				"InitAffiliateLifetime",
				"failed to get main db session",
				err,
			)
		}

		activeAffiliates := []models.ActiveAffiliate{}
		if result := session.Find(&activeAffiliates); result.Error != nil {
			return utils.MakeError(
				"prelude",
				"InitAffiliateLifetime",
				"failed to retrieve all active affiliates",
				err,
			)
		}

		for _, activeAffiliate := range activeAffiliates {
			if result := session.Create(&models.AffiliateLifetime{
				UserID:        activeAffiliate.UserID,
				AffiliateID:   activeAffiliate.AffiliateID,
				LastActivated: activeAffiliate.CreatedAt,
				IsActive:      true,
			}); result.Error != nil {
				log.LogMessage(
					"prelude_init_affiliate_lifetime",
					"failed to create new affiliate lifetime record",
					"error",
					logrus.Fields{
						"userId":        activeAffiliate.UserID,
						"affiliateId":   activeAffiliate.AffiliateID,
						"lastActivated": activeAffiliate.CreatedAt,
						"error":         result.Error.Error(),
					},
				)
			}
		}
	} else {
		log.LogMessage(
			"prelude_init_affiliate_lifetime",
			"already initialized affiliate lifetime",
			"info",
			logrus.Fields{},
		)
	}

	return nil
}
