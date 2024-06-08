package coupon

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
)

func TestRequiredAffiliate(t *testing.T) {
	db := tests.InitMockDB(true, true)
	if db == nil {
		t.Fatal("failed to get mock db")
	}
	db_aggregator.Initialize(db)

	if err := db.Create(
		&[]models.User{
			{
				Name: "creator",
			},
			{
				Name: "having",
			},
			{
				Name: "notHaving",
			},
		},
	).Error; err != nil {
		t.Fatalf("failed to create mock user: %v", err)
	}

	if err := db.Create(
		&models.Affiliate{
			Code:      "affiliate-code",
			CreatorID: 1,
			ActiveAffiliates: []models.ActiveAffiliate{
				{
					UserID: 2,
				},
			},
		},
	).Error; err != nil {
		t.Fatalf("failed to create affiliate code: %v", err)
	}

	affiliateCode := "affiliate-code"
	couponCode, err := Create(CreateCouponRequest{
		Type: models.CouponForSpecUsers,
		AccessUserNames: &[]string{
			"having", "notHaving",
		},
		Balance:               100 * 100000,
		RequiredAffiliateCode: &affiliateCode,
	})
	if err != nil {
		t.Fatalf("failed to create coupon code: %v", err)
	}

	if _, err := Claim(
		2, couponCode.String(),
	); err != nil {
		t.Fatalf("failed to claim coupon code: %v", err)
	}

	if _, err := Claim(
		3, couponCode.String(),
	); !utils.IsErrorCode(
		err,
		ErrCodeMissingRequiredAffiliate,
	) {
		t.Fatalf("should be failed on missing required affiliate: %v", err)
	}
}
