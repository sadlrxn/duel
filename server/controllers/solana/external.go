package solana

import "github.com/gagliardetto/solana-go/rpc"

func Initialize(param *InitParam) error {
	return initialize(param)
}

func DecodeTransactionType(txHashBs58 string) (*DecodedTransaction, error) {
	return decodeTransactionType(txHashBs58)
}

func SendLamports(param *SendLamportsRequest) (string, error) {
	return sendLamports(param)
}

func SendSplTokens(param *SendSplTokenRequest) (string, error) {
	return sendSplTokens(param)
}

func SendNfts(param *SendNftsRequest) (string, error) {
	return sendNfts(param)
}

func SupportedSpls() []SplTokenMeta {
	return supportedSpls()
}

func GetTransactionResult(txHashBs58 string) (*rpc.GetTransactionResult, error) {
	return getTransaction(txHashBs58)
}

func IsSolSplMeta(meta SplTokenMeta) bool {
	return isSolSplMeta(meta)
}
