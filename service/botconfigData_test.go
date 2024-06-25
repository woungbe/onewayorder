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

	viper.SetConfigFile("../.dev.env") // .env 파일 설정
	viper.AutomaticEnv()               // 환경 변수를 자동으로 읽도록 설정
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	accKey := viper.GetString("AccessKey")
	seckey := viper.GetString("SecretKey")

	bacc.Init(accKey, seckey)
	client.Init(bacc, "DOGEUSDT", "20")
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

	// fmt.Sprintf("%+v", basic)
	fmt.Println(basic.GetNotionalPrice()) // 5
	fmt.Println(basic.GetMaxPrice())      // 30
	fmt.Println(basic.GetMinPrice())      // 0.00244
	fmt.Println(basic.GetTickSizePrice()) // 0.000010
	fmt.Println(basic.GetMaxQuantity())   //
	fmt.Println(basic.GetMinQuantity())
	fmt.Println(basic.GetStepSize())
}

// 미체결 리스트 리턴

func TestGetOpenOrder(t *testing.T) {
	basic := GetBasicInit()
	basic.SetSymbol("DOGEUSDT")
	basic.SetLeverage("20")

	res := basic.GetOpenOrder()
	if len(res) == 0 {
		fmt.Println(" 없음 ")
		return
	}

	for _, v := range res {
		fmt.Println(v)
	}

}

func TestGetPositionList(t *testing.T) {
	//	func (ty *BasicInfo) GetPositionList() (map[string]UserPositionInfo, error) {
	basic := GetBasicInit()
	res, err := basic.GetPositionList()
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range res {
		fmt.Sprintf("%+v\n", v)
	}

}

func TestGetCurrentPrice(t *testing.T) {
	//	func (ty *BasicInfo) GetCurrentPrice() (string, error) {
	basic := GetBasicInit()
	res, err := basic.GetCurrentPrice()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)

}

func TestGetEntryPrice(t *testing.T) {
	//	func (ty *BasicInfo) GetEntryPrice(longshort string) (string, error) {
	basic := GetBasicInit()

	res, err := basic.GetEntryPrice("LONG")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)

}

func TestGetCoinQtyForPrice(t *testing.T) {
	//	func (ty *BasicInfo) GetCoinQtyForPrice(price, uset string) (string, error) {
	basic := GetBasicInit()
	res, err := basic.GetCoinQtyForPrice("0.0117", "200")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)

}

func TestGetTPPrice(t *testing.T) {
	//	func (ty *BasicInfo) GetTPPrice(price, positionSide, takePersent string) (string, error) {
	basic := GetBasicInit()
	price := "0.0117"
	positionSide := "LONG"
	takePersent := "0.127"
	res, err := basic.GetTPPrice(price, positionSide, takePersent)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)
}

func TestGetSLPrice(t *testing.T) {
	//	func (ty *BasicInfo) GetSLPrice(price, positionSide, takePersent string) (string, error) {
	basic := GetBasicInit()
	price := "0.0117"
	positionSide := "LONG"
	takePersent := "0.107"
	res, err := basic.GetSLPrice(price, positionSide, takePersent)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)
}

func TestCheckMinTraderPirce(t *testing.T) {
	//	func (ty *BasicInfo) CheckMinTraderPirce(price string, Amount string) (bool, error) {
	basic := GetBasicInit()
	price := "0.11770"
	Amount := "2000"
	b, err := basic.CheckMinTraderPirce(price, Amount)
	if err != nil {
		fmt.Println(err)
	}
	if !b {
		fmt.Println("거래불가 : ", err.Error())
	} else {
		fmt.Println("거래하시면 됩니다. ", b)
	}

}

func TestCheckPriceMinMax(t *testing.T) {
	//	func (ty *BasicInfo) CheckPriceMinMax(price string) (bool, error) {
	basic := GetBasicInit()
	price := "5000" // 도지는 안되것지..
	b, err := basic.CheckPriceMinMax(price)
	if err != nil {
		fmt.Println(err)
	}
	if !b {
		fmt.Println("거래불가 : ", err.Error())
	} else {
		fmt.Println("거래하시면 됩니다. ", b)
	}
}

