package config

import (
	"time"

	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
)

var API_RATE_LIMIT_CONFIGURATION = map[string]types.RateLimit{
	"affiliate/create": {
		Tokens:   1,
		Interval: time.Second,
	},
	"affiliate/delete": {
		Tokens:   2,
		Interval: time.Second,
	},
	"affiliate/activate": {
		Tokens:   1,
		Interval: time.Second,
	},
	"affiliate/claim": {
		Tokens:   30,
		Interval: time.Hour,
	},
	"affiliate/deactivate": {
		Tokens:   1,
		Interval: time.Second,
	},
	"bot/stake": {
		Tokens:   2,
		Interval: time.Second,
	},
	"bot/unstake": {
		Tokens:   2,
		Interval: time.Second,
	},
	"bot/claim": {
		Tokens:   30,
		Interval: time.Hour,
	},
	"dreamtower/bet": {
		Tokens:   1,
		Interval: time.Second,
	},
	"dreamtower/raise": {
		Tokens:   3,
		Interval: time.Second,
	},
	"dreamtower/cashout": {
		Tokens:   1,
		Interval: time.Second,
	},
	"rewards/rakeback": {
		Tokens:   30,
		Interval: time.Hour,
	},
	"seed/rotate": {
		Tokens:   20,
		Interval: time.Hour,
	},
	"user/tip": {
		Tokens:   10,
		Interval: time.Minute,
	},
	"user/update": {
		Tokens:   30,
		Interval: time.Hour,
	},
	"pay/withdraw/sol": {
		Tokens:   10,
		Interval: time.Hour,
	},
	"pay/withdraw/nft": {
		Tokens:   10,
		Interval: time.Hour,
	},
	"coupon/redeem": {
		Tokens:   2,
		Interval: 10 * time.Second,
	},
	"coupon/claim": {
		Tokens:   2,
		Interval: 1 * time.Minute,
	},
}

var WEBSOCKET_RATE_LIMIT_CONFIGURATION = map[string]types.RateLimit{
	"visit/": {
		Tokens:   1,
		Interval: time.Second,
	},
	"visit/coinflip": {
		Tokens:   1,
		Interval: time.Second,
	},
	"visit/jackpot": {
		Tokens:   1,
		Interval: time.Second,
	},
	"visit/grandJackpot": {
		Tokens:   1,
		Interval: time.Second,
	},
	"visit/crash": {
		Tokens:   1,
		Interval: time.Second,
	},
	"event/coinflip": {
		Tokens:   5,
		Interval: time.Second,
	},
	"event/jackpot": {
		Tokens:   3,
		Interval: time.Second,
	},
	"event/grandJackpot": {
		Tokens:   3,
		Interval: time.Second,
	},
	"event/crash": {
		Tokens:   5,
		Interval: time.Second,
	},
	"message/chat": {
		Tokens:   1,
		Interval: time.Second,
	},
	"reply/chat": {
		Tokens:   1,
		Interval: time.Second,
	},
	"delete/chat": {
		Tokens:   5,
		Interval: time.Second,
	},
	"sponsor/chat": {
		Tokens:   5,
		Interval: time.Second,
	},
}

var BALANCE_DECIMALS = 5

const ONE_CHIP_WITH_DECIMALS = int64(100000)

