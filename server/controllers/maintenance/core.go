package maintenance

import (
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/utils"
)

// @Internal
// Saves the detail of the maintenance status.
var _details MaintenanceDetails

// @Internal
// Required constants.
const DEFAULT_STATUS = NotMaintenance

// @Internal
// Initializes maintenance module.
func initialize() error {
	updateStatus(DEFAULT_STATUS, true)
	return nil
}

// @Internal
// Makes a customized error object.
func makeError(category string, reason string, err error) error {
	return utils.MakeError("maintenance-aggregator", category, reason, err)
}

// @Internal
// Check for the state update validity.
// Status Update rules: (exclusive)
// 1. not-maintenance => before-maintenance
// 2. before-maintenance => in-maintenance // 2nd Rule Archived
// 3. in-maintenance => not-maintenance
func isStatusUpdateable(updateStatus MaintenanceStatus) bool {
	// switch updateStatus {
	// case NotMaintenance:
	// 	if currentStatus() != InMaintenance {
	// 		return false
	// 	}
	// case BeforeMaintenance:
	// 	if currentStatus() != NotMaintenance {
	// 		return false
	// 	}
	// case InMaintenance:
	// 	if currentStatus() != BeforeMaintenance {
	// 		return false
	// 	}
	// }
	switch updateStatus {
	case NotMaintenance:
		if currentStatus() != InMaintenance {
			return false
		}
	case InMaintenance:
		if currentStatus() != NotMaintenance {
			return false
		}
	}

	return true
}

// @Internal
// Updates maintenance status.
// Automatically updates the timestamp as the current time.
// Status Update rules: (exclusive)
// 1. not-maintenance => before-maintenance
// 2. before-maintenance => in-maintenance
// 3. in-maintenance => not-maintenance
// If the forceUpdate param is true, checking for rule is not required.
func updateStatus(status MaintenanceStatus, forceUpdate ...bool) error {
	if ((len(forceUpdate) > 0 && !forceUpdate[0]) || len(forceUpdate) == 0) && !isStatusUpdateable(status) {
		return makeError("updateStatus", "the update is against the rules", fmt.Errorf("cannot update from %v status to %v status", currentStatus(), status))
	}

	setStatus(status)
	setTimestamp(time.Now())

	return nil
}

// @Internal
// Get current status.
func currentStatus() MaintenanceStatus {
	return _details.Status
}

// @Internal
// Simply sets maintenance status.
func setStatus(status MaintenanceStatus) {
	_details.Status = status
}

// @Internal
// Simply sets timestamp.
func setTimestamp(timestamp time.Time) {
	_details.Timestamp = timestamp
}

// @External // Archived
// Updates the maintenance details to prepare the maintenance.
// func prepare() error {
// 	if err := updateStatus(BeforeMaintenance); err != nil {
// 		return makeError("prepare", "failed to update status for preparing maintenance", err)
// 	}
// 	return nil
// }

// @External
// Updates the maintenance details to maintain the maintenance.
func maintain() error {
	if err := updateStatus(InMaintenance); err != nil {
		return makeError("maintain", "failed to update status for maintaining maintenance", err)
	}
	return nil
}

// @External
// Updates the maintenance details to finish the maintenance.
func finish() error {
	if err := updateStatus(NotMaintenance); err != nil {
		return makeError("finish", "failed to update status for finishing maintenance", err)
	}
	return nil
}

// @External
// Get the current maintenance details.
func current() MaintenanceDetails {
	return _details
}

// @External
// Returns whether the current status allows new rounds.
func ableToBet() bool {
	return currentStatus() == NotMaintenance
}
