package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/poomipat-k/running-fund/pkg/projects"
)

type StubProjectStore struct {
	list []projects.ReviewDashboardRow
}

func (s *StubProjectStore) GetReviewerDashboard(userId int, fromDate time.Time, toDate time.Time) ([]projects.ReviewDashboardRow, error) {
	var matched []projects.ReviewDashboardRow
	for i := 0; i < len(s.list); i++ {
		if !(s.list[i].ProjectCreatedAt.Before(fromDate) || s.list[i].ProjectCreatedAt.After(toDate)) {
			matched = append(matched, s.list[i])
		}
	}
	return matched, nil
}

func TestGetReviewerDashboard(t *testing.T) {
	now := time.Now()
	tmr := time.Now().Add(time.Duration(24) * time.Hour)
	next2days := time.Now().Add(time.Duration(48) * time.Hour)
	next3days := time.Now().Add(time.Duration(72) * time.Hour)
	p1 := projects.ReviewDashboardRow{
		ProjectId:        1,
		ProjectCode:      "A",
		ProjectName:      "Project_A",
		ProjectCreatedAt: &now,
	}
	p2 := projects.ReviewDashboardRow{
		ProjectId:        2,
		ProjectCode:      "B",
		ProjectName:      "Project_B",
		ProjectCreatedAt: &tmr,
	}
	p3 := projects.ReviewDashboardRow{
		ProjectId:        3,
		ProjectCode:      "C",
		ProjectName:      "Project_C",
		ProjectCreatedAt: &next2days,
	}
	store := StubProjectStore{
		list: []projects.ReviewDashboardRow{
			p1,
			p2,
			p3,
		},
	}

	projectHandler := NewProjectHandler(&store)

	tests := []struct {
		name               string
		payload            projects.GetReviewerDashboardRequest
		expectedHTTPStatus int
		expectedResponse   []projects.ReviewDashboardRow
	}{
		{
			name: "Return all 3 projects",
			payload: projects.GetReviewerDashboardRequest{
				FromDate: now,
				ToDate:   next3days,
			},
			expectedHTTPStatus: 200,
			expectedResponse:   []projects.ReviewDashboardRow{p1, p2, p3},
		},
		{
			name: "Return all projects created today and tomorrow",
			payload: projects.GetReviewerDashboardRequest{
				FromDate: now,
				ToDate:   tmr,
			},
			expectedHTTPStatus: 200,
			expectedResponse:   []projects.ReviewDashboardRow{p1, p2},
		},
		{
			name: "Return all projects created tomorrow and the next 2 days",
			payload: projects.GetReviewerDashboardRequest{
				FromDate: tmr,
				ToDate:   next2days,
			},
			expectedHTTPStatus: 200,
			expectedResponse:   []projects.ReviewDashboardRow{p2, p3},
		},
		{
			name: "Return all the project created tomorrow",
			payload: projects.GetReviewerDashboardRequest{
				FromDate: tmr.Add(time.Duration(-1) * time.Hour),
				ToDate:   tmr.Add(time.Duration(1) * time.Hour),
			},
			expectedHTTPStatus: 200,
			expectedResponse:   []projects.ReviewDashboardRow{p2},
		},
		{
			name: "Return none (to  is earlier than today)",
			payload: projects.GetReviewerDashboardRequest{
				FromDate: tmr.Add(time.Duration(-24) * time.Hour),
				ToDate:   tmr.Add(time.Duration(-1) * time.Hour),
			},
			expectedHTTPStatus: 200,
			expectedResponse:   nil,
		},
		{
			name: "Return none (from  is later than the next2days)",
			payload: projects.GetReviewerDashboardRequest{
				FromDate: tmr.Add(time.Duration(-24) * time.Hour),
				ToDate:   tmr.Add(time.Duration(-1) * time.Hour),
			},
			expectedHTTPStatus: 200,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := newGetReviewerDashboardRequest(tt.payload)
			response := httptest.NewRecorder()

			projectHandler.GetReviewerDashboard(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			gotContentType := response.Header().Get("Content-Type")
			wantContentType := "application/json"
			if gotContentType != wantContentType {
				t.Errorf("Wrong content type, got %s, want %s", gotContentType, wantContentType)
			}
			assertContentTypeHeader(t, gotContentType, "application/json")

			var got []projects.ReviewDashboardRow
			err := json.Unmarshal(response.Body.Bytes(), &got)
			if err != nil {
				log.Fatal(err)
			}
			assertSliceResponseBody(t, got, tt.expectedResponse)
		})
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertContentTypeHeader(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Wrong content type, got %s, want %s", got, want)
	}
}

func assertSliceResponseBody(t testing.TB, got, want []projects.ReviewDashboardRow) {
	t.Helper()
	if got == nil && want == nil {
		return
	}
	if got == nil && want != nil {
		t.Errorf("Failed_1: response body is wrong, got %v want %v", got, want)
		return
	}
	if got != nil && want == nil {
		t.Errorf("Failed_2: response body is wrong, got %v want %v", got, want)
		return
	}
	if len(got) != len(want) {
		t.Errorf("Failed_3: response body size is not equal, got has %d, want has %d", len(got), len(want))
		return
	}
	for i := 0; i < len(got); i++ {
		if got[i].ProjectCode != want[i].ProjectCode {
			t.Errorf("Failed_4: response body is wrong, got %v want %v", got[i], want[i])
		}
	}
}

func newGetReviewerDashboardRequest(body projects.GetReviewerDashboardRequest) *http.Request {
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/api/projects/reviewer", bytes.NewBuffer(out))
	if err != nil {
		log.Fatal(err)
	}
	return req
}