var JACKPOT_MIN_AMOUNT_LOW = int64(ONE_CHIP_WITH_DECIMALS)            // 1 usd
var JACKPOT_MAX_AMOUNT_LOW = int64(25 * ONE_CHIP_WITH_DECIMALS)       // 25 usd
var JACKPOT_MIN_AMOUNT_MEDIUM = int64(25 * ONE_CHIP_WITH_DECIMALS)    // 25 usd
var JACKPOT_MAX_AMOUNT_MEDIUM = int64(250 * ONE_CHIP_WITH_DECIMALS)   // 250 usd
var JACKPOT_MIN_AMOUNT_WILD = int64(ONE_CHIP_WITH_DECIMALS)           // 1 usd
var JACKPOT_MAX_AMOUNT_WILD = int64(1000000 * ONE_CHIP_WITH_DECIMALS) // 100k usd
var JACKPOT_BET_COUNT_LIMIT = uint(5)                                 // Up to 5 bets
var JACKPOT_PLAYER_LIMIT = uint(50)                                   // 50 players
var JACKPOT_COUNTING_TIME = uint(40)                                  // 40 s
var JACKPOT_ROLLING_TIME = uint(30)                                   // 30 s
var JACKPOT_FEE = int64(5)                                            // 5 %
var JACKPOT_TEMP_ID = uint(1001)
var JACKPOT_FEE_ID = uint(1002)
var JACKPOT_TAIL = uint(3)       // 3 s
var JACKPOT_EXTRA_TIME = uint(3) // 3 s

var GRAND_JACKPOT_MIN_AMOUNT = int64(10)

var GRAND_JACKPOT_ROLLING_TIME = uint(1 * 60 * 60)
var GRAND_JACKPOT_COUNTING_TIME = uint(0 * 60)
var GRAND_JACKPOT_BETTING_TIME = uint(23 * 60 * 60)
var GRAND_JACKPOT_ROLLING_DURATION = uint(30)
var GRAND_JACKPOT_FEE = int64(5)
var GRAND_JACKPOT_TEMP_ID = uint(1003)
var GRAND_JACKPOT_FEE_ID = uint(1004)
var GRAND_JACKPOT_NEXT_ROUND = time.Date(2023, 1, 14, 2, 0, 0, 0, time.UTC)

var COINFLIP_ROUND_LIMIT = uint(5)
var COINFLIP_MIN_AMOUNT = int64(float64(0.01) * float64(ONE_CHIP_WITH_DECIMALS))
var COINFLIP_MAX_AMOUNT = int64(250 * ONE_CHIP_WITH_DECIMALS)
var COINFLIP_FEE = int64(2)
var COINFLIP_TEMP_ID = uint(1005)
var COINFLIP_FEE_ID = uint(1006)
var COINFLIP_BOT_ID = uint(1007)

var DREAMTOWER_MIN_AMOUNT = int64(float64(0.01) * float64(ONE_CHIP_WITH_DECIMALS))
var DREAMTOWER_MAX_AMOUNT = int64(100 * ONE_CHIP_WITH_DECIMALS)
var DREAMTOWER_TEMP_ID = uint(1008)
var DREAMTOWER_FEE_ID = uint(1009)
var DREAMTOWER_FEE = int64(3)

var CRASH_TEMP_ID = uint(100002)
var CRASH_FEE_ID = uint(100003)
var CRASH_EVENT_INTERVAL_MILLI = int64(200)
var CRASH_BETTING_DURATION_MILLI = int64(7000)
var CRASH_PENDING_DURATION_MILLI = int64(3000)
var CRASH_PREPARING_DURATION_MILLI = int64(3000)
var CRASH_BET_COUNT_LIMIT = uint(5)
var CRASH_MIN_BET_AMOUNT = int64(float64(0.01) * float64(ONE_CHIP_WITH_DECIMALS))
var CRASH_MAX_BET_AMOUNT = int64(100 * ONE_CHIP_WITH_DECIMALS)
var CRASH_MULTIPLIER_INCREASE_RATE = float64(1.012)
var CRASH_HOUSE_EDGE = int64(500)
var CRASH_MAX_PLAYER_LIMIT = uint(1000)
var CRASH_MIN_CASH_OUT_AT = float64(1.01)
var CRASH_SEED_CHAIN_LENGTH = uint(2500000)
var CRASH_MAX_CASH_OUT = int64(1000 * ONE_CHIP_WITH_DECIMALS)
var CRASH_START_ON_SERVER_STARTUP = false

var BASE_RAKEBACK_RATE = uint(5)       // 5 %
var ADDITIONAL_RAKEBACK_RATE = uint(0) // 0 %
var RAKEBACK_MAX = uint(10)            // 10 %

