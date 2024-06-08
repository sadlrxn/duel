package crash

// Error code range: #106xxx
const ErrCodeBase = "#106"
const ErrCodeInvalidParameter = ErrCodeBase + "001"
const ErrCodeEmptyRound = ErrCodeBase + "002"
const ErrCodeNotPreparingRound = ErrCodeBase + "003"
const ErrCodeNotPreparingStatus = ErrCodeBase + "004"
const ErrCodeNotBettingRound = ErrCodeBase + "005"
const ErrCodeNotBettingStatus = ErrCodeBase + "006"
const ErrCodeNotPendingRound = ErrCodeBase + "007"
const ErrCodeNotPendingStatus = ErrCodeBase + "008"
const ErrCodeNotRunningRound = ErrCodeBase + "009"
const ErrCodeNotRunningStatus = ErrCodeBase + "010"
const ErrCodeNotFoundFinishedRound = ErrCodeBase + "011"

// Error codes for cash specific: #1061xx
const ErrCodeNotStatusForCashIn = ErrCodeBase + "101"
const ErrCodeMaxBetCountExceed = ErrCodeBase + "102"
const ErrCodeLessThanMinBetAmount = ErrCodeBase + "103"
const ErrCodeMoreThanMaxBetAmount = ErrCodeBase + "104"
const ErrCodeBalanceTypeMismatching = ErrCodeBase + "105"
const ErrCodeInsufficientUserBalance = ErrCodeBase + "106"
const ErrCodeNotStatusForCashOut = ErrCodeBase + "107"
const ErrCodeInvalidBetForCashout = ErrCodeBase + "108"
const ErrCodeInsufficientPoolBalance = ErrCodeBase + "109"
