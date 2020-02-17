package util

import (
	"math/rand"
	"time"
)

//随机名字
func GetRandomName() string {
	prefixStr := "abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(prefixStr)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 3; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	str := "0123456789" + prefixStr
	bytes = []byte(str)
	for i := 0; i < 3; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
