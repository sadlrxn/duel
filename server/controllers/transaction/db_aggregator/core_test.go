package db_aggregator

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/lib/pq"
	"gorm.io/gorm/clause"
)

func TestRemoveNftsFromBalanceArray(t *testing.T) {
	{ // Case 1
		// Should remove all containing nfts - [+ + + - -]
		// -: removal, +: leave
		nfts := []Nft{
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
		}
		balanceArray := pq.StringArray{
			"APD4kVba7NQJP8791R8KSkn4xppfNDE2qrr6oAhMx54A",
			"HZdLdpNHrk4QA5ameVixCPdv8ohixWFimi3TEfx2Sh74",
			"HDWp1h9K5hVSox1vDUYoDeZ6AZRC8Hxz76YgEj9bzqP7",
			"29Fj68wD6cogtZgqZ2RJ977WQxapaHd97sEU468nAMPu",
			"C4Ys3WukJEZre2bx17noAejvaNLWz12JCLh8yhb6RZeA",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
		}

		result, err := removeNftsFromBalanceArray(&nfts, &balanceArray)
		resultExp := pq.StringArray{
			"APD4kVba7NQJP8791R8KSkn4xppfNDE2qrr6oAhMx54A",
			"HZdLdpNHrk4QA5ameVixCPdv8ohixWFimi3TEfx2Sh74",
			"HDWp1h9K5hVSox1vDUYoDeZ6AZRC8Hxz76YgEj9bzqP7",
			"29Fj68wD6cogtZgqZ2RJ977WQxapaHd97sEU468nAMPu",
			"C4Ys3WukJEZre2bx17noAejvaNLWz12JCLh8yhb6RZeA",
		}
		if err != nil || len(*result) != len(resultExp) {
			t.Fatalf("Failed %v %v", *result, err)
		}

		allMatch := true
		for i, n := 0, len(*result); i < n; i++ {
			allMatch = allMatch && ((*result)[i] == resultExp[i])
		}

		if !allMatch {
			t.Fatalf("failed to match: %v %d", result, len(*result))
		}
	}

	{ // Case 2
		// Should remove all containing nfts - [- - + + +]
		nfts := []Nft{
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
		}
		balanceArray := pq.StringArray{
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"APD4kVba7NQJP8791R8KSkn4xppfNDE2qrr6oAhMx54A",
			"HZdLdpNHrk4QA5ameVixCPdv8ohixWFimi3TEfx2Sh74",
			"HDWp1h9K5hVSox1vDUYoDeZ6AZRC8Hxz76YgEj9bzqP7",
			"29Fj68wD6cogtZgqZ2RJ977WQxapaHd97sEU468nAMPu",
			"C4Ys3WukJEZre2bx17noAejvaNLWz12JCLh8yhb6RZeA",
		}
		result, err := removeNftsFromBalanceArray(&nfts, &balanceArray)
		resultExp := pq.StringArray{
			"APD4kVba7NQJP8791R8KSkn4xppfNDE2qrr6oAhMx54A",
			"HZdLdpNHrk4QA5ameVixCPdv8ohixWFimi3TEfx2Sh74",
			"HDWp1h9K5hVSox1vDUYoDeZ6AZRC8Hxz76YgEj9bzqP7",
			"29Fj68wD6cogtZgqZ2RJ977WQxapaHd97sEU468nAMPu",
			"C4Ys3WukJEZre2bx17noAejvaNLWz12JCLh8yhb6RZeA",
		}
		if err != nil || len(*result) != len(resultExp) {
			t.Fatalf("Failed %v %v", *result, err)
		}

		allMatch := true
		for i, n := 0, len(*result); i < n; i++ {
			allMatch = allMatch && ((*result)[i] == resultExp[i])
		}

		if !allMatch {
			t.Fatalf("failed to match: %v %d", result, len(*result))
		}
	}

	{ // Case 3
		// Should remove all containing nfts - [+ - + - +]
		nfts := []Nft{
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
		}
		balanceArray := pq.StringArray{
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"HDWp1h9K5hVSox1vDUYoDeZ6AZRC8Hxz76YgEj9bzqP7",
			"29Fj68wD6cogtZgqZ2RJ977WQxapaHd97sEU468nAMPu",
			"APD4kVba7NQJP8791R8KSkn4xppfNDE2qrr6oAhMx54A",
			"HZdLdpNHrk4QA5ameVixCPdv8ohixWFimi3TEfx2Sh74",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
			"C4Ys3WukJEZre2bx17noAejvaNLWz12JCLh8yhb6RZeA",
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
		}
		result, err := removeNftsFromBalanceArray(&nfts, &balanceArray)
		resultExp := pq.StringArray{
			"HDWp1h9K5hVSox1vDUYoDeZ6AZRC8Hxz76YgEj9bzqP7",
			"29Fj68wD6cogtZgqZ2RJ977WQxapaHd97sEU468nAMPu",
			"APD4kVba7NQJP8791R8KSkn4xppfNDE2qrr6oAhMx54A",
			"HZdLdpNHrk4QA5ameVixCPdv8ohixWFimi3TEfx2Sh74",
			"C4Ys3WukJEZre2bx17noAejvaNLWz12JCLh8yhb6RZeA",
		}
		if err != nil || len(*result) != len(resultExp) {
			t.Fatalf("Failed %v %v", *result, err)
		}

		allMatch := true
		for i, n := 0, len(*result); i < n; i++ {
			allMatch = allMatch && ((*result)[i] == resultExp[i])
		}

		if !allMatch {
			t.Fatalf("failed to match: %v %d", result, len(*result))
		}
	}

	{ // Case 4
		// Should leave an empty array - [- - -]
		nfts := []Nft{
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
		}
		balanceArray := pq.StringArray{
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
		}
		result, err := removeNftsFromBalanceArray(&nfts, &balanceArray)
		if err != nil || len(*result) != 0 {
			t.Fatalf("Failed %v %v", *result, err)
		}
	}

	{ // Case 5
		// Should cause an error - [- + - + -]/[- - - ?]
		nfts := []Nft{
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"8bFSU93xdMzqQkpG2CphVCtUE7jV7cv2W1sydG3CyrWV",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
		}
		balanceArray := pq.StringArray{
			"C4Ys3WukJEZre2bx17noAejvaNLWz12JCLh8yhb6RZeA",
			"FcRrqk9RV5ZbZv3tc324D6AyoKStMmdDTA2LETaUQr76",
			"J4VaPkkBYnmXTn27v6zcZnB4SR87WWSSzGHPaNJnviH7",
			"HZdLdpNHrk4QA5ameVixCPdv8ohixWFimi3TEfx2Sh74",
			"5KjtpBXaW9rMxWv6d67bX4ejR1dhv4nocYjx9nTkD9fx",
			"29Fj68wD6cogtZgqZ2RJ977WQxapaHd97sEU468nAMPu",
		}
		result, err := removeNftsFromBalanceArray(&nfts, &balanceArray)
		if err == nil {
			t.Fatalf("Failed %v", *result)
		}
	}

	{ // Case 6
		nfts := []Nft{
			"ESk2VUYmz7LKm9DwMWrxa9AJ4YRiPuPcmMqjNgwmEfwE",
		}
		balanceArray := pq.StringArray{
			"ESk2VUYmz7LKm9DwMWrxa9AJ4YRiPuPcmMqjNgwmEfwE",
			"6mFiFUCDLLiSxTikNtAHtpE2JCgbB7VSkcKv76a7zjTa",
			"EyKa1EXfZk4Fa9kDSSYjjDxDSjYYVbZ6NafH8to3UFQb",
			"7BbLsCh44gqNVyhjXF3umHKUed58Ukkv9nSSDknJZXmZ",
			"EbavJqfhPTcbuErZmGaDuwRBnHjUGFYBZwKyAHbrX8aE",
			"caauoXogcDYVAjGyaEyrKAak7kPyok8KCK59fGgqmnm",
			"CsqXem8HLXnuKaToJcRXXdMbXZgPa1fBdNQMmugwuiCk",
			"Eby6ikiWRevcgWxTDYTA9rc4T9XZQKBnqbiEz69Sh6Y6",
			"371ijDTrtUwgbP5y6nkzCEHjMP3thSsFohr4FJeugPim",
			"DNMKCVGzB6YpoeSn64LinqTFSn51cu8sUPDcEz795tiN",
			"3fPoG48ERaCWxkawQpWsJLmuoRJwzxdrxavrHDGatcrk",
			"6wopTGJvQhCTvG96A15E4Fo7GR2whecRhaSKCPQ2ty2e",
			"BkUf5x3FoZ9FkhatxcrcbZjN34oWV2A7CTTD7sKvsjMR",
			"9z6x4Gm9d4P3j6CF87pZDJtahz5X9aBopyKfTVQVbDQW",
			"WJQZM1vFwbS3tHD9F9awCGEac5BJCyjzDY8g2gyVQQD",
		}

		result, err := removeNftsFromBalanceArray(&nfts, &balanceArray)
		resultExp := pq.StringArray{
			"6mFiFUCDLLiSxTikNtAHtpE2JCgbB7VSkcKv76a7zjTa",
			"EyKa1EXfZk4Fa9kDSSYjjDxDSjYYVbZ6NafH8to3UFQb",
			"7BbLsCh44gqNVyhjXF3umHKUed58Ukkv9nSSDknJZXmZ",
			"EbavJqfhPTcbuErZmGaDuwRBnHjUGFYBZwKyAHbrX8aE",
			"caauoXogcDYVAjGyaEyrKAak7kPyok8KCK59fGgqmnm",
			"CsqXem8HLXnuKaToJcRXXdMbXZgPa1fBdNQMmugwuiCk",
			"Eby6ikiWRevcgWxTDYTA9rc4T9XZQKBnqbiEz69Sh6Y6",
			"371ijDTrtUwgbP5y6nkzCEHjMP3thSsFohr4FJeugPim",
			"DNMKCVGzB6YpoeSn64LinqTFSn51cu8sUPDcEz795tiN",
			"3fPoG48ERaCWxkawQpWsJLmuoRJwzxdrxavrHDGatcrk",
			"6wopTGJvQhCTvG96A15E4Fo7GR2whecRhaSKCPQ2ty2e",
			"BkUf5x3FoZ9FkhatxcrcbZjN34oWV2A7CTTD7sKvsjMR",
			"9z6x4Gm9d4P3j6CF87pZDJtahz5X9aBopyKfTVQVbDQW",
			"WJQZM1vFwbS3tHD9F9awCGEac5BJCyjzDY8g2gyVQQD",
		}
		if err != nil || len(*result) != len(resultExp) {
			t.Fatalf("Failed %v %v", *result, err)
		}

		allMatch := true
		for i, n := 0, len(*result); i < n; i++ {
			allMatch = allMatch && ((*result)[i] == resultExp[i])
		}

		if !allMatch {
			t.Fatalf("failed to match: %v %d", result, len(*result))
		}
	}

}

