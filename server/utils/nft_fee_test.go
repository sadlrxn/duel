package utils

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/davecgh/go-spew/spew"
)

func TestDetermineNftsForFee(t *testing.T) {
	{
		nfts := []types.NftDetails{
			{
				Name:        "Shark",
				Price:       21252309,
				MintAddress: "A",
			},
			{
				Name:        "Donkey",
				Price:       418250,
				MintAddress: "B",
			},
			{
				Name:        "Cat",
				Price:       2398888,
				MintAddress: "C",
			},
		}
		fee := int64(1062615)
		DetermineNFTs4Fee(nfts, fee)
	}

	{
		nfts := []types.NftDetails{
			{
				Name:        "A",
				Price:       8,
				MintAddress: "A",
			},
			{
				Name:        "B",
				Price:       10,
				MintAddress: "B",
			},
			{
				Name:        "C",
				Price:       10,
				MintAddress: "C",
			},
			{
				Name:        "D",
				Price:       23,
				MintAddress: "D",
			},
		}
		fee := int64(31)
		DetermineNFTs4Fee(nfts, fee)
	}

	{
		nfts := []types.NftDetails{
			{
				Name:        "A",
				Price:       8,
				MintAddress: "A",
			},
			{
				Name:        "B",
				Price:       10,
				MintAddress: "B",
			},
			{
				Name:        "C",
				Price:       10,
				MintAddress: "C",
			},
			{
				Name:        "D",
				Price:       10,
				MintAddress: "D",
			},
			{
				Name:        "E",
				Price:       20,
				MintAddress: "E",
			},
		}
		fee := int64(30)
		solution, profits := DetermineNFTs4Fee(nfts, fee)
		spew.Dump(solution)
		spew.Dump(profits)
	}

	t.Fatalf("")
}