func TestReturnPriceForSize(t *testing.T) {
	//	func (ty *BasicInfo) ReturnPriceForSize(price string) (string, error) {
	basic := GetBasicInit()
	price := "0.1170123012301230"
	b, err := basic.ReturnPriceForSize(price)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("거래하시면 됩니다. ", b)
}

func TestCheckMinAmount(t *testing.T) {
	//	func (ty *BasicInfo) CheckMinAmount(Amount string) (bool, error) {
	basic := GetBasicInit()
	// Amount := "1000000"
	Amount := "0.001"

	b, err := basic.CheckMinAmount(Amount)
	if err != nil {
		fmt.Println(err)
	}

	if !b {
		fmt.Println("거래불가 : ", err.Error())
	} else {
		fmt.Println("거래하시면 됩니다. ", b)
	}

}

func TestCheckMaxAmount(t *testing.T) {
	//	func (ty *BasicInfo) CheckMaxAmount(Amount string) (bool, error) {
	basic := GetBasicInit()
	Amount := "10000000"
	b, err := basic.CheckMaxAmount(Amount)
	if err != nil {
		fmt.Println(err)
	}

	if !b {
		fmt.Println("거래불가 : ", err.Error())
	} else {
		fmt.Println("거래하시면 됩니다. ", b)
	}
}

// func TestBeforeOrder(t *testing.T) {
// 	//	func (ty *BasicInfo) BeforeOrder(currentPrice, amout string) (bool, error) {
// 	basic := GetBasicInit()
// }

func TestSendOpenOrder(t *testing.T) {
	//	func (ty *BasicInfo) SendOpenOrder(position, side, price, amount string) (*futures.CreateOrderResponse, error) {
	basic := GetBasicInit()
	position := "LONG"
	side := "OPEN"
	price := "0.01017"
	amount := "1000"
	res, err := basic.SendOpenOrder(position, side, price, amount)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Sprintf("%+v\n", res)

}

func TestSendMarketOrder(t *testing.T) {
	//	func (ty *BasicInfo) SendMarketOrder(symbol, position, openclose, price string) (*futures.CreateOrderResponse, error) {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	position := "LONG"
	openclose := "OPEN"
	price := "300"
	res, err := basic.SendMarketOrder(symbol, position, openclose, price)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Sprintf("%+v\n", res)

}

func TestSendTakeProfit(t *testing.T) {
	//	func (ty *BasicInfo) SendTakeProfit(symbol, position, openclose, price string) (*futures.CreateOrderResponse, error) {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	position := "LONG"
	openclose := "CLOSE"
	price := "0.1270"
	res, err := basic.SendMarketOrder(symbol, position, openclose, price)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Sprintf("%+v\n", res)
}

func TestSendStopLoss(t *testing.T) {
	//	func (ty *BasicInfo) SendStopLoss(symbol, position, openclose, price string) error {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	position := "LONG"
	openclose := "CLOSE"
	price := "0.10700"
	res, err := basic.SendStopLoss(symbol, position, openclose, price)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Sprintf("%+v\n", res)
}

func TestSendStopMarket(t *testing.T) {
	//	func (ty *BasicInfo) SendStopMarket(symbol, positionside, openclose, price, amount string) error {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	positionside := "SHORT"
	openclose := "OPEN"
	price := "0.1070"
	amount := "1000"
	res, err := basic.SendStopMarket(symbol, positionside, openclose, price, amount)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Sprintf("%+v\n", res)
}

func TestCancelOrder(t *testing.T) {
	//	func (ty *BasicInfo) CancelOrder(symbol, ClientOrderID string) (*futures.CancelOrderResponse, error) {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	ClientOrderID := "123123123"
	res, err := basic.CancelOrder(symbol, ClientOrderID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Sprintf("%+v\n", res)
}

func TestSendRemoveOpenOrder(t *testing.T) {
	//	func (ty *BasicInfo) SendRemoveOpenOrder() []string {
	basic := GetBasicInit()
	res := basic.SendRemoveOpenOrder()
	fmt.Sprintf("%+v\n", res)
}

func TestSendClosePosition(t *testing.T) {
	//	func (ty *BasicInfo) SendClosePosition(position string) error {
	basic := GetBasicInit()
	position := "LONG"
	res, err := basic.SendClosePosition(position)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Sprintf("%+v\n", res)

}
