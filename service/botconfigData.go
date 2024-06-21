package service

import (
	"fmt"
	"onewayorder/binance"
	"onewayorder/errors"
)

// 기본정보 가져오기
type BasicInfo struct {
	mbinanceAccount *BinanceAccount // 바이낸스 기본정보
	symbol          string          // 심볼
	leverage        string          // 레버리지
}

// 초기화 설정
func (ty *BasicInfo) Init(m *BinanceAccount, symbol string, Leverage string) {
	ty.mbinanceAccount = m
	ty.symbol = symbol
	ty.leverage = Leverage
}

// 최소단위, 코인 소수점 가져오기
func (ty *BasicInfo) SetExChangeInfo() (string, string, error) {
	for _, v := range ty.mbinanceAccount.ExchangeInfo {
		if v.Symbol == ty.symbol {
			minnotmal := v.MinNotionalFilter()
			lotsize := v.LotSizeFilter()
			return minnotmal.Notional, lotsize.MinQuantity, nil
		}
	}
	return "", "", fmt.Errorf("Error:BasicInfo.SetExchangeInfo() msg:not found symbol ")
}

// 미체결 리스트 리턴 - 걸려있는 것들만.
func (ty *BasicInfo) GetOpenOrder() map[int64]UserOpenOrder {
	// 요청하고,
	b := ty.mbinanceAccount.GetOpenOrderList()
	if b {
		errors.Error("have a openorder")
		return nil
	}

	// 로드하고,
	send := make(map[int64]UserOpenOrder)
	miorderList := ty.mbinanceAccount.GetUserOpenOrders()
	for k, v := range miorderList {
		if v.Symbol == ty.symbol {
			send[k] = v
		}
	}
	return send
}

// 포지션 가져오기
func (ty *BasicInfo) GetPositionList() (map[string]UserPositionInfo, error) {
	// 포지션 요청
	b := ty.mbinanceAccount.GetPositionList()
	if b {
		emsg := "CheckPosition :  GetPositionList"
		errors.Error(emsg)
		return nil, fmt.Errorf(emsg)
	}
	// 포지션 데이터 가져와서 확인
	send := make(map[string]UserPositionInfo)
	res := ty.mbinanceAccount.GetUserPositions()
	for k, v := range res {
		if v.Symbol == ty.mConfigData.Symbol {
			send[k] = v
		}
	}
	return send, nil
}

//////// 데이터 갱신, 조합  ////////

// 현재가 가지오기
func (ty *BasicInfo) GetCurrentPrice() (string, error) {
	symbol := ty.symbol
	res, err := binance.GetTicker(symbol)
	if err != nil {
		errors.Error("Error - OnewayBot GetCurrentPrice ", err)
		return "", err
	}

	for _, v := range res {
		if v.Symbol == symbol {
			return v.LastPrice, nil
		}
	}

	return "", fmt.Errorf("Error BasicInfo.GetCurrentPrice msg=not found Symbol")
}

// 체결이후 평단가 가져오기
func (ty *BasicInfo) GetEntryPrice(longshort string) (string, error) {
	b := ty.mbinanceAccount.GetPositionList()
	if b {
		msg := fmt.Sprintf("GetEntryPrice : GetPositionList")
		errors.Error(msg)
		return "", fmt.Errorf("Error BasicInfo.GetEntryPrice ")
	}

	res := ty.mbinanceAccount.GetUserPositions()
	for _, v := range res {
		if v.PositionSide == longshort && v.Symbol == ty.symbol {
			return v.EntryPrice, nil
		}
	}
	return "", fmt.Errorf("Error BasicInfo.GetEntryPrice msg=not found Symbol")
}

// 최소 금액 계산하기
func (ty *OnewayBot) CheckMinPrice() int64 {
	var send int64
	fmt.Println("최소금액 체크했습니다. ")

	// 1000 을 5회차로 들어간도
	// ty.mConfigData.TotalUSDT
	// 총금액으로 가능한가 ?  를 따져보는거지.

	return send
}

// 구매가격,증거금 기준으로 코인 수량 가져오기

// 가격으로 익절가격 가져오기

// 가격으로 손절가격 가져오기

//////// 이벤트 영역 //////////

// 주문하기

// 마켓주문

// TL 주문

// SL 주문

// Stop Market 주문

// 모든 주문 정리
