package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/paulja/go-work/worker/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("can get leader addr default", func(t *testing.T) {
		assert.Equal(t, "localhost:50051", config.GetLeaderAddr(), "should be able to get default")
	})
	t.Run("can override leader addr", func(t *testing.T) {
		os.Setenv("LEADER_ADDR", "test:9000")
		assert.Equal(t, "test:9000", config.GetLeaderAddr(), "should be able to override value")
		os.Unsetenv("LEADER_ADDR")
	})
	t.Run("can get scheduler addr default", func(t *testing.T) {
		assert.Equal(t, "localhost:50052", config.GetSchedulerAddr(), "should be able to get default")
	})
	t.Run("can override scheduler addr", func(t *testing.T) {
		os.Setenv("SCHEDULER_ADDR", "test:9001")
		assert.Equal(t, "test:9001", config.GetSchedulerAddr(), "should be able to override value")
		os.Unsetenv("SCHEDULER_ADDR")
	})
	t.Run("can get worker name default", func(t *testing.T) {
		assert.Regexp(t, "^WORKER_\\d+?", config.GetName(), "should be able to get default")
	})
	t.Run("can override worker name", func(t *testing.T) {
		os.Setenv("WORKER_NAME", "test_01")
		assert.Equal(t, "test_01", config.GetName(), "should be able to override value")
		os.Unsetenv("WORKER_NAME")
	})
	t.Run("can get same worker name default", func(t *testing.T) {
		name := config.GetName()
		assert.Regexp(t, "^WORKER_\\d+?", name, "should be able to get default")
		assert.Equal(t, name, config.GetName(), "should match")
		os.Unsetenv("WORKER_NAME")
	})
	t.Run("can get addr default", func(t *testing.T) {
		assert.Regexp(t, "^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}:\\d{4,5}$", config.GetAddr(), "should be able to get addr")
	})
	t.Run("can get worker port default", func(t *testing.T) {
		assert.Equal(t, 40041, config.GetWorkerPort(), "should be able to get default")
	})
	t.Run("can override worker port", func(t *testing.T) {
		os.Setenv("WORKER_PORT", "9000")
		assert.Equal(t, 9000, config.GetWorkerPort(), "should be able to override value")
		os.Unsetenv("WORKER_PORT")
	})
	t.Run("can find local addr default", func(t *testing.T) {
		assert.Regexp(t, "^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$", config.GetLocalAddr(), "should be able to get default")
	})
	t.Run("can override local addr", func(t *testing.T) {
		os.Setenv("LOCAL_ADDR", "1.2.3.4")
		assert.Equal(t, "1.2.3.4", config.GetLocalAddr())
		os.Unsetenv("LOCAL_ADDR")
	})
	t.Run("can get heartbeat timeout default", func(t *testing.T) {
		assert.Equal(t, time.Duration(15), config.GetHeartbeatTimeout(), "should be able to get default")
	})
	t.Run("can override heartbeat timeout", func(t *testing.T) {
		os.Setenv("HEARTBEAT_TIMEOUT", "99")
		assert.Equal(t, time.Duration(99), config.GetHeartbeatTimeout(), "should be able to override value")
		os.Unsetenv("HEARTBEAT_TIMEOUT")
	})
	t.Run("can get environment default", func(t *testing.T) {
		assert.Equal(t, "development", config.GetEnvironment(), "should be able to get default")
	})
	t.Run("can override environment", func(t *testing.T) {
		os.Setenv("ENV", "testing")
		assert.Equal(t, "testing", config.GetEnvironment(), "should be able to override value")
		os.Unsetenv("ENV")
	})
}
