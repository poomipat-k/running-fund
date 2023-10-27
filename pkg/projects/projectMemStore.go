package projects

import (
	"time"
)

type MemStore struct {
	list []Project
}

func NewMemStore() *MemStore {
	list := []Project{}
	p1 := Project{
		Id:             1,
		ProjectCode:    "a",
		ProjectVersion: 1,
		CreatedAt:      time.Now(),
	}
	list = append(list, p1)
	return &MemStore{
		list: list,
	}
}

func (m *MemStore) GetAll() ([]Project, error) {
	return m.list, nil
}
