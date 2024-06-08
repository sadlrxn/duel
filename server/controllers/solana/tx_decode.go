package solana

import (
	"context"
	"errors"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	solana_token_account "github.com/gagliardetto/solana-go/programs/associated-token-account"
	solana_system "github.com/gagliardetto/solana-go/programs/system"
	solana_token "github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

// @Internal
// Initialize transaction decoders.
func initializeTxDecode() {
	solana_system.SetProgramID(solana.SystemProgramID)
	solana_token.SetProgramID(solana.TokenProgramID)
	solana_token_account.SetProgramID(solana.SPLAssociatedTokenAccountProgramID)
}

// @Internal
// Get solana transaction from tx hash.
func getTransaction(txHashBs58 string) (*rpc.GetTransactionResult, error) {
	txHash, err := solana.SignatureFromBase58(txHashBs58)
	if err != nil {
		return nil, makeError("getTransaction", "failed to get tx hash", err)
	}

	client := newClient()
	maxSupportedTxVersion := uint64(0)
	txOut, err := client.GetTransaction(
		context.TODO(),
		txHash,
		&rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			MaxSupportedTransactionVersion: &maxSupportedTxVersion,
		},
	)
	if err != nil {
		return nil, makeError("getTransaction", "failed to get transaction result", err)
	}

	return txOut, nil
}

// @Internal
// Get decode transaction from GetTransactionResult.
func getDecodedTransaction(txOut *rpc.GetTransactionResult) (*solana.Transaction, error) {
	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(txOut.Transaction.GetBinary()))
	if err != nil {
		return nil, makeError("getDecodedTransaction", "failed to decode transaction binary", err)
	}

	return tx, nil
}

// @Internal
// Decode a instruction.
func decodeTransactionInstructions(tx *solana.Transaction) (*[]DecodedInstruction, error) {
	if tx == nil {
		return nil, makeError("decodeTransactionInstructions", "Invalid parameter", errors.New("parameter (tx) is nil pointer"))
	}
	result := make([]DecodedInstruction, len(tx.Message.Instructions))
	for i, instruction := range tx.Message.Instructions {
		progKey, err := tx.ResolveProgramIDIndex(instruction.ProgramIDIndex)
		if err != nil {
			result[i] = DecodedInstruction{
				InstructionType: InstructionUnknown,
				Meta:            nil,
			}
			continue
		}

		accounts, err := instruction.ResolveInstructionAccounts(&tx.Message)
		if err != nil {
			result[i] = DecodedInstruction{
				InstructionType: InstructionUnknown,
				Meta:            nil,
			}
			continue
		}
		decodedInstruction, err := solana.DecodeInstruction(
			progKey,
			accounts,
			instruction.Data,
		)
		if err != nil {
			result[i] = DecodedInstruction{
				InstructionType: InstructionUnknown,
				Meta:            nil,
			}
			continue
		}

		if tokenInstrucion, ok := decodedInstruction.(*solana_token.Instruction); ok {
			if instructionMeta, ok := tokenInstrucion.Impl.(*solana_token.Transfer); ok {
				result[i] = DecodedInstruction{
					InstructionType: InstructionTokenTransfer,
					Meta:            instructionMeta,
				}
				continue
			}
			if instructionMeta, ok := tokenInstrucion.Impl.(*solana_token.TransferChecked); ok {
				result[i] = DecodedInstruction{
					InstructionType: InstructionTokenTransferChecked,
					Meta:            instructionMeta,
				}
				continue
			}
		} else if tokenAccountInstruction, ok := decodedInstruction.(*solana_token_account.Instruction); ok {
			if instructionMeta, ok := tokenAccountInstruction.Impl.(*solana_token_account.Create); ok {
				result[i] = DecodedInstruction{
					InstructionType: InstructionTokenCreateAccount,
					Meta:            instructionMeta,
				}
				continue
			}
		} else if systemInstruction, ok := decodedInstruction.(*solana_system.Instruction); ok {
			if instructionMeta, ok := systemInstruction.Impl.(*solana_system.Transfer); ok {
				result[i] = DecodedInstruction{
					InstructionType: InstructionSolTransfer,
					Meta:            instructionMeta,
				}
				continue
			}
		}

		result[i] = DecodedInstruction{
			InstructionType: InstructionUnknown,
			Meta:            nil,
		}
	}

	return &result, nil
}

