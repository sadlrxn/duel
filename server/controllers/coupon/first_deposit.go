package coupon

import (
	"fmt"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

/**
* @External
* Performs first deposit bonus.
* If the depositAmount is higher than max first deposit bonus,
* only gives the maximum.
* Creates coupon code for only that user, and tries to redeem.
* For now only prints the error in case redeem fails due to currently
* activated coupon code.
* Return redeemed bonus balance amount and error object.
 */
func PerformFirstDepositBonus(
	userID uint,
	depositAmount int64,
) (int64, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		depositAmount <= 0 {
		return 0, utils.MakeError(
			"coupon_first_deposit",
			"PerformFirstDepositBonus",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, depositAmount: %d",
				userID, depositAmount,
			),
		)
	}

	// 2. Create coupon code for one spec user.
	if depositAmount > config.COUPON_MAXIMUM_FIRST_DEPOSIT_BONUS {
		depositAmount = config.COUPON_MAXIMUM_FIRST_DEPOSIT_BONUS
	}
	code, err := Create(
		CreateCouponRequest{
			Type:          models.CouponForSpecUsers,
			AccessUserIDs: &[]uint{userID},
			Balance:       depositAmount,
		},
	)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_first_deposit",
			"PerformFirstDepositBonus",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, bonusAmount: %d, err: %v",
				userID, depositAmount, err,
			),
		)
	}

	// 3. Try to redeem code.
	if bonus, err := Claim(
		userID, code.String(),
	); err != nil {
		return 0, utils.MakeError(
			"coupon_first_deposit",
			"PerformFirstDepositBonus",
			"failed to perform redeem coupon code",
			fmt.Errorf(
				"userID: %d, code: %s, err: %v",
				userID, code, err,
			),
		)
	} else {
		return bonus, nil
	}
}
