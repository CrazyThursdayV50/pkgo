package sugar

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := New(DefaultConfig())
	now := time.Now()

	logger.Infof("[%d]time: %s, num: %f", now.Unix(), now.String(), 10.3)
}
