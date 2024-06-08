package maintenance

import "time"

type MaintenanceStatus string

const (
	NotMaintenance MaintenanceStatus = "not-maintenance"
	// BeforeMaintenance MaintenanceStatus = "before-maintenance"
	InMaintenance MaintenanceStatus = "in-maintenance"
)

type MaintenanceDetails struct {
	Status    MaintenanceStatus
	Timestamp time.Time
}
