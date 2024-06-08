package solana

import "github.com/gagliardetto/solana-go"

type InitParam struct {
	TreasuryBs58 string
	Cluster      Cluster
	RpcUrl       string
}

type Cluster string

const (
	ClusterTestNet     Cluster = "testnet"
	ClusterDevNet      Cluster = "devnet"
	ClusterMainNetBeta Cluster = "mainnet-beta"
)

type TransactionType string

const (
	// TransactionSolDeposit  TransactionType = "transaction-sol-deposit"
	TransactionSplDeposit TransactionType = "transaction-spl-deposit"
	TransactionNftDeposit TransactionType = "transaction-nft-deposit"
	// TransactionSolWithdraw TransactionType = "transaction-sol-withdraw"
	TransactionSplWithdraw TransactionType = "transaction-spl-withdraw"
	TransactionNftWithdraw TransactionType = "transaction-nft-withdraw"
	TransactionNothing     TransactionType = "transaction-nothing"
)

type SplTokenType string

const (
	SolSplToken  SplTokenType = "sol_spl"
	BonkSplToken SplTokenType = "bonk_spl"
	UsdcSplToken SplTokenType = "usdc_spl"
	BokuSplToken SplTokenType = "boku_spl"
)

type SplTokenMeta struct {
	Type                 SplTokenType     `json:"type"`
	MintAddress          solana.PublicKey `json:"mintAddress"`
	TreasuryTokenAccount solana.PublicKey `json:"treasuryTokenAccount"`
	Decimals             int              `json:"decimals"`
	Keyword              string           `json:"keyword"`
	Image                string           `json:"image"`
}

type DecodedTransaction struct {
	TransactionType TransactionType
	Participant     string
	Lamports        uint64
	Nfts            []string
	Failed          bool
	SplToken        *SplTokenMeta
}

type InstructionType string

const (
	InstructionSolTransfer          InstructionType = "instruction-sol-transfer"
	InstructionTokenTransfer        InstructionType = "instruction-token-transfer"
	InstructionTokenTransferChecked InstructionType = "instruction-token-transfer-checked"
	InstructionTokenCreateAccount   InstructionType = "instruction-token-create-account"
	InstructionUnknown              InstructionType = "instruction-unknown"
)

type DecodedInstruction struct {
	InstructionType InstructionType
	Meta            interface{}
}

type AnalyseNftTransferResult struct {
	From solana.PublicKey
	To   solana.PublicKey
	Nfts []string
}

type SendLamportsRequest struct {
	To       string
	Lamports uint64
}

type SendNftsRequest struct {
	To   string
	Nfts []string
}

type SendSplTokenRequest struct {
	To     string
	Mint   string
	Amount uint64
}
