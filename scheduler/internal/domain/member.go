package domain

import (
	"errors"
	"strings"
	"sync"
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
	sync.Mutex
	membershipStatus MembershipStatus
	heartbeatStatus  HeartbeatStatus

	Id      string
	Address string

	timeoutFunc func()
}

func NewMember(id, address string) *Member {
	return &Member{
		Id:      id,
		Address: address,
	}
}

func (m *Member) MembershipStatus() MembershipStatus {
	m.Lock()
	defer m.Unlock()

	return m.membershipStatus
}

func (m *Member) SetMembershipStatus(s MembershipStatus) {
	m.Lock()
	defer m.Unlock()

	m.membershipStatus = s
}

func (m *Member) HeartbeatStatus() HeartbeatStatus {
	m.Lock()
	defer m.Unlock()

	return m.heartbeatStatus
}

func (m *Member) SetHeartbeatStatus(s HeartbeatStatus) {
	m.Lock()
	defer m.Unlock()

	m.heartbeatStatus = s
}

func (m *Member) UpdateHeartbeatStatus(status HeartbeatStatus) {
	m.SetMembershipStatus(MembershipStatusAlive)
	m.SetHeartbeatStatus(status)

	m.timeoutFunc = nil
	m.timeoutFunc = func() {
		<-time.After(config.GetHeartbeatTimeout() * time.Second)
		m.SetMembershipStatus(MembershipStatusLeft)
		m.SetHeartbeatStatus(HeartbeatStatusUnknown)
	}
	go m.timeoutFunc()
}

func (m *Member) StatusString() string {
	status := strings.Builder{}

	switch m.MembershipStatus() {
	case MembershipStatusAlive:
		status.WriteString("alive")
	case MembershipStatusLeft:
		status.WriteString("left")
	default:
		status.WriteString("unknown")
	}

	status.WriteString(", ")

	switch m.HeartbeatStatus() {
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
