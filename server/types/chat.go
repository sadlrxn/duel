package types

import "github.com/Duelana-Team/duelana-v1/models"

type ChatContent struct {
	ID          uint   `json:"id"`
	Author      User   `json:"author"`
	Message     string `json:"message"`
	IsDelegated bool   `json:"isDelegated"`
	ReplyTo     *uint  `json:"replyTo"`
	Sponsors    []uint `json:"sponsors"`
	Deleted     bool   `json:"deleted"`
	Time        uint64 `json:"time"`
}

type ChatCommand struct {
	Pattern string      `json:"pattern"`
	RegExp  string      `json:"regExp"`
	Role    models.Role `json:"role"`
}
