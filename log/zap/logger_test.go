package zap

import "testing"

func TestLogger(t *testing.T) {
	cfg := DefaultConfig()
	logger := New(cfg)
	logger.Info("TestLogger")
}
