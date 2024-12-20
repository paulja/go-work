package ports

import "github.com/paulja/go-work/scheduler/internal/domain"

type MembershipPort interface {
	AddMember(*domain.Member) error
	RemoveMember(id string) error
	UpdateMemberStatus(id string, status domain.MembershipStatus) error
	UpdateHeartbeatStatus(id string, status domain.HeartbeatStatus) error
	ListMembers() ([]*domain.Member, error)
}
