package user

// Error code range: #101xxx
const ErrCodeBase = "#101"

const ErrCodeInvalidParameter = ErrCodeBase + "001"
const ErrCodeFailedToCountUserByIP = ErrCodeBase + "002"
const ErrCodeAlreadyReachedAccountLimit = ErrCodeBase + "003"
const ErrCodeFailedToCreateNewUser = ErrCodeBase + "004"
const ErrCodeFailedToGenerateRakeback = ErrCodeBase + "005"
const ErrCodeFailedToCreateStatistics = ErrCodeBase + "006"
const ErrCodeFailedToInitSeedPair = ErrCodeBase + "007"
const ErrCodeFailedToSaveUser = ErrCodeBase + "008"
const ErrCodeFailedToGetTargetUserInfo = ErrCodeBase + "009"
const ErrCodeIsDuelHiddenUser = ErrCodeBase + "010"
const ErrCodeFailedToGetViewerUserInfo = ErrCodeBase + "011"
const ErrCodeInvalidUserName = ErrCodeBase + "012"
