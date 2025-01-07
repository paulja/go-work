package domain_test

import (
	"os"
	"testing"
	"time"

	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDomain(t *testing.T) {
	t.Run("can create member", func(t *testing.T) {
		m := domain.NewMember("1", "localhost")
		assert.NotNil(t, m, "should be able to create members")
		assert.Equal(t, "1", m.Id, "should have set value")
		assert.Equal(t, "localhost", m.Address, "should have set value")
	})
	t.Run("can update heartbeat status", func(t *testing.T) {
		m := domain.NewMember("1", "localhost")
		assert.Equal(t, domain.HeartbeatStatusUnknown, m.HeartbeatStatus())
		m.UpdateHeartbeatStatus(domain.HeartbeatStatusBusy)
		assert.Equal(t, domain.HeartbeatStatusBusy, m.HeartbeatStatus())
	})
	t.Run("update heartbeat status expires", func(t *testing.T) {
		m := domain.NewMember("1", "locahost")
		os.Setenv("HEARTBEAT_TIMEOUT", "1")

		assert.Equal(t, domain.HeartbeatStatusUnknown, m.HeartbeatStatus())
		assert.Equal(t, domain.MembershipStatusUnknown, m.MembershipStatus())
		m.UpdateHeartbeatStatus(domain.HeartbeatStatusBusy)
		assert.Equal(t, domain.HeartbeatStatusBusy, m.HeartbeatStatus())
		assert.Equal(t, domain.MembershipStatusAlive, m.MembershipStatus())

		time.Sleep(1200 * time.Millisecond)

		assert.Equal(t, domain.HeartbeatStatusUnknown, m.HeartbeatStatus())
		assert.Equal(t, domain.MembershipStatusLeft, m.MembershipStatus())
	})
	t.Run("can get status string", func(t *testing.T) {
		m := domain.NewMember("1", "localhost")
		os.Setenv("HEARTBEAT_TIMEOUT", "1")

		assert.Equal(t, "unknown, unknown", m.StatusString(), "should be able to get string")
		m.UpdateHeartbeatStatus(domain.HeartbeatStatusBusy)
		assert.Equal(t, domain.MembershipStatusAlive, m.MembershipStatus())
		assert.Equal(t, "alive, busy", m.StatusString(), "should be able to get string")
		time.Sleep(1 * time.Second)
		assert.Equal(t, "left, unknown", m.StatusString(), "should be able to get string")
	})
}
