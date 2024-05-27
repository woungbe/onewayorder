package utils

import (
	"math/rand"
	"time"
)

var randSeed rand.Source

// RandInitSeed 랜덤 씨드 초기화
func RandInitSeed() {
	randSeed = rand.NewSource(time.Now().UnixNano())

}
func GetRandSeed() rand.Source {
	return randSeed
}

// GetCurrentTimestamp UTC 현재 시간 리턴
func GetCurrentTimestamp() int64 {
	return time.Now().UTC().Unix()
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
