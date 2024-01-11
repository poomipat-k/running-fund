package captcha

import (
	"encoding/base64"
	"math/rand"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

var puzzles = []Puzzle{
	{
		BackgroundPath: "./tmp/desert_70.png",
		PuzzlePath:     "./tmp/desert_jigsaw_70.png",
		Value:          70,
	},
	{
		BackgroundPath: "./tmp/desert_95.png",
		PuzzlePath:     "./tmp/desert_jigsaw_95.png",
		Value:          95,
	},
	{
		BackgroundPath: "./tmp/desert_130.png",
		PuzzlePath:     "./tmp/desert_jigsaw_130.png",
		Value:          130,
	},
}

type store struct {
	data *cache.Cache
}

func NewStore(data map[string]int) *store {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &store{
		data: c,
	}
}

func (s *store) GenerateCaptcha() (Captcha, error) {
	captchaId := utils.RandAlphaNum(12)
	index := rand.Intn(len(puzzles))
	p := puzzles[index]

	base64Background, err := getBase64FromImage(p.BackgroundPath)
	if err != nil {
		return Captcha{}, err
	}

	base64Puzzle, err := getBase64FromImage(p.PuzzlePath)
	if err != nil {
		return Captcha{}, err
	}

	captcha := Captcha{
		CaptchaId:    captchaId,
		Background64: base64Background,
		Puzzle64:     base64Puzzle,
	}
	s.data.Set(captchaId, p.Value, cache.DefaultExpiration)

	return captcha, nil
}

func (s *store) Get(captchaId string) (int, bool) {
	raw, found := s.data.Get(captchaId)
	v, ok := raw.(int)
	if !ok {
		return 0, false
	}
	return v, found
}

func getBase64FromImage(filepath string) (string, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	base64Data := "data:image/png;base64,"
	base64Data += toBase64(bytes)
	return base64Data, nil
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
