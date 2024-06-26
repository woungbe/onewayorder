package service

import (
	"fmt"
	"log"
	"onewayorder/errors"
	"onewayorder/util"
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

	// fmt.Sprintf("%+v", basic)		  // DOGE 일때
	fmt.Println(basic.GetNotionalPrice()) // 5
	fmt.Println(basic.GetMaxPrice())      // 30
	fmt.Println(basic.GetMinPrice())      // 0.00244
	fmt.Println(basic.GetTickSizePrice()) // 0.000010
	fmt.Println(basic.GetMaxQuantity())   // 5e+07
	fmt.Println(basic.GetMinQuantity())   // 1
	fmt.Println(basic.GetStepSize())      // 1
}

// 미체결 리스트 리턴
func TestGetOpenOrder(t *testing.T) {
	basic := GetBasicInit()
	res := basic.GetOpenOrder("")
	if len(res) == 0 {
		fmt.Println(" 없음 ")
		return
	}

	for _, v := range res {
		//fmt.Println(v)
		str := util.JsonData(v)
		fmt.Println(str)
	}

}

func TestGetPositionList(t *testing.T) {
	basic := GetBasicInit()
	res, err := basic.GetPositionList()
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range res {
		fmt.Printf("%+v\n", v)
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

// Entry Price
func TestGetEntryPrice(t *testing.T) {
	basic := GetBasicInit()
	res, err := basic.GetEntryPrice("LONG")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)

}

// 가격,USDT 로 => 코인 수량 가져오기
func TestGetCoinQtyForPrice(t *testing.T) {
	basic := GetBasicInit()
	res, err := basic.GetCoinQtyForPrice("0.117", "200")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)

	amount, err := basic.ReturnCoinForSize(res)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("coin amount : ", amount)
}

// 가격, 포지션, 익절퍼센트 => 가격 가져오기
func TestGetTPPrice(t *testing.T) {
	//	func (ty *BasicInfo) GetTPPrice(price, positionSide, takePersent string) (string, error) {
	basic := GetBasicInit()
	err := basic.SetExChangeInfo()
	if err != nil {
		fmt.Println(err)
	}

	price := "0.0117"
	positionSide := "LONG"
	takePersent := "20"
	res, err := basic.GetTPPrice(price, positionSide, takePersent)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("res : ", res)
}

func TestGetSLPrice(t *testing.T) {
	//	func (ty *BasicInfo) GetSLPrice(price, positionSide, takePersent string) (string, error) {
	basic := GetBasicInit()
	err := basic.SetExChangeInfo()
	if err != nil {
		fmt.Println(err)
	}

	price := "0.11924"
	positionSide := "LONG"
	takePersent := "25"
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

// 가격 잘라주기 !!
func TestReturnPriceForSize(t *testing.T) {
	//	func (ty *BasicInfo) ReturnPriceForSize(price string) (string, error) {
	basic := GetBasicInit()
	err := basic.SetExChangeInfo()
	if err != nil {
		fmt.Println(err)
	}

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

// 미체결 테스트
func TestSendOpenOrder(t *testing.T) {

	basic := GetBasicInit()
	position := "LONG"
	side := "OPEN"
	price := "0.1017"
	amount := "10000"
	res, err := basic.SendOpenOrder(position, side, price, amount)
	if err != nil {
		fmt.Println(err)
	}

	data := util.JsonData(res)
	fmt.Printf("%+v\n", data)

	errmsg := res.Errors[0]
	codemsg, err := util.ParseError(errmsg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("errmsg.msg : %v %s", codemsg.Code, codemsg.Message)

}

// 마켓주문 테스트
func TestSendMarketOrder(t *testing.T) {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	position := "LONG"
	openclose := "OPEN"
	price := "300"
	res, err := basic.SendMarketOrder(symbol, position, openclose, price)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
}

// 익절 테스트
func TestSendTakeProfit(t *testing.T) {
	//	func (ty *BasicInfo) SendTakeProfit(symbol, position, openclose, price string) (*futures.CreateOrderResponse, error) {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	position := "LONG"
	openclose := "CLOSE"
	price := "0.1270"
	res, err := basic.SendTakeProfit(symbol, position, openclose, price)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
}

// 손절 테스트
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

	fmt.Printf("%+v\n", res)
}

// stop Market 조건부 주문
func TestSendStopMarket(t *testing.T) {
	//	func (ty *BasicInfo) SendStopMarket(symbol, positionside, openclose, price, amount string) error {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	positionside := "SHORT"
	openclose := "OPEN"
	price := "0.12258"
	amount := "600"
	res, err := basic.SendStopMarket(symbol, positionside, openclose, price, amount)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
}

// 주문 취소하기 -
func TestCancelOrder(t *testing.T) {
	//	func (ty *BasicInfo) CancelOrder(symbol, ClientOrderID string) (*futures.CancelOrderResponse, error) {
	basic := GetBasicInit()
	symbol := "DOGEUSDT"
	ClientOrderID := "I1XqnYAUifMYT6S0KROc8W" // ClientOrderID
	res, err := basic.CancelOrder(symbol, ClientOrderID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
}

// 해당 Symbol 미체결 모두 제거
func TestSendRemoveOpenOrder(t *testing.T) {
	//	func (ty *BasicInfo) SendRemoveOpenOrder() []string {
	basic := GetBasicInit()
	res := basic.SendRemoveOpenOrder()
	fmt.Printf("%+v\n", res)
}

// 해당 Symbol 포지션 제거
func TestSendClosePosition(t *testing.T) {
	//	func (ty *BasicInfo) SendClosePosition(position string) error {
	basic := GetBasicInit()
	position := "LONG"
	// position := "SHORT"
	res, err := basic.SendClosePosition(position)
	if err != nil {
		fmt.Println(err)
		apierr, err := util.ParseError(err)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("apierr: %s , %s\n", apierr.Code, apierr.Message)
		}
	}

	str := util.JsonData(res)
	fmt.Println("str : ", str)
}
