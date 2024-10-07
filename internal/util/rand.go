package util

import (
	"math/rand"
)

// var r *rand.Rand
//
// func init() {
// 	r = rand.New(rand.NewSource(time.Now().UnixNano()))
// }

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandId() string {
	return RandString(20)
}
