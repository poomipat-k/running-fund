package utils

import "math/rand"

const alphaNumericBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandAlphaNum(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphaNumericBytes[rand.Int63()%int64(len(alphaNumericBytes))]
	}
	return string(b)
}
