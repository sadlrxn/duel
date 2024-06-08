package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	mathRand "math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/sirupsen/logrus"
)

type randomCreateTicketIDResult struct {
	TicketID         string `json:"ticketId"`
	CreationTime     string `json:"creationTime"`
	PreviousTicketID string `json:"previousTicketId"`
	NextTicketID     string `json:"nextTIcketId"`
}

type randomCreateTIcketIDResponse struct {
	JsonRpc string                       `json:"jsonrpc"`
	Results []randomCreateTicketIDResult `json:"result"`
	ID      uint                         `json:"id"`
}

type randomResponseLicense struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	InfoUrl string `json:"infoUrl"`
}

type randomGenerateStringRandom struct {
	Method       string                `json:"string"`
	HashedApiKey string                `json:"hashedApiKey"`
	N            uint                  `json:"n"`
	Length       uint                  `json:"length"`
	Characters   string                `json:"characters"`
	Replacement  bool                  `json:"replacement"`
	Data         []string              `json:"data"`
	License      randomResponseLicense `json:"license"`
}

type randomGenerateStringResult struct {
	Random    randomGenerateStringRandom `json:"random"`
	Signature string                     `json:"signature"`
}

type randomGenerateStringResponse struct {
	JsonRpc string                     `json:"jsonrpc"`
	Result  randomGenerateStringResult `json:"result"`
	ID      uint                       `json:"id"`
}

var randomRequestID uint = 0

const MAX_RANDOM_REQUEST_ID uint = 10000
const CHARACTERS_FOR_RANDOM_STRING = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const CHARACTERS_LENGTH = 16

func getRandomRequestID() uint {
	randomRequestID = (randomRequestID + 1) % MAX_RANDOM_REQUEST_ID
	return randomRequestID
}

func RequestTicketIDOnce() (*randomCreateTIcketIDResponse, error) {
	config := config.Get()

	jsonStr := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "createTickets",
		"params": {
			"apiKey": "%s",
			"n": 1,
			"showResult": true
		},
		"id": %d
	}`, config.RandomKey, getRandomRequestID())

	res, err := http.Post(config.RandomUrl, "application/json", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return nil, MakeError(
			"random.go",
			"RequestTicketIDOnce",
			"failed to send post request",
			fmt.Errorf(
				"res: %v, err: %v",
				res, err,
			),
		)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, MakeError(
			"random.go",
			"RequestTicketIDOnce",
			"failed to parse request body",
			fmt.Errorf(
				"res: %v, body: %v, err: %v",
				res, body, err,
			),
		)
	}

	resObject := randomCreateTIcketIDResponse{}
	json.Unmarshal(body, &resObject)
	if len(resObject.Results) == 0 {
		return nil, MakeError(
			"random.go",
			"RequestTicketIDOnce",
			"failed to send post request",
			fmt.Errorf(
				"res: %v, body: %v, resObject: %v",
				res, body, resObject,
			),
		)
	}

	return &resObject, nil
}

func RequestTicketID() (string, error) {
	for i := 0; i < 5; i++ {
		resObject, err := RequestTicketIDOnce()
		if err != nil {
			log.LogMessage(
				"random.go",
				"RequestTicketID",
				"failed to request ticket id",
				logrus.Fields{
					"error": err.Error(),
				},
			)
			continue
		}
		log.LogMessage(
			"random.org",
			"ticket id generated",
			"success",
			logrus.Fields{
				"ticket": resObject.Results[0].TicketID,
			},
		)
		return resObject.Results[0].TicketID, nil
	}

	return "", MakeError(
		"random.go",
		"RequestTicketID",
		"failed to generate random ticket",
		errors.New("5 retrying all failed"),
	)
}

func GenerateRandomString(ticketID string) (string, error) {
	config := config.Get()

	jsonStr := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "generateSignedStrings",
		"params": {
			"apiKey": "%s",
			"n": 1,
			"length": %d,
			"characters": "%s",
			"replacement": true,
			"userData": null,
			"ticketId": "%s"
		},
		"id": %d
	}`, config.RandomKey, CHARACTERS_LENGTH, CHARACTERS_FOR_RANDOM_STRING, ticketID, getRandomRequestID())

	res, err := http.Post(config.RandomUrl, "application/json", bytes.NewBuffer([]byte(jsonStr)))

	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	resObject := randomGenerateStringResponse{}
	json.Unmarshal(body, &resObject)

	if len(resObject.Result.Random.Data) == 0 {
		return "", errors.New("Failed to generate random string")
	}
	log.LogMessage("random.org", "random string generated", "success", logrus.Fields{"string": resObject.Result.Random.Data[0]})
	return resObject.Result.Random.Data[0], nil
}

