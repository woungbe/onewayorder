package binance

import (
	"fmt"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
)

func GetUsrs() *BinanceUser {
	binance := new(BinanceUser)
	binance.AccessKey = "S4LcIeMURpw9dUGBINKxJ2so7tEfF9seGWOZqf9f9FKU3ffDEohv5Xgx9dgL2Kzs"
	binance.SecritKey = "f08snk1B2BJfRL2gXPF30EenHGCD1KQWFeaCW5HUPC3a9eGR4SM8NRx5bOLQeP5P"
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
	fmt.Println(res)
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
