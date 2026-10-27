package testutil

import (
	"aws-ses-local-go/domain"
	"sync"
)

// FakeMailRepo は domain.IMailRepo のインメモリ実装。Err を設定すると全メソッドがそのエラーを返す。
type FakeMailRepo struct {
	mu         sync.Mutex
	Stored     []domain.Mail
	DeletedIDs []string
	DeleteAllN int
	Err        error
}

var _ domain.IMailRepo = (*FakeMailRepo)(nil)

func (r *FakeMailRepo) Store(m domain.Mail) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return r.Err
	}
	r.Stored = append(r.Stored, m)
	return nil
}

func (r *FakeMailRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return r.Err
	}
	r.DeletedIDs = append(r.DeletedIDs, id)
	return nil
}

func (r *FakeMailRepo) DeleteAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return r.Err
	}
	r.DeleteAllN++
	return nil
}
