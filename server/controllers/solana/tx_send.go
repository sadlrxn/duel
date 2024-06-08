package solana

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/gagliardetto/solana-go"
	solana_token_account "github.com/gagliardetto/solana-go/programs/associated-token-account"
	solana_system "github.com/gagliardetto/solana-go/programs/system"
	solana_token "github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

// Get recent blockhash retry counts.
const GET_RECENT_BLOCKHASH_RETRY = 5

// Get recent blockhash retry wait time in seconds.
const GET_RECENT_BLOCKHASH_WAIT = 2

// Get balance retry counts.
const GET_BALANCE_RETRY = 5

// Get balance retry wait time in seconds.
const GET_BALANCE_WAIT = 5

// Send transaction retry counts.
const SEND_TRANSACTION_RETRY = 5

// Send transaction wait time in seconds.
const SEND_TRANSACTION_WAIT = 5

// Keywords to trigger send transaction retry.
var SEND_TRANSACTION_CONDITIONS = []string{
	"BlockhashNotFound",
	"failed to get recent blockhash",
}

// @Internal
// Returns solana recent blockhash. Try once.
func getRecentBlockHashOnce() (*solana.Hash, error) {
	client := newClient()
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, makeError("getRecentBlockHashOnce", "failing to get recent block hash...", err)
	}

	return &recent.Value.Blockhash, nil
}

// @Internal
// Returns solana recent blockhash. Retires.
func getRecentBlockHashRetry() (*solana.Hash, error) {
	var finalError error = errors.New("")
	for i := 0; i < GET_RECENT_BLOCKHASH_RETRY; i = i + 1 {
		recentBlockHash, err := getRecentBlockHashOnce()
		if err == nil {
			return recentBlockHash, nil
		}
		finalError = fmt.Errorf("%v\n\r%v", finalError.Error(), err.Error())
		if i < GET_RECENT_BLOCKHASH_RETRY-1 {
			time.Sleep(GET_RECENT_BLOCKHASH_WAIT * time.Second)
		}
	}

	return nil, makeError("getRecentBlockHashRetry", "failing to get recent block hash over several retries", finalError)
}

// @Internal
// Builds solana transaction to send lamports to.
func buildSendLamportsTx(to solana.PublicKey, lamports uint64) (*solana.Transaction, error) {
	recentBlockHash, err := getRecentBlockHashRetry()
	if err != nil {
		return nil, makeError("buildSendLamportsTx", "failed to get recent blockhash", err)
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			solana_system.NewTransferInstruction(
				lamports,
				*treasuryPubKey(),
				to,
			).Build(),
		},
		*recentBlockHash,
		solana.TransactionPayer(*treasuryPubKey()),
	)
	if err != nil {
		return nil, makeError("buildSendLamportsTx", "failed to build transaction object", err)
	}

	return tx, nil
}

// @Internal
// Get associated token address of treasury.
func treasuryTokenAccount(mint solana.PublicKey) solana.PublicKey {
	return tokenAccount(*treasuryPubKey(), mint)
}

// @Internal
// Get associated token address of any pubkey.
func tokenAccount(owner solana.PublicKey, mint solana.PublicKey) solana.PublicKey {
	associatedTokenAddress, _, _ := solana.FindAssociatedTokenAddress(
		owner,
		mint,
	)

	return associatedTokenAddress
}

// @Internal
// Builds solana transaction to send spl tokens to.
func buildSendSplTokensTx(to solana.PublicKey, mint solana.PublicKey, amount uint64) (*solana.Transaction, error) {
	recentBlockHash, err := getRecentBlockHashRetry()
	if err != nil {
		return nil, makeError("buildSendSplTokensTx", "failed to get recent blockhash", err)
	}

	insts := []solana.Instruction{}
	destTokenAccount := tokenAccount(to, mint)
	balance, err := getBalanceRetry(destTokenAccount)
	if err != nil {
		return nil, makeError("buildSendSplTokensTx", "failed to get balance of dest token account", err)
	}
	if balance == 0 {
		insts = append(
			insts,
			solana_token_account.NewCreateInstruction(
				*treasuryPubKey(),
				to,
				mint,
			).Build(),
		)
	}
	insts = append(
		insts,
		solana_token.NewTransferInstruction(
			amount,
			treasuryTokenAccount(mint),
			destTokenAccount,
			*treasuryPubKey(),
			[]solana.PublicKey{},
		).Build(),
	)

	tx, err := solana.NewTransaction(
		insts,
		*recentBlockHash,
		solana.TransactionPayer(*treasuryPubKey()),
	)
	if err != nil {
		return nil, makeError("buildSendSplTokensTx", "failed to build transaction object", err)
	}

	return tx, nil
}

