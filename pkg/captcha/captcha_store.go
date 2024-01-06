package captcha

import (
	"math/rand"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

type store struct {
	data *cache.Cache
}

func NewStore(data map[string]int) *store {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &store{
		data: c,
	}
}

func (s *store) GenerateCaptchaId() (string, int) {
	captchaId := utils.RandAlphaNum(12)
	captchaValue := rand.Intn(99) + 1
	s.data.Set(captchaId, captchaValue, cache.DefaultExpiration)
	return captchaId, captchaValue

}
