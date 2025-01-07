package membership

import (
	"sync"

	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/paulja/go-work/scheduler/internal/ports"
)

var _ ports.MembershipPort = (*Adapter)(nil)

type Adapter struct {
	sync.Mutex

	store map[string]*domain.Member
}

func NewAdapter() *Adapter {
	return &Adapter{
		store: make(map[string]*domain.Member),
	}
}

func (a *Adapter) AddMember(m *domain.Member) error {
	if m.Id == "" {
		return domain.ErrIdRequired
	}
	if m.Address == "" {
		return domain.ErrAddressRequired
	}

	a.Lock()
	a.store[m.Id] = m
	a.Unlock()
	return nil
}

func (a *Adapter) RemoveMember(id string) error {
	if id == "" {
		return domain.ErrIdRequired
	}
	if _, ok := a.store[id]; !ok {
		return domain.ErrMemberNotFound
	}

	a.Lock()
	delete(a.store, id)
	a.Unlock()

	return nil
}

func (a *Adapter) UpdateMemberStatus(id string, status domain.MembershipStatus) error {
	if id == "" {
		return domain.ErrIdRequired
	}
	if status == domain.MembershipStatusUnknown {
		return domain.ErrStatusRequired
	}
	if _, ok := a.store[id]; !ok {
		return domain.ErrMemberNotFound
	}

	a.Lock()
	a.store[id].SetMembershipStatus(status)
	a.Unlock()

	return nil
}

func (a *Adapter) UpdateHeartbeatStatus(id string, status domain.HeartbeatStatus) error {
	if id == "" {
		return domain.ErrIdRequired
	}
	if status == domain.HeartbeatStatusUnknown {
		return domain.ErrStatusRequired
	}
	if _, ok := a.store[id]; !ok {
		return domain.ErrMemberNotFound
	}

	a.Lock()
	a.store[id].UpdateHeartbeatStatus(status)
	a.Unlock()

	return nil
}

func (a *Adapter) ListMembers() ([]*domain.Member, error) {
	members := make([]*domain.Member, 0)
	for _, m := range a.store {
		members = append(members, m)
	}
	return members, nil
}
