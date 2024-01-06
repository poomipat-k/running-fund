package captcha_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/captcha"
)

type MockCaptchaStore struct {
	Store                 map[string]int
	GenerateCaptchaIdFunc func() (string, int)
}

func (m *MockCaptchaStore) GenerateCaptchaId() (string, int) {
	return m.GenerateCaptchaIdFunc()
}

func TestGenerateCaptcha(t *testing.T) {
	data := make(map[string]int)
	store := &MockCaptchaStore{
		Store: data,
		GenerateCaptchaIdFunc: func() (string, int) {
			captchaId := "RhEzlSh46ClI"
			captchaValue := 47
			data[captchaId] = captchaValue
			return captchaId, captchaValue
		},
	}
	handler := captcha.NewCaptchaHandler(store)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/captcha/generate", nil)
	res := httptest.NewRecorder()

	handler.GenerateCaptcha(res, req)
	assertStatus(t, res.Code, http.StatusOK)

	var got struct{ CaptchaId string }
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

	wantValue := 47
	if store.Store[wantCaptchaId] != wantValue {
		t.Errorf("captchaValue got %d, want %d", store.Store[wantCaptchaId], wantValue)
	}

}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