// @Internal
// Analyse nft transfer from transaction result.
func analyseNftTransfer(txOut *rpc.GetTransactionResult) (*AnalyseNftTransferResult, error) {
	postBalances := &txOut.Meta.PostTokenBalances
	nfts := make([]string, len(*postBalances)/2)
	nftBucket := make(map[solana.PublicKey]bool, len(nfts))
	var to *solana.PublicKey = nil
	var from *solana.PublicKey = nil
	j := 0

	for _, balance := range *postBalances {
		if balance.UiTokenAmount.Decimals != 0 ||
			(balance.UiTokenAmount.UiAmountString != "1" &&
				balance.UiTokenAmount.UiAmountString != "0") {
			continue
		}
		if _, prs := nftBucket[balance.Mint]; !prs {
			nftBucket[balance.Mint] = true
			nfts[j] = balance.Mint.String()
			j = j + 1
		}

		if balance.UiTokenAmount.UiAmountString == "1" {
			if to == nil {
				to = balance.Owner
			} else if !to.Equals(*balance.Owner) {
				return nil, makeError("analyseNftTransfer", "failed to analyse pre balances", errors.New("mismatching senders"))
			}
		} else if balance.UiTokenAmount.UiAmountString == "0" {
			if from == nil {
				from = balance.Owner
			} else if !from.Equals(*balance.Owner) {
				return nil, makeError("analyseNftTransfer", "failed to analyse pre balances", errors.New("mismatching receivers"))
			}
		}
	}

	return &AnalyseNftTransferResult{
		From: *from,
		To:   *to,
		Nfts: nfts,
	}, nil
}

// @Internal
// Check whether a transaction is succeed one or not.
func isSucceedTransaction(txOut *rpc.GetTransactionResult) bool {
	return txOut.Meta.Err == nil
}

// @Internal
// Analyse instruction array and returns whether it is sol transfer or not.
func isSolTransferInstArray(instArray *[]DecodedInstruction) bool {
	return instArray != nil && len(*instArray) == 1 && (*instArray)[0].InstructionType == InstructionSolTransfer
}

// @Internal
// Analyse instruction array and returns whether it is nft transfer or not.
func isNftTransferInstArray(instArray *[]DecodedInstruction) bool {
	if instArray == nil {
		return false
	}
	for i, inst := range *instArray {
		if i != 0 &&
			inst.InstructionType != InstructionTokenCreateAccount &&
			inst.InstructionType != InstructionTokenTransfer &&
			inst.InstructionType != InstructionTokenTransferChecked {
			return false
		}
	}

	return true
}

// @Internal
// Analyse instruction array and returns whether it is bonk transfer or not.
func isSplTransferInstArray(instArray *[]DecodedInstruction) (*DecodedTransaction, bool) {
	if len(*instArray) > 3 {
		return nil, false
	}
	if (len(*instArray) == 2 || len(*instArray) == 3) &&
		((*instArray)[len(*instArray)-2].InstructionType != InstructionTokenCreateAccount ||
			((*instArray)[len(*instArray)-1].InstructionType != InstructionTokenTransfer &&
				(*instArray)[len(*instArray)-1].InstructionType != InstructionTokenTransferChecked)) {
		return nil, false
	}
	if len(*instArray) == 1 &&
		(*instArray)[0].InstructionType != InstructionTokenTransfer &&
		(*instArray)[0].InstructionType != InstructionTokenTransferChecked {
		return nil, false
	}

	if transferTx, ok := (*instArray)[len(*instArray)-1].Meta.(*solana_token.Transfer); ok {
		for i, supportedSpl := range supportedSplTokens {
			if i == 0 { // skip sol
				continue
			}

			if transferTx.Accounts[0].PublicKey == supportedSpl.TreasuryTokenAccount {
				return &DecodedTransaction{
					TransactionType: TransactionSplWithdraw,
					Lamports:        *transferTx.Amount,
					SplToken:        &supportedSpl,
					Participant:     transferTx.Accounts[2].PublicKey.String(),
				}, true
			}
			if transferTx.Accounts[1].PublicKey == supportedSpl.TreasuryTokenAccount {
				return &DecodedTransaction{
					TransactionType: TransactionSplDeposit,
					Lamports:        *transferTx.Amount,
					SplToken:        &supportedSpl,
					Participant:     transferTx.Accounts[2].PublicKey.String(),
				}, true
			}
		}
	}

	if transferCheckedTx, ok := (*instArray)[len(*instArray)-1].Meta.(*solana_token.TransferChecked); ok {
		for i, supportedSpl := range supportedSplTokens {
			if i == 0 { // skip sol
				continue
			}

			if transferCheckedTx.Accounts[0].PublicKey == supportedSpl.TreasuryTokenAccount {
				return &DecodedTransaction{
					TransactionType: TransactionSplWithdraw,
					Lamports:        *transferCheckedTx.Amount,
					SplToken:        &supportedSpl,
					Participant:     transferCheckedTx.Accounts[3].PublicKey.String(),
				}, true
			}
			if transferCheckedTx.Accounts[2].PublicKey == supportedSpl.TreasuryTokenAccount {
				return &DecodedTransaction{
					TransactionType: TransactionSplDeposit,
					Lamports:        *transferCheckedTx.Amount,
					SplToken:        &supportedSpl,
					Participant:     transferCheckedTx.Accounts[3].PublicKey.String(),
				}, true
			}
		}
	}

	return nil, false
}

