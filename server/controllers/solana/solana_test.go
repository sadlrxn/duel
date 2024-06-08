package solana

import (
	"context"
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TestSolanaPublicKeyConvert(t *testing.T) {

	pubKey := solana.MustPublicKeyFromBase58(config.SOL_SPL_ADDRESS)
	t.Fatalf("pub key string: %s", pubKey.String())

}

func TestGetTransaction(t *testing.T) {
	setRpcUrl("https://solana-mainnet.g.alchemy.com/v2/ZWXSwpVnO3xRBeBuTbWnqW4b6kt27lcA")
	client := newClient()
	for i := 0; i < 50; i++ {
		maxSupportedTxVersion := uint64(0)
		_, err := client.GetTransaction(
			context.TODO(),
			solana.MustSignatureFromBase58("4T9y9y9hMQ37UphKA1vDvXywT2c2kQukQWHgFs6AjUpuN7FMxUMXF2KBMeqiMh3Sp3Ayy7j9ndY2GTrHAn3tohTk"),
			&rpc.GetTransactionOpts{
				Encoding:                       solana.EncodingBase64,
				MaxSupportedTransactionVersion: &maxSupportedTxVersion,
			},
		)

		if err != nil {
			t.Fatalf(
				"failed to get transaction: index: %d, err: %v",
				i,
				err,
			)
		}
	}
}
