package membership_test

import (
	"testing"

	"github.com/paulja/go-work/scheduler/internal/adapters/membership"
	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMembership(t *testing.T) {
	t.Run("list: can list members", func(t *testing.T) {
		a := membership.NewAdapter()
		list, err := a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Len(t, list, 0, "list should be empty")
	})
	t.Run("add: can add member", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.AddMember(domain.NewMember("1", "1.2.3.4"))
		assert.NoError(t, err, "should be able to add member")
		list, err := a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Len(t, list, 1, "list should have one member")
	})
	t.Run("add: should error if member is invalid", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.AddMember(domain.NewMember("", ""))
		assert.ErrorIs(t, err, domain.ErrIdRequired, "should return id required error")
		err = a.AddMember(domain.NewMember("1", ""))
		assert.ErrorIs(t, err, domain.ErrAddressRequired, "should return address required error")
	})
	t.Run("remove: can remove member", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.AddMember(domain.NewMember("1", "1.2.3.4"))
		assert.NoError(t, err, "should be able to add member")
		list, err := a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Len(t, list, 1, "list should have one member")
		err = a.RemoveMember("1")
		assert.NoError(t, err, "should be able to remove member")
		list, err = a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Len(t, list, 0, "list should be empty")
	})
	t.Run("remove: should error if input is invalid", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.RemoveMember("")
		assert.ErrorIs(t, err, domain.ErrIdRequired, "should return id required error")
		err = a.RemoveMember("1")
		assert.ErrorIs(t, err, domain.ErrMemberNotFound, "should return member not found error")
	})
	t.Run("member: can update status", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.AddMember(domain.NewMember("1", "1.2.3.4"))
		list, err := a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Equal(t, domain.MembershipStatusUnknown, list[0].MembershipStatus(), "status should be unknown")
		err = a.UpdateMemberStatus("1", domain.MembershipStatusAlive)
		assert.NoError(t, err, "should be able to update status")
		list, err = a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Equal(t, domain.MembershipStatusAlive, list[0].MembershipStatus(), "status should be alive")
	})
	t.Run("member: should error if input is invalid", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.UpdateMemberStatus("", 0)
		assert.ErrorIs(t, err, domain.ErrIdRequired, "should return id required error")
		err = a.UpdateMemberStatus("1", 0)
		assert.ErrorIs(t, err, domain.ErrStatusRequired, "should return status required error")
		err = a.UpdateMemberStatus("1", domain.MembershipStatusLeft)
		assert.ErrorIs(t, err, domain.ErrMemberNotFound, "should return member not found error")
	})
	t.Run("heartbeat: can update status", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.AddMember(domain.NewMember("1", "1.2.3.4"))
		list, err := a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Equal(t, domain.HeartbeatStatusUnknown, list[0].HeartbeatStatus(), "status should be unknown")
		err = a.UpdateHeartbeatStatus("1", domain.HeartbeatStatusBusy)
		assert.NoError(t, err, "should be able to update status")
		list, err = a.ListMembers()
		assert.NoError(t, err, "should be able to list members")
		assert.Equal(t, domain.HeartbeatStatusBusy, list[0].HeartbeatStatus(), "status should be busy")
	})
	t.Run("heartbeat: should error if input is invalid", func(t *testing.T) {
		a := membership.NewAdapter()
		err := a.UpdateHeartbeatStatus("", 0)
		assert.ErrorIs(t, err, domain.ErrIdRequired, "should return id required error")
		err = a.UpdateHeartbeatStatus("1", 0)
		assert.ErrorIs(t, err, domain.ErrStatusRequired, "should return status required error")
		err = a.UpdateHeartbeatStatus("1", domain.HeartbeatStatusBusy)
		assert.ErrorIs(t, err, domain.ErrMemberNotFound, "should return member not found error")
	})
}