func TestBurn(t *testing.T) {
	db := tests.InitMockDB(true, true)

	initialize(db)

	user, _ := getMockUser()
	if result := db.Create(&user); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	if err := burn((*User)(&user.ID), 17); err != nil {
		t.Fatalf("failed to burn chips: %v", err)
	}

	userInfo := models.User{}
	if result := db.Preload(
		"Wallet.Balance.ChipBalance",
	).First(&userInfo); result.Error != nil {
		t.Fatalf("failed to retrieve user info: %v", result.Error)
	}

	if userInfo.Wallet.Balance.ChipBalance.Balance != 100-17 {
		t.Fatalf(
			"failed to burn chips properly: %d",
			userInfo.Wallet.Balance.ChipBalance.Balance,
		)
	}
}

func TestLockingPreload(t *testing.T) {
	db := tests.InitMockDB(true, true)

	initialize(db)

	user, _ := getMockUser()
	if result := db.Create(&user); result.Error != nil {
		t.Fatalf("failed to create mock user: %v", result.Error)
	}

	sessionId, err := startSession()
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	_, err = getUserWallet((*User)(&user.ID), true, sessionId)
	if err != nil {
		t.Fatalf("failed to get user wallet: %v", err)
	}
	userInfo := models.User{}
	db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&userInfo, user.ID)

	// wallet, err := getUserWallet((*User)(&user.ID), true, sessionId)
	// if err != nil {
	// 	t.Fatalf("failed to get user wallet: %v", err)
	// }
	// walletInfo := models.Wallet{}
	// db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&walletInfo, wallet)
}
