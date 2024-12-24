package grpc_test

import (
	"os"
	"testing"
	"time"

	"github.com/paulja/go-work/worker/internal/adapters/grpc"
	"github.com/stretchr/testify/assert"
)

func TestHeartbeat(t *testing.T) {
	leader := NewLeaderMock(t)
	assert.NoError(t, leader.Start())

	t.Run("can start and stop heartbeat", func(t *testing.T) {
		leader.Reset()

		hb := grpc.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")
		assert.NoError(t, hb.Stop(), "should be able to stop heartbeat")
		assert.Equal(t, 1, leader.JoinCallCount, "unexpected join call count")
		assert.Equal(t, 1, leader.LeaveCallCount, "unexpected leave call count")
	})
	t.Run("can apply status", func(t *testing.T) {
		hb := grpc.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")
		hb.ApplyStatus(grpc.HeartbeatStatusBusy)
		assert.NoError(t, hb.Stop(), "should be able to stop heartbeat")
	})
	t.Run("cannot apply invalid status", func(t *testing.T) {
		leader.Reset()

		os.Setenv("HEARTBEAT_TIMEOUT", "1")

		hb := grpc.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")

		hb.ApplyStatus(9)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 0, leader.HeartbeatCallCount, "unexpected heartbeat call count")
		assert.Equal(t, grpc.HeartbeatStatusUnknown, leader.Status, "unexpected heartbeat status")
	})
	t.Run("heartbeat handler sends correct status", func(t *testing.T) {
		leader.Reset()

		os.Setenv("HEARTBEAT_TIMEOUT", "1")

		hb := grpc.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")

		hb.ApplyStatus(grpc.HeartbeatStatusIdle)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 1, leader.HeartbeatCallCount, "unexpected heartbeat call count")
		assert.Equal(t, grpc.HeartbeatStatusIdle, leader.Status, "unexpected heartbeat status")

		hb.ApplyStatus(grpc.HeartbeatStatusBusy)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 2, leader.HeartbeatCallCount, "unexpected heartbeat call count")
		assert.Equal(t, grpc.HeartbeatStatusBusy, leader.Status, "unexpected heartbeat status")

		hb.ApplyStatus(grpc.HeartbeatStatusFailed)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 3, leader.HeartbeatCallCount, "unexpected heartbeat call count")
		assert.Equal(t, grpc.HeartbeatStatusFailed, leader.Status, "unexpected heartbeat status")

		assert.NoError(t, hb.Stop(), "should be able to stop heartbeat")
	})

	assert.NoError(t, leader.Stop())
}
