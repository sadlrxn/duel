package redis

import (
	"fmt"
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
)

func TestRecentlyWagered(t *testing.T) {
	if err := InitializeMockRedis(true); err != nil {
		t.Fatalf("failed to initialize mock redis: %v", err)
	}

	for i := 1; i < 300; i++ {
		fmt.Println("adding for userID: ", i)
		if err := ZAddRecentlyWagered(uint(i)); err != nil {
			t.Fatalf(
				"failed to perform zadd: index: %d, err: %v",
				i,
				err,
			)
		}
	}

	result := ZRevRangeRecentlyWagered()
	if len(result) != int(config.RAIN_MAX_SPLIT_COUNT) {
		t.Fatalf(
			"failed to retrieve recently wagered properly: %d, %d",
			len(result), config.RAIN_MAX_SPLIT_COUNT,
		)
	}

	for i, item := range result {
		if item != uint(299-i) {
			t.Fatalf(
				"failed to retrieve properly: %d, %d",
				item, 299-i,
			)
		}
	}
}
