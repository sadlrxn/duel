package coupon

import (
	"fmt"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/google/uuid"
)

/**
* @External
* Creates a short cut to couponID.
* Returns error when:
*  - Invalid parameter. `ErrCodeInvalidParameter`
*  - couponID is not found. `ErrCodeCouponCodeNotFound`
*  - Coupon shortcut duplicated for couponID or shortcut. `ErrCodeCouponShortcutDuplicated`
 */
func CreateShortcut(
	couponID uuid.UUID,
	shortcut string,
) error {
	// 1. Validate parameter.
	if couponID == uuid.Nil ||
		shortcut == "" {
		return utils.MakeErrorWithCode(
			"coupon_shortcut",
			"CreateShortcut",
			"invalid parameter",
			ErrCodeInvalidParameter,
			fmt.Errorf(
				"couponID: %v, shortcut: %s",
				couponID, shortcut,
			),
		)
	}

	// 2. Try to retrieve coupon record.
	if _, err := lockAndRetrieveCoupon(
		couponID,
		db_aggregator.MainSessionId(),
	); err != nil {
		return utils.MakeError(
			"coupon_shortcut",
			"CreateShortcut",
			"failed to retrieve coupon record",
			fmt.Errorf(
				"couponID: %v, err: %v",
				couponID, err,
			),
		)
	}

	// 3. Create coupon shortcut.
	if err := createCouponShortcutUnchecked(
		couponID,
		shortcut,
	); err != nil {
		return utils.MakeErrorWithCode(
			"coupon_shortcut",
			"CreateShortcut",
			"failed to create coupon record",
			ErrCodeCouponShortcutDuplicated,
			fmt.Errorf(
				"couponID: %v, shortcut: %s, err: %v",
				couponID, shortcut, err,
			),
		)
	}

	return nil
}

/**
* @External
* Deletes shortcut record.
* Returns error object for
*  - Not found shortcut. `ErrCodeCouponShortcutNotFound`
 */
func DeleteShortcut(
	shortcut string,
) error {
	if err := deleteCouponShortcutUnchecked(
		shortcut,
	); err != nil {
		return utils.MakeErrorWithCode(
			"coupon_shortcut",
			"DeleteShortcut",
			"failed to delete coupon shortcut",
			ErrCodeCouponShortcutNotFound,
			fmt.Errorf(
				"shortcut: %s, err: %v",
				shortcut, err,
			),
		)
	}

	return nil
}

/**
* @Internal
* Gets UUID form code from both UUID string and shortcut.
* Returns error for
*  - Not found shortcut. `ErrCodeCouponShortcutNotFound`
 */
func getCouponAsUUID(
	codeLike string,
) (uuid.UUID, error) {
	// 1. Return parsed result if codeLike is form of UUID.
	if code, err := uuid.Parse(codeLike); err == nil {
		return code, nil
	}

	// 2. Retrieve coupon shortcut record.
	if couponShortcut, err := retrieveCouponShortcut(
		codeLike,
	); err != nil {
		return uuid.Nil, utils.MakeErrorWithCode(
			"coupon_shortcut",
			"getCouponAsUUID",
			"failed to retrieve coupon shortcut",
			ErrCodeCouponShortcutNotFound,
			fmt.Errorf(
				"codeLike: %s, err: %v",
				codeLike, err,
			),
		)
	} else {
		return couponShortcut.CouponID, nil
	}
}
