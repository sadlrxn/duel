package utils

import (
	"sort"

	"github.com/Duelana-Team/duelana-v1/types"
)

type NFTs []types.NftDetails

func (nfts NFTs) Len() int {
	return len(nfts)
}
func (nfts NFTs) Swap(i, j int) {
	nfts[i], nfts[j] = nfts[j], nfts[i]
}
func (nfts NFTs) Less(i, j int) bool {
	return nfts[i].Price < nfts[j].Price || (nfts[i].Price == nfts[j].Price && nfts[i].MintAddress < nfts[j].MintAddress)
}

func betterSolution(solution, candidate NFTs) NFTs {
	if len(solution) > len(candidate) {
		return solution
	} else if len(solution) < len(candidate) {
		return candidate
	} else {
		if solution[len(solution)-1].Price <= candidate[len(candidate)-1].Price {
			return solution
		} else {
			return candidate
		}
	}
}

var solution NFTs
var max int64

func collect(fullList NFTs, subList NFTs, cursor int, limit int64) {
	var sum int64
	for i := 0; i < len(subList); i++ {
		sum += subList[i].Price
	}
	if sum > limit {
		return
	} else if sum > max && sum <= limit {
		solution = NFTs{}
		max = sum
		collect(fullList, NFTs{}, -1, limit)
	} else if sum == max {
		solution = betterSolution(solution, subList)
		if cursor >= len(fullList)-1 {
			return
		}
		collect(fullList, append(subList, fullList[cursor+1]), cursor+1, limit)
	} else if sum < max {
		for i := cursor + 1; i < len(fullList); i++ {
			var list NFTs
			list = append(list, subList...)
			list = append(list, fullList[i])
			collect(fullList, list, i, limit)
		}
	}
}

func DetermineNFTs4Fee(nfts []types.NftDetails, fee int64) ([]types.NftDetails, []types.NftDetails) {
	sort.Sort(NFTs(nfts))

	var nftCandidates NFTs
	for i := 0; i < len(nfts); i++ {
		if nfts[i].Price > fee {
			break
		}
		nftCandidates = append(nftCandidates, nfts[i])
	}

	solution = NFTs{}
	if len(nftCandidates) == 0 {
		return []types.NftDetails{}, nfts
	}
	max = nftCandidates[len(nftCandidates)-1].Price
	collect(nftCandidates, NFTs{}, -1, fee)
	profits := []types.NftDetails{}
	isSolution := make(map[string]bool)
	for i := 0; i < len(solution); i++ {
		isSolution[solution[i].MintAddress] = true
	}
	count := len(nfts)
	for i := 0; i < count; i++ {
		if !isSolution[nfts[i].MintAddress] {
			profits = append(profits, nfts[i])
		}
	}
	return solution, profits
}