// @Internal
// Sign and send transaction.
func signAndSendTx(tx *solana.Transaction) (*solana.Signature, error) {
	if tx == nil {
		return nil, makeError("signAndSendTx", "invalid parameter", errors.New("parameter (tx): null pointer"))
	}
	if _, err := tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if equalPubKeyToTreasury(key) {
				return treasuryKeyPair()
			}
			return nil
		},
	); err != nil {
		return nil, makeError("signAndSendTx", "failed to sign transaction", err)
	}

	client := newClient()
	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		return nil, makeError("signAndSendTx", "failed to send transaction", err)
	}

	return &sig, nil
}

// @Internal
// Convert solana signature type to base58 string.
func convertSigToBs58(sig solana.Signature) string {
	hex := make([]byte, 64)
	for i, b := range sig {
		hex[i] = byte(b)
	}
	return base58.Encode(hex)
}

// @Internal
// Check whether the error is for retry.
func retry(err error) bool {
	if err == nil {
		return false
	}
	for _, condition := range SEND_TRANSACTION_CONDITIONS {
		if strings.Contains(err.Error(), condition) {
			return true
		}
	}
	return false
}

// @Internal
// Send lamports. Try Once.
func sendLamportsOnce(param *SendLamportsRequest) (string, error) {
	nothing := ""
	if param == nil {
		return "", makeError("sendLamports", "invalid parameter", errors.New("parameter (param): null pointer"))
	}

	toPubKey, err := solana.PublicKeyFromBase58(param.To)
	if err != nil {
		return nothing, makeError("sendLamports", "failed to get public key from base58", err)
	}

	tx, err := buildSendLamportsTx(toPubKey, param.Lamports)
	if err != nil {
		return nothing, makeError("sendLamports", "failed to build send lamports tx", err)
	}

	sig, err := signAndSendTx(tx)
	if err != nil {
		return nothing, makeError("sendLamports", "failed to sign and send tx", err)
	}

	return convertSigToBs58(*sig), nil
}

// @External
// Send lamports. Retries.
func sendLamports(param *SendLamportsRequest) (string, error) {
	var finalError error = errors.New("")
	for i := 0; i < SEND_TRANSACTION_RETRY; i = i + 1 {
		txHash, err := sendLamportsOnce(param)
		if !retry(err) {
			return txHash, err
		}
		finalError = fmt.Errorf("%v\n\r%v", finalError.Error(), err.Error())
		time.Sleep(SEND_TRANSACTION_WAIT * time.Second)
	}
	return "", makeError("sendLamports", "failing to send lamports over several retries", finalError)
}

// @Internal
// Send spl tokens. Try once.
func sendSplTokensOnce(param *SendSplTokenRequest) (string, error) {
	if param.Mint == config.SOL_SPL_ADDRESS {
		return sendLamportsOnce(&SendLamportsRequest{
			To:       param.To,
			Lamports: param.Amount,
		})
	}
	nothing := ""
	if param == nil {
		return "", makeError("sendSplTokensOnce", "invalid parameter", errors.New("parameter (param): null pointer"))
	}

	toPubKey, err := solana.PublicKeyFromBase58(param.To)
	if err != nil {
		return nothing, makeError("sendSplTokensOnce", "failed to get public key from base58 of to address", err)
	}

	mintPubKey, err := solana.PublicKeyFromBase58(param.Mint)
	if err != nil {
		return nothing, makeError("sendSplTokensOnce", "failed to get public key from base58 of mint address", err)
	}

	tx, err := buildSendSplTokensTx(toPubKey, mintPubKey, param.Amount)
	if err != nil {
		return nothing, makeError("sendSplTokensOnce", "failed to build send lamports tx", err)
	}

	sig, err := signAndSendTx(tx)
	if err != nil {
		return nothing, makeError("sendSplTokensOnce", "failed to sign and send tx", err)
	}

	return convertSigToBs58(*sig), nil
}

// @External
// Send spl tokens. Retries.
func sendSplTokens(param *SendSplTokenRequest) (string, error) {
	var finalError error = errors.New("")
	for i := 0; i < SEND_TRANSACTION_RETRY; i = i + 1 {
		txHash, err := sendSplTokensOnce(param)
		if !retry(err) {
			return txHash, err
		}
		finalError = fmt.Errorf("%v\n\r%v", finalError.Error(), err.Error())
		time.Sleep(SEND_TRANSACTION_WAIT * time.Second)
	}
	return "", makeError("sendSplTokens", "failing to send spl tokens over several retries", finalError)
}

// @Internal
// Get balance of a pubkey. Retries.
func getBalanceRetry(account solana.PublicKey) (uint64, error) {
	var finalError error = errors.New("")
	client := newClient()
	for i := 0; i < GET_BALANCE_RETRY; i = i + 1 {
		balance, err := client.GetBalance(context.TODO(), account, rpc.CommitmentFinalized)
		if err == nil {
			return balance.Value, nil
		}
		finalError = fmt.Errorf("%v\n\r%v", finalError.Error(), err.Error())
		if i < GET_BALANCE_RETRY-1 {
			time.Sleep(GET_BALANCE_WAIT * time.Second)
		}
	}

	return 0, makeError("getBalanceRetry", "failing to get recent block hash over several retries", finalError)
}

