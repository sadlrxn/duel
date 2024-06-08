package grand_jackpot

import (
	"testing"
	"time"
)

func TestSetGlobalTimeZone(t *testing.T) {
	time.Local = time.UTC
	t.Log(time.Now().String())
}
