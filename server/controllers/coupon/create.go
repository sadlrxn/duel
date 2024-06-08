package coupon

import (
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// @External
// For admin.
// Creates coupon code according to the request info.
// Returns newly generated coupon code.
// Throws error in case of
//   - Access user names are missing for spec user coupon type.
//   - Limit user count is missing for limit user coupon type.
//   - Balance is 0 or negative.
func Create(createRequest CreateCouponRequest) (uuid.UUID, error) {
	// 1. Validate parameters
	if createRequest.Type == models.CouponForSpecUsers &&
		((createRequest.AccessUserNames == nil ||
			len(*createRequest.AccessUserNames) == 0) &&
			(createRequest.AccessUserIDs == nil ||
				len(*createRequest.AccessUserIDs) == 0)) {
		return uuid.Nil, utils.MakeErrorWithCode(
			"coupon",
			"createCoupon",
			"invalid parameter",
			ErrCodeAccessUserNamesMissingForCreate,
			errors.New("type is CouponForSpecUsers, but missing user names"),
		)
	}
	if createRequest.Type == models.CouponForLimitUsers &&
		(createRequest.AccessUserLimit == nil ||
			*createRequest.AccessUserLimit == 0) {
		return uuid.Nil, utils.MakeErrorWithCode(
			"coupon",
			"createCoupon",
			"invalid parameter",
			ErrCodeLimitUserCountMissingForCreate,
			errors.New("type is CouponForLimitUsers, but missing user limit"),
		)
	}
	if createRequest.Balance <= 0 {
		return uuid.Nil, utils.MakeErrorWithCode(
			"coupon",
			"createCoupon",
			"invalid parameter",
			ErrCodeInvalidBalanceForCreate,
			errors.New("should provide balance greater than 0"),
		)
	}

	// 2. Retrieve main session
	session, err := db_aggregator.GetSession()
	if err != nil {
		return uuid.Nil, utils.MakeError(
			"coupon",
			"createCoupon",
			"failed to retrieve session",
			err,
		)
	}

	// 3. Create a code and return
	newCoupon := models.Coupon{
		Type:         createRequest.Type,
		BonusBalance: createRequest.Balance,
	}
	if createRequest.Type == models.CouponForLimitUsers {
		newCoupon.AccessUserLimit = uint(*createRequest.AccessUserLimit)
	} else if createRequest.Type == models.CouponForSpecUsers {
		if createRequest.AccessUserIDs != nil {
			for _, userID := range *createRequest.AccessUserIDs {
				newCoupon.AccessUserIDs = append(
					newCoupon.AccessUserIDs,
					int64(userID),
				)
			}
		}
		if createRequest.AccessUserNames != nil {
			userInfos := []models.User{}
			if result := session.Where(
				"name in ?",
				*createRequest.AccessUserNames,
			).Find(&userInfos); result.Error != nil {
				return uuid.Nil, utils.MakeError(
					"coupon",
					"createCoupon",
					"failed to retrieve user info from names",
					result.Error,
				)
			}
			if len(userInfos) != len(*createRequest.AccessUserNames) {
				return uuid.Nil, utils.MakeErrorWithCode(
					"coupon",
					"createCoupon",
					"retrieved user info count mismatching with provided names",
					ErrCodeAccessUserNameNotExistForCreate,
					fmt.Errorf(
						"provided: %d, retrieved: %d",
						len(*createRequest.AccessUserNames),
						len(userInfos),
					),
				)
			}
			newCoupon.AccessUserIDs = pq.Int64Array{}
			for _, userInfo := range userInfos {
				newCoupon.AccessUserIDs = append(newCoupon.AccessUserIDs, int64(userInfo.ID))
			}
		}
		newCoupon.AccessUserLimit = uint(len(newCoupon.AccessUserIDs))
	}
	if createRequest.RequiredAffiliateCode != nil &&
		len(*createRequest.RequiredAffiliateCode) > 0 {
		newCoupon.RequiredAffiliateCode = createRequest.RequiredAffiliateCode
	}
	if result := session.Create(&newCoupon); result.Error != nil {
		return uuid.Nil, utils.MakeError(
			"coupon",
			"createCoupon",
			"failed to create new coupon",
			result.Error,
		)
	}
	if newCoupon.Code == nil {
		return uuid.Nil, utils.MakeError(
			"coupon",
			"createCoupon",
			"default uuid is not assigned",
			errors.New("coupon code is nil pointer"),
		)
	}

	return *newCoupon.Code, nil
}
