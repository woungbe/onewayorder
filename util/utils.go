package util

import (
	"fmt"
	"math/rand"
	"strconv"
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

func FormatPrice(priceStr string, decimalStr string) (string, error) {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return "", fmt.Errorf("invalid price: %v", err)
	}

	// Convert decimal string to float64
	decimal, err := strconv.ParseFloat(decimalStr, 64)
	if err != nil {
		return "", fmt.Errorf("invalid decimal places: %v", err)
	}

	// Calculate the number of decimal places
	decimalPlaces := 0
	for decimal < 1 {
		decimal *= 10
		decimalPlaces++
	}

	// Format price with specified decimal places
	format := fmt.Sprintf("%%.%df", decimalPlaces)
	formattedPrice := fmt.Sprintf(format, price)

	return formattedPrice, nil
}
