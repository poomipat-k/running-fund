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

// func TestCheckCaptcha(t *testing.T) {
// 	data := make(map[string]float64)

// 	tests := []struct {
// 		name           string
// 		payload        captcha.CheckCaptchaRequest
// 		store          *MockCaptchaStore
// 		expectedStatus int
// 		expectedError  error
// 	}{
// 		{
// 			name: "should fail when captcha is missing",
// 			payload: captcha.CheckCaptchaRequest{
// 				CaptchaId:    "",
// 				CaptchaValue: 20,
// 			},
// 			store: &MockCaptchaStore{
// 				Store: data,
// 				GetFunc: func(captchaId string) (float64, bool) {
// 					return 0, false
// 				},
// 			},
// 			expectedStatus: 400,
// 			expectedError:  &captcha.CaptchaIdRequiredError{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			handler := captcha.NewCaptchaHandler(tt.store)
// 			reqPayload := checkCaptchaPayloadToJSON(tt.payload)
// 			req := httptest.NewRequest(http.MethodPost, "/api/v1/captcha/check", reqPayload)
// 			res := httptest.NewRecorder()

// 			handler.CheckCaptcha(res, req)

// 			assertStatus(t, res.Code, tt.expectedStatus)

// 			t.Log(res.Body.String())
// 			if tt.expectedError != nil {
// 				errBody := getErrorResponse(t, res)
// 				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
// 			}

// 		})
// 	}
// }

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

// func getErrorResponse(t testing.TB, res *httptest.ResponseRecorder) ErrorBody {
// 	t.Helper()
// 	var body ErrorBody
// 	err := json.Unmarshal(res.Body.Bytes(), &body)
// 	if err != nil {
// 		t.Errorf("Error unmarshal ErrorResponse")
// 	}
// 	return body
// }

// func assertErrorMessage(t testing.TB, got, want string) {
// 	t.Helper()
// 	if got != want {
// 		t.Errorf("did not get correct error, got %v, want %v", got, want)
// 	}
// }

// func checkCaptchaPayloadToJSON(payload captcha.CheckCaptchaRequest) *strings.Reader {
// 	userJson, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return strings.NewReader(string(userJson))
// }
