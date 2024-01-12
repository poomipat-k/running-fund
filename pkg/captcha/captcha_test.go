package captcha_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/captcha"
)

type MockCaptchaStore struct {
	Store               map[string]float64
	GenerateCaptchaFunc func() (captcha.Captcha, error)
	GetFunc             func(captchaId string) (float64, bool)
}

func (m *MockCaptchaStore) GenerateCaptcha() (captcha.Captcha, error) {
	return m.GenerateCaptchaFunc()
}

func (m *MockCaptchaStore) Get(captchaId string) (float64, bool) {
	return m.GetFunc(captchaId)
}

type ErrorBody struct {
	Error   bool
	Message string
	Name    string
}

func TestGenerateCaptcha(t *testing.T) {
	data := make(map[string]float64)
	store := &MockCaptchaStore{
		Store: data,
		GenerateCaptchaFunc: func() (captcha.Captcha, error) {
			captchaId := "RhEzlSh46ClI"
			captchaValue := 47.2
			data[captchaId] = captchaValue
			return captcha.Captcha{
				CaptchaId:    captchaId,
				Background64: "abc",
				Puzzle64:     "def",
			}, nil
		},
	}
	handler := captcha.NewCaptchaHandler(store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/captcha/generate", nil)
	res := httptest.NewRecorder()

	handler.GenerateCaptcha(res, req)
	assertStatus(t, res.Code, http.StatusOK)

	var got captcha.Captcha
	err := json.Unmarshal(res.Body.Bytes(), &got)
	if err != nil {
		t.Errorf("fail to unmarshal err: %+v", err)
	}

	if len(got.CaptchaId) != 12 {
		t.Errorf("CaptchaId length is not valid got %d, want %d", len(got.CaptchaId), 24)
	}

	wantCaptchaId := "RhEzlSh46ClI"
	if got.CaptchaId != wantCaptchaId {
		t.Errorf("CaptchaId is not match got %s, want %s", got.CaptchaId, wantCaptchaId)
	}

	wantValue := 47.2
	if store.Store[wantCaptchaId] != wantValue {
		t.Errorf("captchaValue got %f, want %f", store.Store[wantCaptchaId], wantValue)
	}

	wantBase64Background := "abc"
	if got.Background64 != wantBase64Background {
		t.Errorf("base64 got %s, want %s", got.Background64, wantBase64Background)
	}

	wantBase64Puzzle := "def"
	if got.Puzzle64 != wantBase64Puzzle {
		t.Errorf("base64 got %s, want %s", got.Background64, wantBase64Puzzle)
	}
}

func TestCheckCaptcha(t *testing.T) {
	store := &MockCaptchaStore{}

	// Payload validation cases
	tests := []struct {
		name          string
		captchaId     string
		captchaValue  float64
		expectedError error
	}{
		{
			name:          "should error when captchaId is missing",
			captchaId:     "",
			captchaValue:  40,
			expectedError: &captcha.CaptchaIdRequiredError{},
		},
		{
			name:          "should error when captchaValue is missing",
			captchaId:     "abc",
			expectedError: &captcha.CaptchaValueRequiredError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := captcha.CheckCaptcha(tt.captchaId, tt.captchaValue, store)
			if tt.expectedError != nil {
				assertErrorMessage(t, err.Error(), tt.expectedError.Error())
			}
		})
	}

	// Mock data
	data := make(map[string]float64)
	store = &MockCaptchaStore{
		Store: data,
		GenerateCaptchaFunc: func() (captcha.Captcha, error) {
			captchaId := "RhEzlSh46ClI"
			captchaValue := 47.2
			data[captchaId] = captchaValue
			return captcha.Captcha{
				CaptchaId:    captchaId,
				Background64: "abc",
				Puzzle64:     "def",
			}, nil
		},
		GetFunc: func(captchaId string) (float64, bool) {
			v, ok := data[captchaId]
			return v, ok
		},
	}

	t.Run("should error when captchaId is not found", func(t *testing.T) {
		_, _ = store.GenerateCaptchaFunc()
		id := "RhEzlSh_TEST"
		value := 46.5
		_, err := captcha.CheckCaptcha(id, value, store)
		want := &captcha.CaptchaNotFoundError{}
		assertErrorMessage(t, err.Error(), want.Error())
	})

	t.Run("should error when captchaValue is out of accept limit", func(t *testing.T) {
		_, _ = store.GenerateCaptchaFunc()
		id := "RhEzlSh46ClI"
		value := 12.5
		_, err := captcha.CheckCaptcha(id, value, store)
		want := &captcha.CaptchaValueNotValidError{}
		assertErrorMessage(t, err.Error(), want.Error())
	})

	t.Run("should be ok", func(t *testing.T) {
		_, _ = store.GenerateCaptchaFunc()
		id := "RhEzlSh46ClI"
		value := 42.5
		_, err := captcha.CheckCaptcha(id, value, store)
		if err != nil {
			t.Errorf("should be ok but got error: %v", err)
		}
	})
}

func assertErrorMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct error, got %v, want %v", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