// @External
// Analyse transaction type.
func decodeTransactionType(txHashBs58 string) (*DecodedTransaction, error) {
	txOut, err := getTransaction(txHashBs58)
	nothing := DecodedTransaction{
		TransactionType: TransactionNothing,
	}
	if err != nil {
		return &nothing, makeError("decodeTransactionType", "failed to get transaction", err)
	}

	if !isSucceedTransaction(txOut) {
		return &DecodedTransaction{
			TransactionType: TransactionNothing,
			Failed:          true,
		}, makeError("decodeTransactionType", "failed transaction", errors.New("provided transaction is a failed one"))
	}

	tx, err := getDecodedTransaction(txOut)
	if err != nil {
		return &nothing, makeError("decodeTransactionType", "failed to decode transaction", err)
	}

	decodedInsts, err := decodeTransactionInstructions(tx)
	if err != nil {
		return &nothing, makeError("decodeTransactionType", "failed to decode instructions", err)
	}

	if isSolTransferInstArray(decodedInsts) {
		solTransfer, ok := (*decodedInsts)[0].Meta.(*solana_system.Transfer)
		if !ok {
			return &nothing, makeError("decodeTransactionType", "failed to get system transfer meta", errors.New("failed to parse assigned meta"))
		}

		var from *solana.PublicKey = solTransfer.AccountMetaSlice[0].PublicKey.ToPointer()
		var to *solana.PublicKey = solTransfer.AccountMetaSlice[1].PublicKey.ToPointer()
		if !solTransfer.AccountMetaSlice[0].IsSigner {
			from, to = to, from
		}

		if equalPubKeyToTreasury(*from) {
			return &DecodedTransaction{
				TransactionType: TransactionSplWithdraw,
				Lamports:        *solTransfer.Lamports,
				Participant:     to.String(),
				SplToken:        &supportedSplTokens[0],
			}, nil
		}

		if equalPubKeyToTreasury(*to) {
			return &DecodedTransaction{
				TransactionType: TransactionSplDeposit,
				Lamports:        *solTransfer.Lamports,
				Participant:     from.String(),
				SplToken:        &supportedSplTokens[0],
			}, nil
		}

		return &nothing, nil
	}

	fmt.Printf("decode insts: %v", decodedInsts)

	if decodedTx, ok := isSplTransferInstArray(decodedInsts); ok {
		return decodedTx, nil
	}

	if isNftTransferInstArray(decodedInsts) {
		nftTransfer, err := analyseNftTransfer(txOut)
		if err != nil {
			return &nothing, makeError("decodeTransactionType", "failed to analyse nft transfer", err)
		}

		if equalPubKeyToTreasury(nftTransfer.From) {
			return &DecodedTransaction{
				TransactionType: TransactionNftWithdraw,
				Nfts:            nftTransfer.Nfts,
				Participant:     nftTransfer.To.String(),
			}, nil
		}

		if equalPubKeyToTreasury(nftTransfer.To) {
			return &DecodedTransaction{
				TransactionType: TransactionNftDeposit,
				Nfts:            nftTransfer.Nfts,
				Participant:     nftTransfer.From.String(),
			}, nil
		}

		return &nothing, nil
	}

	return &nothing, nil
}
