package jackpot

import (
	"fmt"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/types"
	"golang.org/x/sync/syncmap"
)

func TestRoundEnd(t *testing.T) {
	var c Controller
	// piggy_amount, piggy := utils.GetNftDetailsFromMintAddresses([]string{"BA2rHhvMH1Xt23cB8Wsxj6KpToKHWRE2VKnrhdhFRAr6"})
	// chicken_amount, chicken := utils.GetNftDetailsFromMintAddresses([]string{"2VZKHR6iQAhtyb1yhHhu4srbH86K5EpZaSkB3pRo8Hzo"})
	piggy_amount := int64(663)
	piggy := types.NftDetails{
		MintAddress: "BA2rHhvMH1Xt23cB8Wsxj6KpToKHWRE2VKnrhdhFRAr6",
		Price:       piggy_amount,
	}
	chicken_amount := int64(590)
	chicken := types.NftDetails{
		MintAddress: "2VZKHR6iQAhtyb1yhHhu4srbH86K5EpZaSkB3pRo8Hzo",
		Price:       chicken_amount,
	}

	winner := uint(2)
	c.status = Started
	c.fee = 5
	c.winnerID = &winner
	c.playerToBets = syncmap.Map{}
	c.playerToBets.Store(uint(1), Bet4Player{
		TotalUsdAmount: 0,
		TotalNftAmount: piggy_amount,
		Bets: []BetData{
			{
				Amount:    0,
				NftAmount: piggy_amount,
				Nfts:      []types.NftDetails{piggy},
			},
		},
	})
	c.playerToBets.Store(uint(2), Bet4Player{
		TotalUsdAmount: 0,
		TotalNftAmount: chicken_amount,
		Bets: []BetData{
			{
				Amount:    0,
				NftAmount: chicken_amount,
				Nfts:      []types.NftDetails{chicken},
			},
		},
	})

	usdProfit, nfts4Profit, totalProfit, usdFee, nfts4Fee, totalFee, totalAmount := c.calculateJackpots()

	fmt.Printf("usdProfit: %d", usdProfit)
	fmt.Printf("NftProfit: %v", nfts4Profit)
	fmt.Printf("totalProfit: %d", totalProfit)
	fmt.Printf("usdFee: %d", usdFee)
	fmt.Printf("nftfee: %v", nfts4Fee)
	fmt.Printf("totalFee: %d", totalFee)
	t.Fatalf("totalAmount: %d", totalAmount)
}

func TestTimer(t *testing.T) {
	fmt.Println(time.Now().String())
	var timer *time.Timer
	timer = time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		fmt.Println(time.Now().String())
	}()
	timer = time.NewTimer(15 * time.Second)
	// go func() {
	// 	<-timer.C
	// 	fmt.Println(time.Now().String())
	// }()
	time.Sleep(20 * time.Second)
	t.Fatalf("")
}
