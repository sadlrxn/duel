package types

type NftDetails struct {
	Name            string `json:"name"`
	MintAddress     string `json:"mintAddress"`
	Image           string `json:"image"`
	CollectionName  string `json:"collectionName"`
	CollectionImage string `json:"collectionImage"`
	Price           int64  `json:"price"`
}

type PaymentSubscriptionData struct {
	Address          string `json:"address"`
	TxID             string `json:"txId"`
	BlockNumber      uint64 `json:"blockNumber"`
	Type             string `json:"type"`
	Amount           string `json:"amount"`
	SubscriptionType string `json:"subscriptionType"`
}

type PaymentSubscriptionDataV2 struct {
	TxID string `json:"txId"`
}

type TransactionDataMessage struct {
	AccountKeys     []string `json:"accountKeys"`
	RecentBlockhash string   `json:"recentBlockhash"`
}

type TransactionUiTokenAmount struct {
	Amount   string  `json:"amount"`
	Decimals uint    `json:"decimals"`
	UiAmount float32 `json:"uiAmount"`
}

type TransactionTokenBalance struct {
	AccountIndex  uint                     `json:"accountIndex"`
	Mint          string                   `json:"mint"`
	Owner         string                   `json:"owner"`
	ProgramID     string                   `json:"programId"`
	UiTokenAmount TransactionUiTokenAmount `json:"uiTokenAmount"`
}

type TransactionDataMeta struct {
	Error             interface{}               `json:"err"`
	Fee               uint64                    `json:"fee"`
	PostBalances      []uint64                  `json:"postBalances"`
	PreBalances       []uint64                  `json:"preBalances"`
	PostTokenBalances []TransactionTokenBalance `json:"postTokenBalances"`
	PreTokenBalances  []TransactionTokenBalance `json:"preTokenBalances"`
}

type TransactionDataTransaction struct {
	Message    TransactionDataMessage `json:"message"`
	Signatures []string               `json:"signatures"`
}

type TransactionDataResult struct {
	BlockTime   uint64                     `json:"blockTime"`
	Meta        TransactionDataMeta        `json:"meta"`
	Slot        uint64                     `json:"slot"`
	Transaction TransactionDataTransaction `json:"transaction"`
}
