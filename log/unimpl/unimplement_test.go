package unimpl

import "testing"

func TestUnimplement(t *testing.T) {
	logger := New()
	logger.Debug(1, 2, 3)
	logger.Info(1, 2, 3)
	logger.Warn(1, 2, 3)
	logger.Error(1, 2, 3)
	logger.Debugf("debugf: %d, %d, %d\n", 1, 2, 3)
	logger.Infof("debugf: %d, %d, %d\n", 1, 2, 3)
	logger.Warnf("debugf: %d, %d, %d\n", 1, 2, 3)
	logger.Errorf("debugf: %d, %d, %d\n", 1, 2, 3)
}
