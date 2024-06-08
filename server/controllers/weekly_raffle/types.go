package weekly_raffle

import "time"

type WinnerInWeeklyRafflePrizingResult struct {
	UserID   uint  `json:"userId"`
	TicketID uint  `json:"ticketId"`
	Rank     uint  `json:"rank"`
	Prize    int64 `json:"prize"`
}

type WeeklyRafflePrizingResult struct {
	StartedAt time.Time                           `json:"startedAt"`
	EndedAt   time.Time                           `json:"endedAt"`
	Winners   []WinnerInWeeklyRafflePrizingResult `json:"winners"`
}

type WeeklyRaffleTicketInStatus struct {
	CreatedAt time.Time `json:"date"`
	TicketID  uint      `json:"ticketId"`
}
type UserInWeeklyRaffleStatus struct {
	ID          uint                         `json:"id"`
	Name        string                       `json:"name"`
	Avatar      string                       `json:"avatar"`
	Rank        int                          `json:"rank"`
	TicketCount uint                         `json:"ticketCount"`
	Tickets     []WeeklyRaffleTicketInStatus `json:"tickets"`
}

type WeeklyRaffleRunningStatus string

const (
	WeeklyRaffleStatusRunning WeeklyRaffleRunningStatus = "running"
	WeeklyRaffleStatusPending WeeklyRaffleRunningStatus = "pending"
)

type WeeklyRaffleStatus struct {
	Me                 UserInWeeklyRaffleStatus   `json:"me"`
	Players            []UserInWeeklyRaffleStatus `json:"players"`
	TotalTicketsIssued uint                       `json:"totalTickets"`
	Remaining          uint                       `json:"remaining"`
	ChipsPerTicket     int64                      `json:"chipsPerTicket"`
	TotalPrize         int64                      `json:"totalPrize"`
	Status             WeeklyRaffleRunningStatus  `json:"status"`
}

type WeeklyRaffleRewardsStatus struct {
	ID    uint      `json:"id"`
	Date  time.Time `json:"date"`
	Prize int64     `json:"prize"`
	Rank  uint      `json:"rank"`
}
