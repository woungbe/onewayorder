package service

import (
	"fmt"
	"log"
	"onewayorder/errors"
	"testing"

	"github.com/spf13/viper"
)

func GetBasicInit() *BasicInfo {

	errors.Path = "../logs/error_log"

	client := new(BasicInfo)
	bacc := new(BinanceAccount)

	viper.SetConfigFile("../.env") // .env 파일 설정
	viper.AutomaticEnv()           // 환경 변수를 자동으로 읽도록 설정
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	accKey := viper.GetString("AccessKey")
	seckey := viper.GetString("SecretKey")

	bacc.Init(accKey, seckey)
	client.Init(bacc, "BTCUSDT", "20")
	return client
}

// 최소 거래금액,  코인 소수점
func TestGetExchangeInfo(t *testing.T) {
	basic := GetBasicInit()
	basic.SetSymbol("DOGEUSDT")
	basic.SetLeverage("20")
	err := basic.SetExChangeInfo()
	if err != nil {
		fmt.Println(err)
	}
	// 5 USDT , 0.001
	fmt.Println(err)
}