// @Internal
// Gets instructions to transfer mint to account.
func getInstsForTransfer(to solana.PublicKey, mint solana.PublicKey) (*[]solana.Instruction, error) {
	instructions := []solana.Instruction{}
	toTokenAccount, _, err := solana.FindAssociatedTokenAddress(to, mint)
	if err != nil {
		return nil, makeError("getInstsForTransfer", "failed to get associated token address of dest", err)
	}
	fromTokenAccount, _, err := solana.FindAssociatedTokenAddress(*treasuryPubKey(), mint)
	if err != nil {
		return nil, makeError("getInstsForTransfer", "failed to get associated token address of source", err)
	}

	lamports, err := getBalanceRetry(toTokenAccount)
	if err != nil {
		return nil, makeError("getInstsForTransfer", "failed to get balance of the token accoutn", err)
	}

	if lamports == 0 {
		instructions = append(instructions, solana_token_account.Create{
			Payer:  *treasuryPubKey(),
			Wallet: to,
			Mint:   mint,
		}.Build())
	}

	instructions = append(instructions, solana_token.NewTransferInstruction(
		1,
		fromTokenAccount,
		toTokenAccount,
		*treasuryPubKey(),
		[]solana.PublicKey{},
	).Build())

	return &instructions, nil
}

// @Internal
// Builds send nft transaction.
func buildSendNftsTx(to solana.PublicKey, mints *[]solana.PublicKey) (*solana.Transaction, error) {
	instructions := []solana.Instruction{}
	for _, mint := range *mints {
		insts, err := getInstsForTransfer(to, mint)
		if err != nil {
			return nil, makeError("buildSendNftsTx", fmt.Sprintf("failed to get insts for mint(%v) to account(%v)", mint.String(), to.String()), err)
		}

		instructions = append(instructions, *insts...)
	}

	recentBlockHash, err := getRecentBlockHashRetry()
	if err != nil {
		return nil, makeError("buildSendNftsTx", "failed to get recent blockhash", err)
	}

	tx, err := solana.NewTransaction(
		instructions,
		*recentBlockHash,
		solana.TransactionPayer(*treasuryPubKey()),
	)
	if err != nil {
		return nil, makeError("buildSendNftsTx", "failed to build transaction object", err)
	}

	return tx, nil
}

// @Internal
// Convert string array to pubkey array.
func convertBs58sToPubKeys(bs58s *[]string) (*[]solana.PublicKey, error) {
	if bs58s == nil {
		return nil, makeError("convertBs58sToPubKeys", "invalid parameter", errors.New("parameter (bs58): null pointer"))
	}
	pubKeys := make([]solana.PublicKey, len(*bs58s))

	for i, bs58 := range *bs58s {
		pubKey, err := solana.PublicKeyFromBase58(bs58)
		if err != nil {
			return nil, makeError("convertBs58sToPubKeys", fmt.Sprintf("failed to get pubkey from base58(%v)", bs58), err)
		}

		pubKeys[i] = pubKey
	}

	return &pubKeys, nil
}

// @Internal
// Send nfts. Try Once.
func sendNftsOnce(param *SendNftsRequest) (string, error) {
	nothing := ""
	if param == nil {
		return "", makeError("sendNfts", "invalid parameter", errors.New("parameter (param): null pointer"))
	}

	toPubKey, err := solana.PublicKeyFromBase58(param.To)
	if err != nil {
		return nothing, makeError("sendNfts", "failed to get public key from base58", err)
	}

	mints, err := convertBs58sToPubKeys(&param.Nfts)
	if err != nil {
		return nothing, makeError("sendNfts", "failed to convert bs58 mints to pubkeys", err)
	}

	tx, err := buildSendNftsTx(toPubKey, mints)
	if err != nil {
		return nothing, makeError("sendNfts", "failed to build send lamports tx", err)
	}

	sig, err := signAndSendTx(tx)
	if err != nil {
		return nothing, makeError("sendNfts", "failed to sign and send tx", err)
	}

	return convertSigToBs58(*sig), nil
}

// @External
// Send nfts. Retries
func sendNfts(param *SendNftsRequest) (string, error) {
	var finalError error = errors.New("")
	for i := 0; i < SEND_TRANSACTION_RETRY; i = i + 1 {
		txHash, err := sendNftsOnce(param)
		if !retry(err) {
			return txHash, err
		}
		finalError = fmt.Errorf("%v\n\r%v", finalError.Error(), err.Error())
		time.Sleep(SEND_TRANSACTION_WAIT * time.Second)
	}
	return "", makeError("sendNfts", "failing to send nfts over several retries", finalError)
}
