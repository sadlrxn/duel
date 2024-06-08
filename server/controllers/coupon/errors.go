package coupon

// Error code range: #101xxx
const ErrCodeBase = "#101"
const ErrCodeInvalidParameter = ErrCodeBase + "001"
const ErrCodeAccessUserNamesMissingForCreate = ErrCodeBase + "002"
const ErrCodeLimitUserCountMissingForCreate = ErrCodeBase + "003"
const ErrCodeInvalidBalanceForCreate = ErrCodeBase + "004"
const ErrCodeAccessUserNameNotExistForCreate = ErrCodeBase + "005"
const ErrCodeCouponCodeNotFound = ErrCodeBase + "006"
const ErrCodeCouponAlreadyClaimed = ErrCodeBase + "007"
const ErrCodeCouponNotAllowedToClaim = ErrCodeBase + "008"
const ErrCodeCouponClaimReachedLimit = ErrCodeBase + "009"
const ErrCodeNotActiveCodeForExchange = ErrCodeBase + "010"
const ErrCodeNotReachingExchangeWager = ErrCodeBase + "011"
const ErrCodeExchangeTransactionFailure = ErrCodeBase + "012"
const ErrCodeAlreadyExistingActiveCoupon = ErrCodeBase + "013"
const ErrCodeInsufficientBonusBalance = ErrCodeBase + "014"
const ErrCodeZeroBonusBalance = ErrCodeBase + "015"
const ErrCodeInsufficientAdminBalance = ErrCodeBase + "016"
const ErrCodeCouponShortcutDuplicated = ErrCodeBase + "017"
const ErrCodeCouponShortcutNotFound = ErrCodeBase + "018"
const ErrCodeExistingPlayingRounds = ErrCodeBase + "019"
const ErrCodeMissingRequiredAffiliate = ErrCodeBase + "020"

// API Response errors
const ErrResponseCouponAlreadyHasActiveCode = 140001
const ErrResponseCouponNotFoundCode = 140002
const ErrResponseCouponInvalidPermission = 140003
const ErrResponseCouponClaimLimitExceed = 140004
const ErrResponseCouponAlreadyClaimedCode = 140005
const ErrResponseCouponNotReachedWagerLimit = 140006
const ErrResponseCouponInsufficientAdminBalance = 140007
const ErrResponseCouponExistingPlayingRounds = 140008
const ErrResponseCouponMssingRequiredAffiliate = 140009
