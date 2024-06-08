package utils

import (
	"math"

	"github.com/Duelana-Team/duelana-v1/config"
)

/*
* @External
* Dividing decimals to get chip amount from additional decimals.
 */
func ConvertBalanceToChip(value int64) int64 {
	return value / int64(math.Pow10(config.BALANCE_DECIMALS))
}

/*
* @External
* Multipling additional decimals to get balance from chip amount.
*/
func ConvertChipToBalance(value int64) int64 {
	return value * int64(math.Pow10(config.BALANCE_DECIMALS))
}

func GetMinimumChipInBalance() int64 {
	return int64(1000)
}
