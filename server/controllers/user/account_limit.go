package user

import (
	"errors"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/utils"
)

// @Internal
// Check whether the account limit is reached with the ipAddress.
func checkIpAccountLimit(ipAddress string) (bool, error) {
	if ipAddress == "" {
		return false, utils.MakeErrorWithCode(
			"user_account_limit",
			"checkIpAccountLimit",
			"invalid parameter",
			ErrCodeInvalidParameter,
			errors.New("provided ipAddress is empty string"),
		)
	}

	limit := getAccountLimitForSpecificIp(ipAddress)

	count, err := getUserCountByIPAddress(ipAddress)
	if err != nil {
		return false, utils.MakeError(
			"user_account_limit",
			"checkIpAccountLimit",
			"failed to count user by IP address",
			err,
		)
	}

	return count >= int64(limit), nil
}

// @Internal
// Set account limit for specific ip address.
func setAccountLimitForSpecificIp(ipAddress string, limit uint) error {
	if ipAddress == "" {
		return utils.MakeError(
			"user_account_limit",
			"setAccountLimitForSpecificIp",
			"invalid parameter",
			errors.New("provided ipAddress is empty string"),
		)
	}
	config.ACCOUNT_LIMIT_FOR_SPECIFIC_IP[ipAddress] = limit
	return nil
}

// @Internal
// Get account limit for specific ip address. Returns base limit in case
// the ip address is not registered in the specific ip map.
func getAccountLimitForSpecificIp(ipAddress string) uint {
	limit, prs := config.ACCOUNT_LIMIT_FOR_SPECIFIC_IP[ipAddress]
	if prs {
		return limit
	}
	return config.ACCOUNT_LIMIT_PER_IP
}
