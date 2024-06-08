package solana

import (
	"fmt"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// ===== Constants =====
// Default cluster.
const DEFAULT_CLUSTER = ClusterDevNet

// ===== Treasury Wallet =====
// Treasury wallet private key in base58 format.
var _treasuryBs58 string

// @Internal
// Setter for `_treasuryBs58`.
func setTreasuryBs58(treasuryBs58 string) {
	_treasuryBs58 = treasuryBs58
}

// @Internal
// Getter for `_treasuryBs58`.
func treasuryBs58() string {
	return _treasuryBs58
}

// @Internal
// Get KeyPair from `_treasuryBs58`.
func treasuryKeyPair() *solana.PrivateKey {
	treasuryKeyPair := solana.MustPrivateKeyFromBase58(treasuryBs58())
	return &treasuryKeyPair
}

// @Internal
// Get PublicKey from `_treasuryBs58`.
func treasuryPubKey() *solana.PublicKey {
	treasuryKeyPair := treasuryKeyPair()
	treasuryPubKey := treasuryKeyPair.PublicKey()
	return &treasuryPubKey
}

// @Internal
// Check whether a PubKey equals to treasury PubKey.
func equalPubKeyToTreasury(pubKey solana.PublicKey) bool {
	return treasuryPubKey().Equals(pubKey)
}

// @Internal
// Check whether a bs58 equals to treasury bs58.
// func equalBs58ToTreasury(bs58 string) bool {
// 	return treasuryPubKey().Equals(solana.MustPublicKeyFromBase58(bs58))
// }

// ===== Cluster =====
// Solana Cluster.
var _cluster Cluster

// @Internal
// Setter for `_cluster`
func setCluster(cluster Cluster) {
	_cluster = cluster
}

// @Internal
// Getter for `_cluster`
func cluster() Cluster {
	return _cluster
}

// ===== Rpc Url =====
// Solana rpc url.
var _rpcUrl string

// @Internal
// Setter for `_rpcUrl`
func setRpcUrl(rpcUrl string) {
	if rpcUrl == "" {
		rpcUrl = getDefaultRpcUrl(cluster())
	}

	_rpcUrl = rpcUrl
}

// @Internal
// Get solana rpc.
func rpcUrl() string {
	return _rpcUrl
}

// @Internal
// Get default rpc url for given cluster.
// If am empty string is provided as cluster, it refers to `DEFAULT_CLUSTER` constant.
func getDefaultRpcUrl(clusters ...Cluster) string {
	var cluster Cluster

	if len(clusters) == 0 || clusters[0] == "" {
		cluster = DEFAULT_CLUSTER
	} else {
		cluster = clusters[0]
	}

	clusterToRpcUrl := make(map[Cluster]string, 3)
	clusterToRpcUrl[ClusterTestNet] = rpc.TestNet_RPC
	clusterToRpcUrl[ClusterDevNet] = rpc.DevNet_RPC
	clusterToRpcUrl[ClusterMainNetBeta] = rpc.MainNetBeta_RPC

	return clusterToRpcUrl[cluster]
}

// ===== Client =====
// @Internal
// Create a client object and return pointer.
func newClient() *rpc.Client {
	client := rpc.New(rpcUrl())

	return client
}

// @Internal
// Supported spl tokens.
var supportedSplTokens = []SplTokenMeta{}

// @Internal
// Initialize supported spl tokens.
func initializeSplTokens() {
	supportedSplTokens = []SplTokenMeta{
		{
			Type:                 SolSplToken,
			MintAddress:          solana.MustPublicKeyFromBase58(config.SOL_SPL_ADDRESS),
			TreasuryTokenAccount: *treasuryPubKey(),
			Decimals:             9,
			Keyword:              "SOL",
			Image:                "https://duelana-bucket-prod.s3.us-east-2.amazonaws.com/coins/SOL.png",
		},
		{
			Type:                 BonkSplToken,
			MintAddress:          solana.MustPublicKeyFromBase58(config.BONK_SPL_ADDRESS),
			TreasuryTokenAccount: treasuryTokenAccount(solana.MustPublicKeyFromBase58(config.BONK_SPL_ADDRESS)),
			Decimals:             5,
			Keyword:              "Bonk",
			Image:                "https://duelana-bucket-prod.s3.us-east-2.amazonaws.com/coins/Bonk.png",
		},
		{
			Type:                 UsdcSplToken,
			MintAddress:          solana.MustPublicKeyFromBase58(config.USDC_SPL_ADDRESS),
			TreasuryTokenAccount: treasuryTokenAccount(solana.MustPublicKeyFromBase58(config.USDC_SPL_ADDRESS)),
			Decimals:             6,
			Keyword:              "USDC",
			Image:                "https://duelana-bucket-prod.s3.us-east-2.amazonaws.com/coins/USDC.png",
		},
		{
			Type:                 BokuSplToken,
			MintAddress:          solana.MustPublicKeyFromBase58(config.BOKU_SPL_ADDRESS),
			TreasuryTokenAccount: treasuryTokenAccount(solana.MustPublicKeyFromBase58(config.BOKU_SPL_ADDRESS)),
			Decimals:             9,
			Keyword:              "BOKU",
			Image:                "https://duelana-bucket-prod.s3.us-east-2.amazonaws.com/coins/BOKU.png",
		},
	}
}

// @External
// Returns whether spl token meta is for sol or spl.
func isSolSplMeta(meta SplTokenMeta) bool {
	if meta.Keyword == "SOL" &&
		meta.MintAddress == solana.MustPublicKeyFromBase58(config.SOL_SPL_ADDRESS) {
		return true
	}
	return false
}

// @External
// Retrieve supported spl tokens.
func supportedSpls() []SplTokenMeta {
	return supportedSplTokens
}

// @External
// Initializer.
func initialize(param *InitParam) error {
	if param == nil {
		return fmt.Errorf("solana-aggregator: initializer param is nil pointer")
	}

	setTreasuryBs58(param.TreasuryBs58)
	setCluster(param.Cluster)
	setRpcUrl(param.RpcUrl)
	initializeTxDecode()
	initializeSplTokens()

	return nil
}

// @Internal
// Makes an error object with given topic and reasons.
func makeError(category string, reason string, err error) error {
	return utils.MakeError("solana-aggregator", category, reason, err)
}
