package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
)

func SwapTokens(amount float32, source string, target string) (float32, error) {
	config := config.Get()
	var data = []byte(fmt.Sprintf(`{
		"inputMint": "%s",
		"outputMint": "%s",
		"amount": %f
		}`, source, target, amount))
	req, _ := http.NewRequest("POST", config.JupiterAggregater, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-access-key", config.JupiterApiAccessKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	body, _ := io.ReadAll(res.Body)
	var resBody struct {
		TxID         string  `json:"txid"`
		InputAmount  float32 `json:"inputAmount"`
		OutputAmount float32 `json:"outputAmount"`
	}
	json.Unmarshal(body, &resBody)

	if res.StatusCode == 200 {
		return resBody.OutputAmount, nil
	} else {
		return 0, errors.New("Failed to swap " + source + " to " + target + " : " + string(body))
	}
}

// func SwapSolToUsdc(amount float32) (float32, error) {
// 	config := config.Get()
// 	var data = []byte(fmt.Sprintf(`{
// 		"inputMint": "%s",
// 		"outputMint": "%s",
// 		"amount": %f
// 		}`, config.Sol, config.Usdc, amount))
// 	req, _ := http.NewRequest("POST", config.JupiterAggregater, bytes.NewBuffer(data))
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("x-api-access-key", config.JupiterApiAccessKey)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return 0, err
// 	}
// 	body, _ := io.ReadAll(res.Body)
// 	var resBody struct {
// 		TxID         string  `json:"txid"`
// 		InputAmount  float32 `json:"inputAmount"`
// 		OutputAmount float32 `json:"outputAmount"`
// 	}
// 	json.Unmarshal(body, &resBody)

// 	if res.StatusCode == 200 {
// 		return resBody.OutputAmount, nil
// 	} else {
// 		return 0, errors.New("Failed to swap SOL to USDC : " + string(body))
// 	}
// }

// func SwapUsdcToSol(amount float32) (float32, error) {
// 	config := config.Get()
// 	var data = []byte(fmt.Sprintf(`{
// 		"inputMint": "%s",
// 		"outputMint": "%s",
// 		"amount": %f
// 		}`, config.Usdc, config.Sol, amount))
// 	req, _ := http.NewRequest("POST", config.JupiterAggregater, bytes.NewBuffer(data))
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("x-api-access-key", config.JupiterApiAccessKey)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return 0, err
// 	}
// 	body, _ := io.ReadAll(res.Body)
// 	var resBody struct {
// 		TxID         string  `json:"txid"`
// 		InputAmount  float32 `json:"inputAmount"`
// 		OutputAmount float32 `json:"outputAmount"`
// 	}
// 	json.Unmarshal(body, &resBody)

// 	if res.StatusCode == 200 {
// 		return resBody.OutputAmount, nil
// 	} else {
// 		return 0, errors.New("Failed to swap USDC to SOL : " + string(body))
// 	}
// }
