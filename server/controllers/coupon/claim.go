package coupon

import (
	"fmt"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

// @External
// For client.
// Claims bonus balance through coupon code provided by admin.
// Returns claimed bonus balance and error object.
// Returns error on
//   - Already have active coupon not expired. ErrCode: ErrCodeAlreadyExistingActiveCoupon
//   - Provided code not found in coupon table. ErrCode: ErrCodeCouponCodeNotFound
//   - Provided code is already claimed by this user. ErrCode: ErrCodeCouponAlreadyClaimed
//   - Provided code is not allowed by this user to be claimed. ErrCode: ErrCodeCouponNotAllowedToClaim
//   - Provided code is already reached Claim limit. ErrCode: ErrCodeCouponClaimReachedLimit
//   - For the expired coupon codes(14 days passed after creation) will be handled by
//     ErrCodeCouponCodeNotFound error.
func Claim(userID uint, codeLike string) (int64, error) {
	// 1. Validate parameter.
	if userID == 0 ||
		codeLike == "" {
		return 0, utils.MakeErrorWithCode(
			"coupon_claim",
			"Claim",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf("userID: %d, code: %v", userID, codeLike),
		)
	}

	// 2. Get coupon code as UUID.
	code, err := getCouponAsUUID(codeLike)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"failed to get coupon as UUID",
			fmt.Errorf(
				"codeLike: %s, err: %v",
				codeLike, err,
			),
		)
	}

	// 3. Start a session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"failed to start session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 4. Try retrieve active coupon record.
	if activeCoupon, err := lockAndRetrieveActiveCoupon(userID, sessionId); err == nil {
		return 0, utils.MakeErrorWithCode(
			"coupon_claim",
			"Claim",
			"already have active coupon",
			ErrCodeAlreadyExistingActiveCoupon,
			fmt.Errorf(
				"userId: %d, codeParam: %v, activeCoupon: %v",
				userID,
				code,
				activeCoupon,
			),
		)
	} else if !utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"failed to try get active coupon",
			err,
		)
	}

	// 5. Lock and retrieve coupon record.
	coupon, err := lockAndRetrieveCoupon(code, sessionId)
	if err != nil {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"failed to retrieve coupon",
			fmt.Errorf("code: %v, error: %v", code, err),
		)
	}

	// 6. Verify about coupon accessibility.
	if len(coupon.ClaimedCoupons) >= int(coupon.AccessUserLimit) {
		return 0, utils.MakeErrorWithCode(
			"coupon_claim",
			"Claim",
			"coupon claim reached limit",
			ErrCodeCouponClaimReachedLimit,
			fmt.Errorf(
				"coupon: %v",
				coupon,
			),
		)
	}
	if coupon.Type == models.CouponForSpecUsers {
		isAccessUser := false
		for _, accessUserID := range coupon.AccessUserIDs {
			if uint(accessUserID) == userID {
				isAccessUser = true
				break
			}
		}
		if !isAccessUser {
			return 0, utils.MakeErrorWithCode(
				"coupon_claim",
				"Claim",
				"coupon not allowed to claim",
				ErrCodeCouponNotAllowedToClaim,
				fmt.Errorf(
					"coupon: %v, userID: %d",
					coupon, userID,
				),
			)
		}
	} else if coupon.Type == models.CouponForLimitUsers {
		// Nothing to check here for now since the limit is checked for all cases.
	} else {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"unrecognized coupon type",
			fmt.Errorf("type: %s", coupon.Type),
		)
	}

	isClaimed := false
	for _, claimedCoupon := range coupon.ClaimedCoupons {
		if claimedCoupon.ClaimedUserID == userID {
			isClaimed = true
			break
		}
	}
	if isClaimed {
		return 0, utils.MakeErrorWithCode(
			"coupon_claim",
			"Claim",
			"already claimed by this user",
			ErrCodeCouponAlreadyClaimed,
			fmt.Errorf(
				"coupon: %v, userID: %d",
				coupon, userID,
			),
		)
	}

	// 7. Check for missing required affiliate code activation.
	if missingRequiredAffiliateCode(userID, coupon) {
		return 0, utils.MakeErrorWithCode(
			"coupon_claim",
			"Claim",
			"missing required affiliate code",
			ErrCodeMissingRequiredAffiliate,
			fmt.Errorf(
				"coupon: %v, userID: %d",
				coupon, userID,
			),
		)
	}

	// 8. Create a new claimed coupon record.
	newActiveCoupon := models.ClaimedCoupon{
		CouponID:      code,
		ClaimedUserID: userID,
	}
	if err := createClaimedCouponUnchecked(
		&newActiveCoupon,
		sessionId,
	); err != nil {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"failed to create claimed coupon record",
			fmt.Errorf(
				"newClaimedCoupon: %v, err: %v",
				newActiveCoupon, err,
			),
		)
	}

	// 9. Update balance and record transaction.
	if _, err := performTransactionInSession(
		CouponTransactionRequest{
			Type:          models.CpTxClaimCode,
			UserID:        userID,
			Balance:       coupon.BonusBalance,
			ToBeConfirmed: true,
		},
		sessionId,
	); err != nil {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"failed to mint bonus balance for claiming",
			err,
		)
	}

	// 10. Commit session.
	if err := db_aggregator.CommitSession(sessionId); err != nil {
		return 0, utils.MakeError(
			"coupon_claim",
			"Claim",
			"failed to commit session",
			err,
		)
	}

	return coupon.BonusBalance, nil
}

/**
* @Internal
* Checks whether the user is missing required affiliate code activation.
* If the targeted coupon is not requiring any affiliate code activation,
* returns false.
 */
func missingRequiredAffiliateCode(
	userID uint,
	coupon *models.Coupon,
) bool {
	if coupon == nil ||
		coupon.RequiredAffiliate == nil {
		return false
	}

	for _, activeUser := range coupon.RequiredAffiliate.ActiveAffiliates {
		if activeUser.UserID == userID {
			return false
		}
	}

	return true
}
