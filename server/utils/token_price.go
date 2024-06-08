package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

var TokenPrices = syncmap.Map{}

func FetchTokenPrice(tokens map[string]string) (result map[string]float64) {
	config := config.Get()
	result = make(map[string]float64)
	url := config.TokenPriceApi
	for token, address := range tokens {
		// 	url += token + ","
		// }
		res, err := http.Get(fmt.Sprintf("%v/simple/token_price/solana?contract_addresses=%v&vs_currencies=%v", url, address, "usd"))
		if err != nil || res.StatusCode == http.StatusTooManyRequests {
			var errMessage string
			if err != nil {
				errMessage = err.Error()
			}
			log.LogMessage(
				"token price",
				"because of an Error occured fetching token price, returning old one",
				"error",
				logrus.Fields{
					"error": errMessage,
					"url":   url,
				})
			for token, _ := range tokens {
				price, prs := TokenPrices.Load(token)
				if !prs {
					price = float64(0)
				}
				result[token] = price.(float64)
			}
			return
		}

		body, err := io.ReadAll(res.Body)

		if err != nil {
			log.LogMessage("token price", "because of an Error occured fetching token price, returning old one", "info", logrus.Fields{})
			for token, _ := range tokens {
				price, prs := TokenPrices.Load(token)
				if !prs {
					price = float64(0)
				}
				result[token] = price.(float64)
			}
			return
		}
		// var response struct {
		// 	Data map[string]struct {
		// 		ID            string  `json:"id"`
		// 		MintSymbol    string  `json:"mintSymbol"`
		// 		VsToken       string  `json:"vsToken"`
		// 		VsTokenSymbol string  `json:"vsTokenSymbol"`
		// 		Price         float64 `json:"price"`
		// 	} `json:"data"`
		//}
		var response map[string]struct {
			USD float64 `json:"usd"`
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			for token, _ := range tokens {
				price, prs := TokenPrices.Load(token)
				if !prs {
					price = float64(0)
				}
				result[token] = price.(float64)
			}
			return
		}
		// for key, token := range response.Data {
		// 	if key == "USDC" {
		// 		TokenPrices.Store(key, float64(1))
		// 		continue
		// 	}
		// 	TokenPrices.Store(key, token.Price)
		// }
		if token == "USDC" {
			TokenPrices.Store(token, float64(1))
		} else {
			TokenPrices.Store(token, response[address].USD)
		}

		for token, _ := range tokens {
			price, prs := TokenPrices.Load(token)
			if !prs {
				price = float64(0)
			}
			result[token] = price.(float64)
		}
	}
	return
}
