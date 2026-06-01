package random

import (
	"math/rand"
	"time"
)

const tokenLength = 6

type base62Generator struct {
	length int
}

func NewBase62Generator() *base62Generator {
	return &base62Generator{length: tokenLength}
}

// NewRandomString generates random string with given size.
func (g *base62Generator) RandomString() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, g.length)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
