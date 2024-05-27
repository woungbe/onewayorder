package service

import (
	"onewayorder/errors"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/woungbe/utils"
)

type UserPositionInfo struct {
	Symbol           string
	PositionSide     string //포지션 ( SHORT / LONG)
	EntryPrice       string //진입가 (AVG)
	MarginType       bool   //마진 타입(격리인경우 true)
	IsolatedMargin   string //격리마진 금
	PositionAmt      string //포지션 수량
	Leverage         int64  //레버리지
	LiquidationPrice string //청산가
	OrderTime        int64  //주문시간
	UnRealizedProfit string //수익가
}

// 주문 기입
func (ty *UserPositionInfo) SetPosition(args *futures.PositionRisk) (string, *UserPositionInfo) {
	key := ty.setPositionKey(args.Symbol, args.PositionSide)
	ty.Symbol = args.Symbol
	ty.PositionSide = args.PositionSide
	ty.EntryPrice = args.EntryPrice
	ty.MarginType = false
	if args.MarginType == "isolated" {
		ty.MarginType = true
	}
	ty.IsolatedMargin = args.IsolatedWallet
	ty.PositionAmt = args.PositionAmt

	ty.Leverage = ty.getLeverage(args.Leverage)
	ty.UnRealizedProfit = args.UnRealizedProfit
	ty.LiquidationPrice = args.LiquidationPrice
	return key, ty
}

// key 만들기
func (ty *UserPositionInfo) setPositionKey(symbol, posside string) string {
	if symbol == "" || posside == "" {
		return ""
	}

	send := symbol + "_" + posside
	return send
}

func (ty *UserPositionInfo) getLeverage(Leverage string) int64 {
	leverage, err := utils.Int64(Leverage)
	if err != nil {
		errors.Error("Crit Panic", "UserPositionInfo.getLeverage  ", err)
	}
	return leverage
}
