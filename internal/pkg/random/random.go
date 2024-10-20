package random

import (
	"math/rand"
)

// var r *rand.Rand
//
// func init() {
// 	r = rand.New(rand.NewSource(time.Now().UnixNano()))
// }

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateID() string {
	return RandString(20)
}
