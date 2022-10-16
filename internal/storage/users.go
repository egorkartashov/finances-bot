package storage

import (
	"sync"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/users"
)

type Users struct {
	mu    *sync.Mutex
	cache map[int64]users.User
}

func NewUsers() *Users {
	return &Users{
		mu:    new(sync.Mutex),
		cache: make(map[int64]users.User),
	}
}

func (s *Users) Save(user users.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[user.GetID()] = user
	return nil
}

func (s *Users) Get(id int64) (user users.User, ok bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok = s.cache[id]
	err = nil
	return
}