var CHAT_MAX_COUNT = int(100)                                 // 100 messages
var CHAT_MAX_LENGTH = uint(200)                               // 200 letters
var CHAT_WAGER_LIMIT = int64(50 * ONE_CHIP_WITH_DECIMALS)     // 50 usd
var CHAT_COOL_DOWN = int(0)                                   // 0 s
var CHAT_RAIN_MIN_WAGER = int64(100 * ONE_CHIP_WITH_DECIMALS) // 100 usd

var WITHDRAW_MIN_LIMIT = int64(ONE_CHIP_WITH_DECIMALS)                           // 1 usd
var WITHDRAW_FEE_PER_SPL = int64(float64(0.1) * float64(ONE_CHIP_WITH_DECIMALS)) // 0.1 usd

const DREAMTOWER_HEIGHT = uint(9)

var DUEL_BOT_STAKE_ID = uint(10001)
var DUEL_BOT_TOTAL_SHARE = int64(80) // 80 %

var AFFILIATE_RATE_MIN = uint(5)                                          // 5 %
var AFFILIATE_RATE_MAX = uint(20)                                         // 20 %
var AFFILIATE_WAGER_LIMIT_FOR_CREATION = uint(0 * ONE_CHIP_WITH_DECIMALS) // 0 usd should be set in chips with 5 decimals
var AFFILIATE_RESERVED_WORDS = []string{
	"duel",
}
var AFFILIATE_ACTIVATION_TIMELINE_IN_HOURS = 24 //24 hours

var RAIN_MAX_SPLIT_COUNT = int64(100)

var COUPON_REQUIRED_WAGER_TIMES = 30                                  // claim x 30
var COUPON_BALANCE_LIFE_TIME_IN_HOURS = 8                             // 8 hours
var COUPON_CODE_LIFE_TIME_IN_DAYS = 14                                // 14 days
var COUPON_EXCHANGE_CHIP_RATE = 100                                   // 100 %
var COUPON_MAXIMUM_EXCHANGE = 50 * ONE_CHIP_WITH_DECIMALS             // Maximum exchange amount
var COUPON_MAXIMUM_FIRST_DEPOSIT_BONUS = 100 * ONE_CHIP_WITH_DECIMALS // Maximum bonus balance on first deposit
var COUPON_TEMP_ID = uint(100001)

var HIDDEN_USERS = []string{"DuelBot"}
var ACCOUNT_LIMIT_PER_IP = uint(20)
var ACCOUNT_LIMIT_FOR_SPECIFIC_IP = map[string]uint{}

var DAILY_RACE_TEMP_ID = uint(100004)

var WEEKLY_RAFFLE_TEMP_ID = uint(100005)                                // maximum reserved user_id flag here
var WEEKLY_RAFFLE_MAXIMUM_TICKET_ID = uint(9999999)                     // Currently set as infinite.
var WEEKLY_RAFFLE_MAXIMUM_TICKET_PER_USER = uint(9999999)               // Currently set as infinite
var WEEKLY_RAFFLE_MAXIMUM_PARTICIPANTS = uint(9999999)                  // Currently set as infinite
var WEEKLY_RAFFLE_CHIPS_WAGER_PER_TICKET = 500 * ONE_CHIP_WITH_DECIMALS // 100 chips per ticket
var WEEKLY_RAFFLE_DURATION_IN_DAYS = 7                                  // 7 days playing
var WEEKLY_RAFFLE_PENDING_IN_MINUTES = 60                               // 60 minutes pending before starting
var WEEKLY_RAFFLE_DEFAULT_PRIZES = []int64{                             // One time used weekly raffle prizes.
	1000 * ONE_CHIP_WITH_DECIMALS,
	700 * ONE_CHIP_WITH_DECIMALS,
	300 * ONE_CHIP_WITH_DECIMALS,
}
var WEEKLY_RAFFLE_OPEN = true

