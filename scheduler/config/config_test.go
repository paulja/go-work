package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/paulja/go-work/scheduler/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("can get leader port default", func(t *testing.T) {
		assert.Equal(t, 50051, config.GetLeaderPort(), "should be able to get default")
	})
	t.Run("can override leader addr", func(t *testing.T) {
		os.Setenv("LEADER_PORT", "9000")
		assert.Equal(t, 9000, config.GetLeaderPort(), "should be able to override value")
		os.Unsetenv("LEADER_PORT")
	})
	t.Run("can get rpc port default", func(t *testing.T) {
		assert.Equal(t, 50052, config.GetRPCPort(), "should be able to get default")
	})
	t.Run("can override rpc port", func(t *testing.T) {
		os.Setenv("RPC_PORT", "9000")
		assert.Equal(t, 9000, config.GetRPCPort(), "should be able to override value")
		os.Unsetenv("RPC_PORT")
	})
	t.Run("can get heartbeat timeout default", func(t *testing.T) {
		assert.Equal(t, time.Duration(30), config.GetHeartbeatTimeout(), "should be able to get default")
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
