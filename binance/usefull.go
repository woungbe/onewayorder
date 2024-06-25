package binance

import (
	"fmt"
	"strings"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/woungbe/utils"
)

// 미체결 주문 하기
func (ty *BinanceUser) SendOpenOrder(symbol, position, openclose, price, amount string) (*futures.CreateBatchOrdersResponse, error) {
	var send []*futures.CreateOrderService
	var order OpenOrder
	order.Symbol = symbol
	order.PositionSide = PositionSide(position)         // PositionSideTypeLong PositionSideTypeShort
	order.Side = SideType(GetSide(position, openclose)) // SideTypeBuy SideTypeSell
	order.Type = futures.OrderTypeLimit                 // OrderTypeLimit OrderTypeMarket
	order.Price = price
	order.Quantity = amount
	order.TimeInForce = futures.TimeInForceTypeGTC
	createOrderService := ty.CreateOrderLimitMarket(order)
	send = append(send, createOrderService)
	res, err := ty.CreateMuiOrder(send)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	if err != nil {
		return nil, err
	}
	return res, err
}

// 마켓 주문
func (ty *BinanceUser) SendOrderMarket(symbol, position, openclose, amount string) (*futures.CreateOrderResponse, error) {
	var order OrderType
	order.Symbol = symbol
	order.PositionSide = PositionSide(position)         // PositionSideTypeLong PositionSideTypeShort
	order.Side = SideType(GetSide(position, openclose)) // SideTypeBuy SideTypeSell
	order.OrderType = futures.OrderTypeMarket           // OrderTypeLimit OrderTypeMarket
	order.Quantity = amount
	order.TimeInForce = futures.TimeInForceTypeGTC
	res, err := ty.CreateOrderService(order)
	if err != nil {
		return nil, err
	}
	return res, err
}

