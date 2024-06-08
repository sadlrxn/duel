package admin

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/utils"
)

type UserWagerStatus struct {
	UserName        string `json:"userName"`
	DreamTowerWager uint   `json:"dreamTowerWager"`
	DreamTowerWin   uint   `json:"dreamTowerWin"`
	CoinflipWager   uint   `json:"coinflipWager"`
	CoinflipWin     uint   `json:"coinflipWin"`
}

func getUserLoss(userName string, startTime time.Time, endTime time.Time) (*UserWagerStatus, error) {
	db := db.GetDB()
	if db == nil {
		return nil, utils.MakeError(
			"admin_db_aggregator",
			"getUserLoss",
			"failed to retrieve db pointer",
			errors.New("retrieved db pointer is nil"),
		)
	}

	if userName == "" {
		return nil, utils.MakeError(
			"admin_db_aggregator",
			"getUserLoss",
			"failed to get user loss",
			errors.New("provided user name is empty or start time is after end time"),
		)
	}

	const getUserLossQuery = `select coinflip_state.name user_name, coinflip_state.coinflip_wager, coinflip_state.coinflip_win,
dream_tower_state.dream_tower_wager, dream_tower_state.dream_tower_win from
(select coinflip_bet.name, coinflip_bet.coinflip_wager, coinflip_profit.coinflip_win from
(
	select '%s' as name, sum(balance) coinflip_wager from chip_balances where id in
	(
		select chip_balance_id from balances where owner_type = 'in-transaction' and owner_id in
		(
			select id from transactions where from_wallet in
			(
				select id from wallets where user_id in
				(
					select id from users where name = '%s'
				)
			) and 
			type = 'coinflip_bet' and
			created_at > '%s' and
			created_at < '%s'
		)
	)
) coinflip_bet
left join
(
	select '%s' as name, sum(balance) coinflip_win from chip_balances where id in
	(
		select chip_balance_id from balances where owner_type = 'in-transaction' and owner_id in
		(
			select id from transactions where to_wallet in
			(
				select id from wallets where user_id in
				(
					select id from users where name = '%s'
				)
			) and 
			type = 'coinflip_profit' and
			created_at > '%s' and
			created_at < '%s'
		)
	)
) coinflip_profit
on coinflip_bet.name = coinflip_profit.name) coinflip_state
left join
(select dream_tower_bet.name, dream_tower_bet.dream_tower_wager, dream_tower_profit.dream_tower_win from
(
	select '%s' as name, sum(balance) dream_tower_wager from chip_balances where id in
	(
		select chip_balance_id from balances where owner_type = 'in-transaction' and owner_id in
		(
			select id from transactions where from_wallet in
			(
				select id from wallets where user_id in
				(
					select id from users where name = '%s'
				)
			) and 
			type = 'dreamtower_bet' and
			created_at > '%s' and
			created_at < '%s'
		)
	)
) dream_tower_bet
left join
(
	select '%s' as name, sum(balance) dream_tower_win from chip_balances where id in
	(
		select chip_balance_id from balances where owner_type = 'in-transaction' and owner_id in
		(
			select id from transactions where to_wallet in
			(
				select id from wallets where user_id in
				(
					select id from users where name = '%s'
				)
			) and 
			type = 'dreamtower_profit' and
			created_at > '%s' and
			created_at < '%s'
		)
	)
) dream_tower_profit
on dream_tower_bet.name = dream_tower_profit.name) dream_tower_state
on coinflip_state.name = dream_tower_state.name`

	queryWithArgs := fmt.Sprintf(
		getUserLossQuery,
		userName,
		userName,
		startTime.String()[:10],
		endTime.String()[:10],
		userName,
		userName,
		startTime.String()[:10],
		endTime.String()[:10],
		userName,
		userName,
		startTime.String()[:10],
		endTime.String()[:10],
		userName,
		userName,
		startTime.String()[:10],
		endTime.String()[:10],
	)

	userState := UserWagerStatus{}
	if result := db.Raw(
		queryWithArgs,
	).Scan(&userState); result.Error != nil {
		return nil, utils.MakeError(
			"admin_db_aggregator",
			"getUserLoss",
			"failed to get user wager status",
			result.Error,
		)
	}

	return &userState, nil
}
