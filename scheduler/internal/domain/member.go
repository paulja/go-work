package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/paulja/go-work/scheduler/config"
)

var (
	ErrIdRequired      = errors.New("id required")
	ErrAddressRequired = errors.New("address required")
	ErrStatusRequired  = errors.New("status required")
	ErrMemberNotFound  = errors.New("member not found")
)

type MembershipStatus int

const (
	MembershipStatusUnknown MembershipStatus = iota
	MembershipStatusAlive
	MembershipStatusLeft
)

type HeartbeatStatus int

const (
	HeartbeatStatusUnknown HeartbeatStatus = iota
	HeartbeatStatusIdle
	HeartbeatStatusBusy
	HeartbeatStatusFailed
)

type Member struct {
	Id               string
	Address          string
	MembershipStatus MembershipStatus
	HeartbeatStatus  HeartbeatStatus

	timeoutFunc func()
}

func NewMember(id, address string) *Member {
	return &Member{
		Id:      id,
		Address: address,
	}
}

func (m *Member) UpdateHeartbeatStatus(status HeartbeatStatus) {
	m.MembershipStatus = MembershipStatusAlive
	m.HeartbeatStatus = status

	m.timeoutFunc = nil
	m.timeoutFunc = func() {
		<-time.After(config.GetHeartbeatTimeout() * time.Second)
		m.MembershipStatus = MembershipStatusLeft
		m.HeartbeatStatus = HeartbeatStatusUnknown
	}
	go m.timeoutFunc()
}

func (m *Member) StatusSting() string {
	status := strings.Builder{}

	switch m.MembershipStatus {
	case MembershipStatusAlive:
		status.WriteString("alive")
	case MembershipStatusLeft:
		status.WriteString("left")
	default:
		status.WriteString("unknown")
	}

	status.WriteString(", ")

	switch m.HeartbeatStatus {
	case HeartbeatStatusIdle:
		status.WriteString("idle")
	case HeartbeatStatusBusy:
		status.WriteString("busy")
	case HeartbeatStatusFailed:
		status.WriteString("failed")
	case HeartbeatStatusUnknown:
		status.WriteString("unknown")
	}

	return status.String()
}