// take profit - 익절 주문
func (ty *BinanceUser) SendOrderTakeProfit(symbol, position, openclose, price string) (*futures.CreateOrderResponse, error) {
	var order OrderType
	order.Symbol = symbol
	order.PositionSide = PositionSide(position)         // PositionSideTypeLong PositionSideTypeShort
	order.Side = SideType(GetSide(position, openclose)) // SideTypeBuy SideTypeSell
	order.OrderType = futures.OrderTypeTakeProfitMarket // OrderTypeLimit OrderTypeMarket OrderTypeTakeProfitMarket
	order.StopPrice = price
	order.ClosePosition = true
	order.WorkingType = futures.WorkingTypeContractPrice // WorkingTypeContractPrice
	res, err := ty.CreateOrderService(order)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// StopLoss - 모든 수량을 팔아서 포지션을 정리함. - 손절 주문
func (ty *BinanceUser) SendOrderStopLoss(symbol, position, openclose, price string) (*futures.CreateOrderResponse, error) {
	/*
		symbol: BIGTIMEUSDT
		side: SELL
		positionSide: LONG
		type: STOP_MARKET
		stopPrice: 0.2051
		closePosition: true
	*/
	var order OrderType
	order.Symbol = symbol
	order.PositionSide = PositionSide(position)         // PositionSideTypeLong PositionSideTypeShort
	order.Side = SideType(GetSide(position, openclose)) // SideTypeBuy SideTypeSell
	order.OrderType = futures.OrderTypeStopMarket       // OrderTypeLimit OrderTypeMarket OrderTypeTakeProfitMarket
	order.StopPrice = price
	order.ClosePosition = true
	order.WorkingType = futures.WorkingTypeContractPrice // WorkingTypeContractPrice
	res, err := ty.CreateOrderService(order)
	return res, err

}

// stopmarket 은 포지션을 종료하지 않음 - 수량대로 더 사거나, 더 팔 수 있음
// BITTIMEUSDT, SHORT, OPEN, "0.210", "500"
func (ty *BinanceUser) SendOrderStopMarket(symbol, position, openclose, price, amount string) (*futures.CreateOrderResponse, error) {
	var order OrderType
	order.Symbol = symbol
	order.PositionSide = PositionSide(position)         // PositionSideTypeLong PositionSideTypeShort
	order.Side = SideType(GetSide(position, openclose)) // SideTypeBuy SideTypeSell
	order.OrderType = futures.OrderTypeStopMarket       // OrderTypeLimit OrderTypeMarket OrderTypeTakeProfitMarket
	order.StopPrice = "0.2150"
	order.Quantity = "500"
	order.ClosePosition = false
	order.WorkingType = futures.WorkingTypeContractPrice // WorkingTypeContractPrice
	res, err := ty.CreateOrderService(order)
	return res, err
}

// 트레일링 스탑
func (ty *BinanceUser) SendOrderTrailingStop(symbol, position, openclose, price, amount string) (*futures.CreateOrderResponse, error) {
	/*
		symbol: BIGTIMEUSDT
		side: SELL
		positionSide: LONG
		type: TRAILING_STOP_MARKET
		quantity: 500
		callbackRate: 1
		workingType: CONTRACT_PRICE
		activationPrice: 0.2104
	*/
	var order OrderType
	order.Symbol = symbol
	order.PositionSide = PositionSide(position)         // PositionSideTypeLong PositionSideTypeShort
	order.Side = SideType(GetSide(position, openclose)) // SideTypeBuy SideTypeSell
	order.OrderType = futures.OrderTypeStopMarket       // OrderTypeLimit OrderTypeMarket OrderTypeTakeProfitMarket
	order.ActivationPrice = price
	order.Quantity = amount
	order.WorkingType = futures.WorkingTypeContractPrice // WorkingTypeContractPrice
	res, err := ty.CreateOrderService(order)
	return res, err
}

// 해당 심볼의 미체결을 모두 취소하기
func (ty *BinanceUser) SendRemoveOpenOrderForSymbol(symbol string) []string {
	// ([]*futures.Order, error)
	res, err := ty.GetListOpenOrdersService(symbol)
	if err != nil {
		fmt.Println(err)
	}

	var orderID []int64
	for _, v := range res {
		val := utils.String(v.ClientOrderID)
		oriClientID, err := utils.Int64(val)
		if err != nil {
			fmt.Println("utils.Int64 : ", err)
		}
		orderID = append(orderID, oriClientID)
	}

	msg := ty.SendRemoveOpenOrder(symbol, orderID)
	if len(msg) != 0 {
		return msg
	}
	var send []string
	return send
}

// 해당 미체결 모두 정리
func (ty *BinanceUser) SendRemoveOpenOrder(symbol string, orderID []int64) []string {
	var send []string
	for _, v := range orderID {
		val := utils.String(v)
		_, err := ty.CancelOrder(symbol, val)
		if err != nil {
			send = append(send, fmt.Sprintf("%s %s", val, err))
		}
	}
	return send
}

func SideType(sidetype string) futures.SideType {
	// Side = futures.SideTypeBuy                  // SideTypeBuy SideTypeSell
	if sidetype == "BUY" {
		return futures.SideTypeBuy
	}

	if sidetype == "BUY" {
		return futures.SideTypeSell
	}
	return ""
}

func PositionSide(positionside string) futures.PositionSideType {
	// PositionSide = futures.futures.PositionSideTypeLong // PositionSideTypeLong PositionSideTypeShort
	if positionside == "LONG" {
		return futures.PositionSideTypeLong
	}

	if positionside == "SHORT" {
		return futures.PositionSideTypeShort
	}

	return ""
}

// side 가져오기
func GetSide(posside, openClose string) string {

	tmpposside := strings.ToUpper(posside)
	tmpopenClose := strings.ToUpper(openClose)

	if tmpposside == "LONG" && tmpopenClose == "OPEN" {
		return "BUY"
	}
	if tmpposside == "LONG" && tmpopenClose == "CLOSE" {
		return "SELL"
	}
	if tmpposside == "LONG" && tmpopenClose == "OPEN" {
		return "SELL"
	}
	if tmpposside == "LONG" && tmpopenClose == "CLOSE" {
		return "BUY"
	}

	return ""

}
