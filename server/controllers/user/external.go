package user

import "github.com/Duelana-Team/duelana-v1/types"

func GetNftDetailsFromMintAddresses(
	mintAddresses []string,
) (int64, []types.NftDetails) {
	return getNftDetailsFromMintAddresses(mintAddresses)
}
