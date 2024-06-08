package user

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
)

func TestSetAccountLimitForSpecificIp(t *testing.T) {
	if err := setAccountLimitForSpecificIp("", 100); err == nil ||
		!strings.Contains(err.Error(), "invalid parameter") {
		t.Fatalf("should be failed in empty string: %v", err)
	}
	if err := setAccountLimitForSpecificIp("some-ip-address", 0); err != nil {
		t.Fatalf("failed to set account limit for specific ip: %v", err)
	}
	if err := setAccountLimitForSpecificIp("some-ip-address", 100); err != nil {
		t.Fatalf("failed to set account limit for specific ip: %v", err)
	}
}

func TestGetAccountLimitForSpecificIp(t *testing.T) {
	if limit := getAccountLimitForSpecificIp("some-ip-address"); limit != config.ACCOUNT_LIMIT_PER_IP {
		t.Fatalf(
			"failed to get account limit for specific ip. expected: %d, actual: %d",
			config.ACCOUNT_LIMIT_PER_IP,
			limit,
		)
	}
	TestSetAccountLimitForSpecificIp(t)
	if limit := getAccountLimitForSpecificIp("some-ip-address"); limit != 100 {
		t.Fatalf(
			"failed to get account limit for specific ip. expected: %d, actual: %d",
			100,
			limit,
		)
	}
	if limit := getAccountLimitForSpecificIp("another-ip-address"); limit != config.ACCOUNT_LIMIT_PER_IP {
		t.Fatalf(
			"failed to get account limit for another specific ip. expected: %d, actual: %d",
			config.ACCOUNT_LIMIT_PER_IP,
			limit,
		)
	}
}

func TestCheckIpAccountLimit(t *testing.T) {
	db := tests.InitMockDB(true, true)
	db_aggregator.Initialize(db)

	for i := 0; i < int(config.ACCOUNT_LIMIT_PER_IP); i++ {
		if isLimited, err := checkIpAccountLimit("some-ip-address"); err != nil || isLimited {
			t.Fatalf(
				"index: %d, err: %v, should not be limited",
				i, err,
			)
		}
		if result := db.Create(&models.User{
			Name:      fmt.Sprintf("user%d", i),
			IpAddress: "some-ip-address",
		}); result.Error != nil {
			t.Fatalf(
				"failed to create a user record. index: %d, err: %v",
				i, result.Error,
			)
		}
	}

	if isLimited, err := checkIpAccountLimit("some-ip-address"); err != nil || !isLimited {
		t.Fatalf("should be limited after reaching limit: err: %v", err)
	}

	TestSetAccountLimitForSpecificIp(t) // `some-ip-address` has 100 limit here

	for i := int(config.ACCOUNT_LIMIT_PER_IP); i < 100; i++ {
		if isLimited, err := checkIpAccountLimit("some-ip-address"); err != nil || isLimited {
			t.Fatalf(
				"index: %d, err: %v, should not be limited",
				i, err,
			)
		}
		if result := db.Create(&models.User{
			Name:      fmt.Sprintf("user%d", i),
			IpAddress: "some-ip-address",
		}); result.Error != nil {
			t.Fatalf(
				"failed to create a user record. index: %d, err: %v",
				i, result.Error,
			)
		}
	}

	if isLimited, err := checkIpAccountLimit("some-ip-address"); err != nil || !isLimited {
		t.Fatalf("should be limited after reaching limit: err: %v", err)
	}

	if result := db.Where("id > 20").Delete(&models.User{}); result.Error != nil {
		t.Fatalf("failed to remove 80 users: %v", result.Error)
	}

	if isLimited, err := checkIpAccountLimit("some-ip-address"); err != nil || isLimited {
		t.Fatalf(
			"should not be limited after user removal err: %v", err,
		)
	}

	if err := setAccountLimitForSpecificIp("some-ip-address", 10); err != nil {
		t.Fatalf("failed to set account limit: %v", err)
	}
	if isLimited, err := checkIpAccountLimit("some-ip-address"); err != nil || !isLimited {
		t.Fatalf("should be blocked if the number accounts are bigger than limit: err: %v", err)
	}
}