var DREAMTOWER_DIFFICULTIES = map[string]models.DreamTowerDifficulty{
	"Easy": {
		Level:       models.LevelEasy,
		BlocksInRow: 4,
		StarsInRow:  3,
	},
	"Medium": {
		Level:       models.LevelMedium,
		BlocksInRow: 3,
		StarsInRow:  2,
	},
	"Hard": {
		Level:       models.LevelHard,
		BlocksInRow: 2,
		StarsInRow:  1,
	},
	"Expert": {
		Level:       models.LevelExpert,
		BlocksInRow: 3,
		StarsInRow:  1,
	},
	// "Master": {
	// 	Level:       models.LevelMaster,
	// 	BlocksInRow: 4,
	// 	StarsInRow:  1,
	// },
}

var TIP_MIN_AMOUNT = int64(float64(0.01) * float64(ONE_CHIP_WITH_DECIMALS))
var TIP_MAX_AMOUNT = int64(10000 * ONE_CHIP_WITH_DECIMALS)
var MUTE_DURATION = time.Duration(15) * time.Minute

var CHAT_COMMANDS = []types.ChatCommand{
	{Pattern: "/mute user", RegExp: `/\/mute \w+/`, Role: models.ModeratorRole},
	{Pattern: "/unmute user", RegExp: `/\/unmute \w+/`, Role: models.ModeratorRole},
	{Pattern: "/ban user", RegExp: `/\/ban \w+/`, Role: models.ModeratorRole},
	{Pattern: "/unban user", RegExp: `/\/unban \w+/`, Role: models.ModeratorRole},
	{Pattern: "/setWagerLimit amount", RegExp: `/\/setWagerLimit \d+/`, Role: models.AdminRole},
	{Pattern: "/setMaxLength amount", RegExp: `/\/setMaxLength \d+/`, Role: models.AdminRole},
	{Pattern: "/tip user amount", RegExp: `/\/tip \w+ \d+/`, Role: models.UserRole},
	{Pattern: "/details user", RegExp: `/\/details \w+/`, Role: models.UserRole},
	{Pattern: "/setChatCooldown time", RegExp: `/\/setChatCooldown \d+/`, Role: models.AdminRole},
	{Pattern: "/rain split amount", RegExp: `/\/rain \d+ \d+/`, Role: models.UserRole},
}

var BLACK_LIST = map[string]bool{"AF": true, "AU": true, "BY": true, "BE": true, "CI": true, "CU": true, "CW": true, "CZ": true, "CD": true, "FR": true, "DE": true, "GR": true, "IR": true, "IQ": true, "IT": true, "LR": true, "LY": true, "LT": true, "NL": true, "KP": true, "PT": true, "RS": true, "SK": true, "SS": true, "ES": true, "SD": true, "SE": true, "SY": true, "GB": true, "US": true, "ZW": true}

var BAD_WORDS = []string{
	"hidden", "anal", "anus", "arse", "ass", "ballsack", "balls", "bastard", "bitch", "biatch", "bloody", "blowjob", "blow", "job", "bollock", "bollok", "boner", "boob", "bugger", "clitoris", "cock", "coon", "cunt", "dick", "dildo", "dyke", "fag", "fellate", "fellatio", "felching", "fuck", "f u c k", "fudgepacker", "fudge packer", "flange", "homo", "jerk", "jizz", "knobend", "knob end", "labia", "muff", "nigger", "nigga", "penis", "piss", "poop", "prick", "pube", "pussy", "queer", "scrotum", "sex", "shit", "s hit", "sh1t", "slut", "smegma", "spunk", "tit", "tosser", "turd", "twat", "vagina", "wank", "whore",
}

var BONK_SPL_ADDRESS = "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"
var USDC_SPL_ADDRESS = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
var SOL_SPL_ADDRESS = "So11111111111111111111111111111111111111112"
var BOKU_SPL_ADDRESS = "CN7qFa5iYkHz99PTctvT4xXUHnxwjQ5MHxCuTJtPN5uS"