func CalculateRandomOutput(randomString string) uint64 {
	sum := sha256.Sum256([]byte(randomString))

	randomUint := binary.BigEndian.Uint64(sum[:8])
	log.LogMessage("random.org", "random number calculated", "info", logrus.Fields{"number": randomUint})
	return randomUint
}

type WinnerCandidate[T any] struct {
	Weight uint64
	ID     uint
	Entity T
}

type CandidateWithCount[T any] struct {
	Entity T
	Count  uint8
}

type PickWinnerResult[T any] struct {
	CandidatesWithCount []CandidateWithCount[T]
	Winner              T
}

type WinnerCandidates[T any] []WinnerCandidate[T]

func (candidates WinnerCandidates[T]) Len() int {
	return len(candidates)
}

func (candidates WinnerCandidates[T]) Swap(i, j int) {
	candidates[i], candidates[j] = candidates[j], candidates[i]
}

func (candidates WinnerCandidates[T]) Less(i, j int) bool {
	return candidates[i].Weight < candidates[j].Weight || candidates[i].Weight == candidates[j].Weight && candidates[i].ID > candidates[j].ID
}

func GenerateWinnerWithArray[T any](randomString string, candidates WinnerCandidates[T], expectedEntityCount uint) PickWinnerResult[T] {
	sort.Sort(candidates)

	var weights []uint64
	var totalWeight uint64
	for _, candidate := range candidates {
		weights = append(weights, candidate.Weight)
		totalWeight += candidate.Weight
	}

	gcd := findGCD(weights)
	totalWeight /= gcd
	for i := range weights {
		weights[i] /= gcd
	}

	randomOutput := CalculateRandomOutput(randomString) % totalWeight
	winnerIndex := 0

	for randomOutput >= weights[winnerIndex] && winnerIndex < candidates.Len() {
		randomOutput -= weights[winnerIndex]
		winnerIndex++
	}

	if winnerIndex == candidates.Len() {
		winnerIndex--
	}

	result := PickWinnerResult[T]{
		Winner: candidates[winnerIndex].Entity,
	}
	result.CandidatesWithCount = make([]CandidateWithCount[T], candidates.Len())

	weightPerCount := totalWeight * gcd / uint64(expectedEntityCount)

	for index, candidate := range candidates {
		result.CandidatesWithCount[index] = CandidateWithCount[T]{
			Count:  uint8(math.Ceil(float64(weights[index] * gcd / weightPerCount))),
			Entity: candidate.Entity,
		}
	}

	return result
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateServerSeed(n int) (string, string, error) {
	b, err := generateRandomBytes(n)
	hashBytes := sha256.Sum256(b)
	return hex.EncodeToString(b), hex.EncodeToString(hashBytes[:]), err
}

func GenerateClientSeed(n int) (string, error) {
	b, err := generateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func gcd(a uint64, b uint64) uint64 {
	if a == 0 {
		return b
	}
	return gcd(b%a, a)
}

func findGCD(arr []uint64) uint64 {
	var result = arr[0]
	for i := 1; i < len(arr); i++ {
		result = gcd(arr[i], result)

		if result == 1 {
			return 1
		}
	}
	return result
}

func ShuffleSlice[T any](slice []T) {
	mathRand.Seed(time.Now().UnixNano())
	mathRand.Shuffle(len(slice), func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })
}
