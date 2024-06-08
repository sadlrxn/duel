package types

import (
	"time"

	"github.com/Duelana-Team/duelana-v1/models"
)

// update balance
type UpdateType int

const (
	Increase UpdateType = 1
	Decrease UpdateType = -1
)

type BalanceUpdatePayload struct {
	UpdateType  UpdateType                `json:"type"`
	Nfts        []NftDetails              `json:"nfts"`
	Balance     int64                     `json:"balance"`
	Wagered     int64                     `json:"wagered"`
	BalanceType models.PaidBalanceForGame `json:"balanceType"`
	Delay       float32                   `json:"delay"`
}

// coinflip
type CoinflipRoundDataPayload struct {
	RoundID         uint                      `json:"roundId"`
	EndedAt         time.Time                 `json:"endedAt"`
	HeadsUser       User                      `json:"headsUser"`
	TailsUser       User                      `json:"tailsUser"`
	Amount          int64                     `json:"amount"`
	Prize           int64                     `json:"prize"`
	TicketID        string                    `json:"ticketId"`
	SignedString    string                    `json:"signedString"`
	WinnerID        uint                      `json:"winnerId"`
	CreatorID       uint                      `json:"creatorId"`
	PaidBalanceType models.PaidBalanceForGame `json:"paidBalanceType"`
}

type CoinflipRoundDataPayloads []CoinflipRoundDataPayload

func (payloads CoinflipRoundDataPayloads) Len() int { return len(payloads) }

func (payloads CoinflipRoundDataPayloads) Swap(i, j int) {
	payloads[i], payloads[j] = payloads[j], payloads[i]
}

func (payloads CoinflipRoundDataPayloads) Less(i, j int) bool {
	return payloads[i].RoundID < payloads[j].RoundID
}

// jackpot
type JackpotBetPayload struct {
	RoundID   uint         `json:"roundId"`
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	Avatar    string       `json:"avatar"`
	Role      models.Role  `json:"role"`
	UsdAmount int64        `json:"usdAmount"`
	NftAmount int64        `json:"nftAmount"`
	Nfts      []NftDetails `json:"nfts"`
}

type JackpotEndPayload struct {
	RoundID         uint         `json:"roundId"`
	TicketID        string       `json:"ticketId"`
	Winner          User         `json:"winner"`
	Chance          float32      `json:"chance"`
	UsdProfit       int64        `json:"usdProfit"`
	NftProfit       []NftDetails `json:"nftProfit"`
	UsdFee          int64        `json:"usdFee"`
	NftFee          []NftDetails `json:"nftFee"`
	Prize           int64        `json:"prize"`
	Candidates      []User       `json:"candidates"`
	RollingDuration uint64       `json:"rollingDuration"`
}

type PlayerInJackpotRound struct {
	ID        uint         `json:"id"`
	Role      models.Role  `json:"role"`
	Name      string       `json:"name"`
	Avatar    string       `json:"avatar"`
	UsdAmount int64        `json:"usdAmount"`
	NftAmount int64        `json:"nftAmount"`
	Nfts      []NftDetails `json:"nfts"`
	BetCount  uint         `json:"betCount"`
}

type JackpotPayload struct {
	Status          string                 `json:"status,omitempty"`
	RoundID         uint                   `json:"roundId,omitempty"`
	TicketID        *string                `json:"ticketId,omitempty"`
	SignedString    *string                `json:"signedString,omitempty"`
	Players         []PlayerInJackpotRound `json:"players,omitempty"`
	Winner          User                   `json:"winner,omitempty"`
	Offset          uint64                 `json:"offset,omitempty"`
	UsdProfit       int64                  `json:"usdProfit,omitempty"`
	NftProfit       []NftDetails           `json:"nftProfit,omitempty"`
	UsdFee          int64                  `json:"usdFee,omitempty"`
	NftFee          []NftDetails           `json:"nftFee,omitempty"`
	Prize           int64                  `json:"prize,omitempty"`
	Candidates      []User                 `json:"candidates,omitempty"`
	RollingDuration uint64                 `json:"rollingDuration,omitempty"`
	CountingTime    uint                   `json:"countingTime,omitempty"`
}

type JackpotHistoryPayload struct {
	RoundID      uint      `json:"roundId"`
	TicketID     string    `json:"ticketId"`
	SignedString string    `json:"signedString"`
	EndedAt      time.Time `json:"endedAt"`
	Winner       User      `json:"winner"`
	Players      []User    `json:"players"`
	Chance       float32   `json:"chance"`
	Prize        int64     `json:"prize"`
}

// payment
type PaymentHistoryPayload struct {
	Time      time.Time            `json:"time"`
	Type      string               `json:"type"`
	Status    models.PaymentStatus `json:"status"`
	SolDetail models.SolDetail     `json:"solDetail"`
	NftDetail []NftDetails         `json:"nftDetail"`
	TxID      string               `json:"txId"`
}

// error
type ErrorCode string

const (
	ErrLocked          ErrorCode = "locked"
	ErrBanned          ErrorCode = "banned"
	ErrMuted           ErrorCode = "muted"
	ErrNotEnoughPeople ErrorCode = "notEnoughPeople"
	ErrNotEnoughChip   ErrorCode = "notEnoughChip"
)

type ErrorMessagePayload struct {
	Message   string    `json:"message"`
	RoundID   uint      `json:"roundId"`
	ErrorCode ErrorCode `json:"errorCode"`
}
