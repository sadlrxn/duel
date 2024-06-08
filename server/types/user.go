package types

import "github.com/Duelana-Team/duelana-v1/models"

type User struct {
	ID            uint        `json:"id"`
	Name          string      `json:"name"`
	Role          models.Role `json:"role"`
	Avatar        string      `json:"avatar"`
	Count         uint        `json:"count,omitempty"`
	WalletAddress string      `json:"walletAddress"`
	Banned        bool        `json:"banned"`
	Muted         bool        `json:"muted"`
}

type Users []User

func (users Users) Len() int {
	return len(users)
}

func (users Users) Swap(i, j int) {
	users[i], users[j] = users[j], users[i]
}

func (users Users) Less(i, j int) bool {
	return users[i].ID < users[j].ID
}
