package server

import (
	"testing"
	"time"

	"github.com/poomipat-k/running-fund/pkg/projects"
)

type StubProjectStore struct {
	list []projects.Project
}

func (s *StubProjectStore) GetReviewerDashboard() ([]projects.Project, error) {
	return s.list, nil
}

func TestGetReviewerDashboard(t *testing.T) {
	p1 := projects.Project{
		Id:             1,
		ProjectCode:    "A",
		ProjectName:    "Project_A",
		ProjectVersion: 1,
		CreatedAt:      time.Now(),
	}
	p2 := projects.Project{
		Id:             2,
		ProjectCode:    "B",
		ProjectName:    "Project_B",
		ProjectVersion: 1,
		CreatedAt:      time.Now().Add(time.Duration(24) * time.Hour),
	}
	p3 := projects.Project{
		Id:             3,
		ProjectCode:    "C",
		ProjectName:    "Project_C",
		ProjectVersion: 1,
		CreatedAt:      time.Now().Add(time.Duration(48) * time.Hour),
	}
	store := StubProjectStore{
		list: []projects.Project{
			p1,
			p2,
			p3,
		},
	}
	// server := &Server{}

}
