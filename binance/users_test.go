package binance

import (
	"fmt"
	"log"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/spf13/viper"
)

func GetEnv() (string, string) {
	viper.SetConfigFile("../.env") // .env 파일 설정
	viper.AutomaticEnv()           // 환경 변수를 자동으로 읽도록 설정
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	accKey := viper.GetString("AccessKey")
	seckey := viper.GetString("SecretKey")
	return accKey, seckey
}

func GetUsrs() *BinanceUser {
	AccessKey, SecritKey := GetEnv()
	binance := new(BinanceUser)
	binance.Init(AccessKey, SecritKey)
	return binance
}

func TestGetStartUserStreamService(t *testing.T) {
	res, err := GetUsrs().GetStartUserStreamService()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestGetLeverageBracket(t *testing.T) {
	res, err := GetUsrs().GetLeverageBracket("ETHUSDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

// positionModel
func TestGetChangePositionModeService(t *testing.T) {
	err := GetUsrs().GetChangePositionModeService(true)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetBalanceService(t *testing.T) {
	res, err := GetUsrs().GetBalanceService()
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range res {
		// fmt.Println(v.AccountAlias, v.Asset, v.Balance, v.CrossWalletBalance, v.CrossUnPnl, v.AvailableBalance, v.MaxWithdrawAmount)
		if v.Asset == "USDT" {
			fmt.Println("v.AvailableBalance ", v.AvailableBalance)
		}
	}
}

func TestGetPositionRiskService(t *testing.T) {
	res, err := GetUsrs().GetPositionRiskService("ETHUSDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestGetListOpenOrdersService(t *testing.T) {
	res, err := GetUsrs().GetListOpenOrdersService("ETHUSDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCreateMuiOrder(t *testing.T) {
	var send []*futures.CreateOrderService
	var order OpenOrder
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeBuy                   // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeShort // PositionSideTypeLong PositionSideTypeShort
	order.Type = futures.OrderTypeLimit                // OrderTypeLimit OrderTypeMarket
	order.Quantity = "500"
	order.Price = "0.2100"
	order.TimeInForce = futures.TimeInForceTypeGTC

	createOrderService := CreateOrderLimitMarket(order)
	send = append(send, createOrderService)
	res, err := GetUsrs().CreateMuiOrder(send)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCreateMuiOrderMarket(t *testing.T) {
	var send []*futures.CreateOrderService
	var order OpenOrder
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeSell                  // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeShort // PositionSideTypeLong PositionSideTypeShort
	order.Type = futures.OrderTypeMarket               // OrderTypeLimit OrderTypeMarket
	order.Quantity = "500"
	// order.Price = "0.2177"
	// order.TimeInForce = futures.TimeInForceTypeGTC

	createOrderService := CreateOrderLimitMarket(order)
	send = append(send, createOrderService)
	res, err := GetUsrs().CreateMuiOrder(send)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCreateMuiOrderTakeProfit(t *testing.T) {
	// var send []*futures.CreateOrderService
	var order OrderType
	order.Symbol = "BIGTIMEUSDT"
	order.Side = futures.SideTypeBuy                    // SideTypeBuy SideTypeSell
	order.PositionSide = futures.PositionSideTypeShort  // PositionSideTypeLong PositionSideTypeShort
	order.OrderType = futures.OrderTypeTakeProfitMarket // OrderTypeLimit OrderTypeMarket OrderTypeTakeProfitMarket
	order.StopPrice = "0.210"
	order.ClosePosition = true
	order.WorkingType = futures.WorkingTypeContractPrice // WorkingTypeContractPrice
	// createOrderService := CreateOrderLimitMarket(order)
	// send = append(send, createOrderService)
	// res, err := GetUsrs().CreateMuiOrder(send)
	res, err := GetUsrs().CreateOrderService(order)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
