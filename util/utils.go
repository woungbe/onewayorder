package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

func JsonData(aa interface{}) string {
	jsonData, err := json.MarshalIndent(aa, "", "    ")
	if err != nil {
		fmt.Printf("Error encoding JSON: %s\n", err.Error())
		return ""
	}
	return string(jsonData)
}

type APIError struct {
	Code    string
	Message string
}

func ParseError(err error) (*APIError, error) {
	errMsg := err.Error()

	// 에러 메시지에서 'code'와 'msg' 부분을 분리합니다.
	parts := strings.Split(errMsg, ", ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid error format")
	}

	// 'code=' 부분과 'msg=' 부분을 파싱합니다.
	codePart := strings.TrimPrefix(parts[0], "<APIError> code=")
	msgPart := strings.TrimPrefix(parts[1], "msg=")

	return &APIError{
		Code:    codePart,
		Message: msgPart,
	}, nil
}
