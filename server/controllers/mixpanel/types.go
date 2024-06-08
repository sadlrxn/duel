package mixpanel

import "github.com/Duelana-Team/duelana-v1/models"

type TransactionBody struct {
	FromUser    *uint
	ToUser      *uint
	ChipBalance *int64
	NftBalance  *int64
	Type        models.TransactionType
}
