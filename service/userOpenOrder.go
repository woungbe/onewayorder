package service

import "github.com/adshao/go-binance/v2/futures"

type UserOpenOrder struct {
	Symbol        string
	OrderID       int64
	ClientOrderID string //주문번호(될수있음 이걸로 주문취소나 찾기를 하자)
	Type          string //주문 타입 (STOP_MARKET ,TAKE_PROFIT_MARKET 이두가지만 사용 )
	Side          string //주문 사이드 (BUY,SELL)
	PositionSide  string //포지션 사이드 (LONG,SHROT)
	Amount        string //주문 수량
	Price         string //주문가
	StopPrice     string //
}

// 주문 기입
func (ty *UserOpenOrder) SetOpenOrder(args futures.Order) *UserOpenOrder {
	ty.Symbol = args.Symbol
	ty.OrderID = args.OrderID
	ty.ClientOrderID = args.ClientOrderID
	ty.Type = string(args.Type)
	ty.Side = string(args.Side)
	ty.PositionSide = string(args.PositionSide)
	ty.Amount = args.OrigQuantity
	ty.Price = args.Price
	ty.StopPrice = args.StopPrice
	return ty
}
