package dreamtower

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
)

func TestCalculateOdd(t *testing.T) {
	expectedOdds := []float32{1, 1.29, 1.67, 2.15, 2.78, 3.58, 4.6, 5.91, 7.59, 9.72}

	for i := 0; i <= 9; i++ {
		actualOdd := calculateMutiplierV2(
			config.DREAMTOWER_DIFFICULTIES["Easy"],
			uint(3),
			i,
		)
		if expectedOdds[i] != actualOdd {
			t.Fatalf(
				"level: %d, expected: %f, actual: %f",
				i, expectedOdds[i], actualOdd,
			)
		}
	}
}
