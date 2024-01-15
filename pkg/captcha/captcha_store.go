package captcha

import (
	"encoding/base64"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const resourcesBasePath = "./home"

var puzzles = []Puzzle{
	// Value and YPosition are values from photoshop used to generate the background and puzzle
	{
		BackgroundPath: filepath.Join(resourcesBasePath, "captcha/desert_80_80.png"),
		PuzzlePath:     filepath.Join(resourcesBasePath, "captcha/desert_puzzle_80_80.png"),
		Value:          80,
		YPosition:      80,
	},
	{
		BackgroundPath: filepath.Join(resourcesBasePath, "captcha/desert_100_50.png"),
		PuzzlePath:     filepath.Join(resourcesBasePath, "captcha/desert_puzzle_100_50.png"),
		Value:          100,
		YPosition:      50,
	},
	{
		BackgroundPath: filepath.Join(resourcesBasePath, "captcha/desert_200_30.png"),
		PuzzlePath:     filepath.Join(resourcesBasePath, "captcha/desert_puzzle_200_30.png"),
		Value:          200,
		YPosition:      30,
	},
	{
		BackgroundPath: filepath.Join(resourcesBasePath, "captcha/rocky_60_80.png"),
		PuzzlePath:     filepath.Join(resourcesBasePath, "captcha/rocky_puzzle_60_80.png"),
		Value:          60,
		YPosition:      80,
	},
	{
		BackgroundPath: filepath.Join(resourcesBasePath, "captcha/rocky_120_60.png"),
		PuzzlePath:     filepath.Join(resourcesBasePath, "captcha/rocky_puzzle_120_60.png"),
		Value:          120,
		YPosition:      60,
	},
	{
		BackgroundPath: filepath.Join(resourcesBasePath, "captcha/rocky_170_40.png"),
		PuzzlePath:     filepath.Join(resourcesBasePath, "captcha/rocky_puzzle_170_40.png"),
		Value:          170,
		YPosition:      40,
	},
}

type store struct {
	data *cache.Cache
}

func NewStore() *store {
	c := cache.New(3*time.Minute, 5*time.Minute)
	return &store{
		data: c,
	}
}

func (s *store) GenerateCaptcha() (Captcha, error) {
	captchaId := utils.RandAlphaNum(12)

	// Create a new random number generator with a custom seed (e.g., current time)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	index := rng.Intn(len(puzzles))

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
		YPosition:    p.YPosition,
	}
	s.data.Set(captchaId, p.Value, cache.DefaultExpiration)

	return captcha, nil
}

func (s *store) Get(captchaId string) (float64, bool) {
	raw, found := s.data.Get(captchaId)
	v, ok := raw.(float64)
	if !ok {
		return 0, false
	}
	return v, found
}

func (s *store) Delete(captchaId string) {
	s.data.Delete(captchaId)
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
