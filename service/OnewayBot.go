package service

import (
	"fmt"
	"math"
	"math/rand"
	"onewayorder/binance"
	"onewayorder/errors"
	"strings"
	"time"

	"github.com/woungbe/utils"
)

// 설정 정보 데이터
type BotConfigData struct {
	BotName        string `json:"BotName"`        // 봇이름
	TotalUSDT      string `json:"TotalUSDT"`      // 총금액
	Symbol         string `json:"Symbol"`         // 종목
	Leverage       string `json:"Leverage"`       // 레버리지
	Increment      string `json:"Increment"`      // 배율
	OrderCount     int    `json:"OrderCount"`     // 진입횟수
	LSSpacing      string `json:"LSSpacing"`      // 롱숏 간격
	TakeProfit     string `json:"TakeProfit"`     // 익절 간격 (롱익절, 숏익절 통일 )
	RepeatFlg      bool   `json:"RepeatFlg"`      // 반복실행 여부 : true: 반복실행, false: 안함
	LastOrderHedge bool   `json:"LastOrderHedge"` // 막회차 헷지 - 아직 작업 안함
	FirstOrder     string `json:"FirstOrder"`     // 첫번째 - LONG SHORT, RANDOM
}

// 봇 실행 정보
type ProcData struct {
}

type InterfaceBot interface {
	GetBotName() string
	SetConfigData(map[string]interface{}) (bool, error)
}

type OnewayBot struct {
	// AccessKey string
	// SecritKey string
	mBotName    string
	mConfigData BotConfigData // 설정 정보 데이터
	// mProcData   ProcData      // 봇 실행 정보

	mActFlg bool // 동작 여부
	mRound  int  // 현재 회차

	mbinanceAccount *BinanceAccount
	FirstPosside    string // LONG , SHORT

	NormalPrice string // 최소수량
	CoinDecimal string // 코인 수량
}

func (ty *OnewayBot) GetBotName() string {
	return ty.mBotName
}

// 초기화 필요한 선언은 여기서 하세요
func (ty *OnewayBot) Init(args *BinanceAccount) {
	ty.mActFlg = false
	ty.mRound = 1
	ty.mbinanceAccount = args
}

// 세팅 정보 저장하기
func (ty *OnewayBot) SetConfigData(args map[string]interface{}) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			v := fmt.Sprintf("%+v\n", args)
			errors.Error("Crit defer", "OnewayBot.SetConfigData ", err, v)
		}
	}()

	err := utils.MapToStruct("json", args, &ty.mConfigData)
	if err != nil {
		errors.Error("Crit Panic", "OnewayBot.SetConfigData ", err)
		return false, err
	}

	return true, nil
}

func (ty *OnewayBot) SetExChangeInfo() error {
	for _, v := range ty.mbinanceAccount.ExchangeInfo {
		if v.Symbol == ty.mConfigData.Symbol {
			minnotmal := v.MinNotionalFilter()
			lotsize := v.LotSizeFilter()
			ty.NormalPrice = minnotmal.Notional
			ty.CoinDecimal = lotsize.MinQuantity
		}
	}
	return nil
}

// 실행하기
func (ty *OnewayBot) Run() bool {
	ty.mActFlg = true
	ty.mRound = 1

	// 미체결 체크
	if chkOpenOrder := ty.CheckOpenOrder(); !chkOpenOrder {
		return chkOpenOrder
	}

	// 포지션 체크
	if chkPosition := ty.CheckPosition(); !chkPosition {
		return chkPosition
	}

	// 금액 확인 - 에이씨
	// minprice := ty.CheckMinPrice()
	// marginAmount, err := utils.Int64(ty.mConfigData.TotalUSDT)
	// if err != nil {
	// 	errors.Error("Crit ", "OnewayBot Run : ", marginAmount)
	// 	return false
	// }

	// 내 마진금액이 최소금액보다 낮으면 에러 ( 동일도 하지 말자 )
	// if marginAmount < minprice {
	// 	errors.Error("Crit ", "marginAmount , minprice : ", marginAmount, minprice)
	// 	return false
	// }

	// 주문하기
	err := ty.SendFirstOrder()
	if err != nil {
		errors.Error("Crit ", "SendFirstOrder : ", err)
		return false
	}

	return true
}

// 미체결이 있는지 체크
func (ty *OnewayBot) CheckOpenOrder() bool {
	// var b bool = false
	// 미체결 작업
	// res, err := ty.mbinanceAccount.mBinanceAPI.GetListOpenOrdersService(ty.mConfigData.Symbol)
	// 포지션 정보 갱신
	b := ty.mbinanceAccount.GetOpenOrderList()
	if b {
		errors.Error("have a openorder")
		return false
	}

	// 로드는 잘됐고
	miorderList := ty.mbinanceAccount.GetUserOpenOrders()
	for _, v := range miorderList {
		if v.Symbol == ty.mConfigData.Symbol {
			return true
		}
	}
	return b
}

// 포지션이 있는지 체크
func (ty *OnewayBot) CheckPosition() bool {
	// 포지션 요청
	b := ty.mbinanceAccount.GetPositionList()
	if b {
		errors.Error("CheckPosition :  GetPositionList")
		return false
	}
	// 포지션 데이터 가져와서 확인
	res := ty.mbinanceAccount.GetUserPositions()
	for _, v := range res {
		if v.Symbol == ty.mConfigData.Symbol {
			return true
		}
	}
	return false
}

// 체결이 되고 나서 평단가 구하는 구하기
func (ty *OnewayBot) GetEntryPrice(longshort string) string {
	b := ty.mbinanceAccount.GetPositionList()
	if b {
		msg := fmt.Sprintf("GetEntryPrice : GetPositionList")
		errors.Error(msg)
		return ""
	}

	res := ty.mbinanceAccount.GetUserPositions()
	for _, v := range res {
		if v.PositionSide == longshort && v.Symbol == ty.mConfigData.Symbol {
			return v.EntryPrice
		}
	}
	return ""
}

// 최소 금액 계산
func (ty *OnewayBot) CheckMinPrice() int64 {
	var send int64
	fmt.Println("최소금액 체크했습니다. ")

	// 1000 을 5회차로 들어간도
	// ty.mConfigData.TotalUSDT
	// 총금액으로 가능한가 ?  를 따져보는거지.

	return send
}

// 주문하기
func (ty *OnewayBot) SendFirstOrder() error {

	ty.mActFlg = true
	ty.mRound = 1

	// ticker 로 현재가를 가져온다.

	symbol := ty.mConfigData.Symbol
	positionside := ty.GetPositionSide(true)
	side := ty.GetSide(positionside, "OPEN")
	currentprice := ty.GetCurrentPrice()
	amount := ty.GetAmount(currentprice)

	// 시장가 주문
	err := ty.sendMarketOrder(symbol, positionside, side, amount)
	if err != nil {
		errors.Error("Crit ", "OnewayBot sendMarketOrder ", err)
		return err
	}

	time.Sleep(time.Second * 2)

	// 포지션을 가져옵니다. - 현재 포지션을 체크합니다.
	entryPrice := ty.GetEntryPrice(positionside)

	// 익절 - symbol string, position string, openclose string, price string
	err = ty.sendTakeProfit(symbol, positionside, "CLOSE", ty.GetTakePrice(entryPrice, positionside))
	if err != nil {
		errors.Error("Crit ", "OnewayBot sendTakeProfit ", err)
		return err
	}

	// 손절
	err = ty.sendStopLoss(symbol, positionside, "CLOSE", ty.GetStopLoss(entryPrice, positionside))
	if err != nil {
		errors.Error("Crit ", "OnewayBot sendStopLoss ", err)
		return err
	}

	// 스탑 마켓 - symbol string, positionside string, openclose string, price string, amount string
	ty.mRound++
	err = ty.sendStopMarket(symbol, ty.GetPositionSide(false), "OPEN", ty.GetOtherPrice(entryPrice, positionside), ty.GetAmount(entryPrice))
	if err != nil {
		errors.Error("Crit ", "OnewayBot sendStopMarket ", err)
		return err
	}
	return nil
}

// 두번째가 체결된 이후 조건 주문 - 체결된 포지션 side
func (ty *OnewayBot) SendSecond(position string) error {
	// fmt.Println("두번째 체결 이후 작업 ")

	// 라운드
	// 금액설정 => 수량 변환
	// 넣을 가격
	ty.mRound++
	symbol := ty.mConfigData.Symbol

	// 체결된 포지션의 반대를 넣어야되는 상황일세 ?
	// positionSide := ty.GetPositionSide(true)
	var positionSide string
	if position == "LONG" {
		positionSide = "SHORT"
	} else {
		positionSide = "LONG"
	}

	stopPrice := ty.GetPrice(positionSide)
	amount := ty.GetAmount(stopPrice)

	// 스탑마켓 넣어야지
	err := ty.sendStopMarket(symbol, positionSide, "OPEN", stopPrice, amount)
	if err != nil {
		errors.Error("Crit ", "OnewayBot sendStopMarket ", err)
		return err
	}

	/*
		익절손절을 재설정 하는 방법도 있고
		취소하고 다시 등록하는 방법도 있지만.
		그냥 유지하는 방법도 있음. - 여러개 하지뭐..

		상황에 따라서 각각 다를 수 있음.
	*/

	// other side 의 익절 손절을 걸자
	otherSide := position
	entryPrice := ty.GetPrice(otherSide)

	// 익절 - symbol string, position string, openclose string, price string
	err = ty.sendTakeProfit(symbol, otherSide, "CLOSE", ty.GetTakePrice(entryPrice, otherSide))
	if err != nil {
		errors.Error("Crit ", "OnewayBot sendTakeProfit ", err)
		return err
	}

	// 손절
	err = ty.sendStopLoss(symbol, otherSide, "CLOSE", ty.GetStopLoss(entryPrice, otherSide))
	if err != nil {
		errors.Error("Crit ", "OnewayBot sendStopLoss ", err)
		return err
	}

	return nil
}

// 마켓주문하기
func (ty *OnewayBot) sendMarketOrder(symbol, position, openclose, price string) error {
	res, err := ty.mbinanceAccount.mBinanceAPI.SendOrderMarket(symbol, position, openclose, price)
	if err != nil {
		errors.Error("OnewayBot : SendOrderMarket", err)
		return err
	}
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendTakeProfit ", log)
	return nil
}

// TP 주문
func (ty *OnewayBot) sendTakeProfit(symbol, position, openclose, price string) error {
	errors.Log("sendTakeProfit : ", symbol, position, openclose, price)
	res, err := ty.mbinanceAccount.mBinanceAPI.SendOrderTakeProfit(symbol, position, openclose, price)
	if err != nil {
		errors.Error("OnewayBot : SendOrderTakeProfit", err)
		return err
	}
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendTakeProfit ", log)
	return nil
}

// SL 주문
func (ty *OnewayBot) sendStopLoss(symbol, position, openclose, price string) error {
	res, err := ty.mbinanceAccount.mBinanceAPI.SendOrderStopLoss(symbol, position, openclose, price)
	if err != nil {
		errors.Error("OnewayBot : sendStopLoss", err)
		return err
	}
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendStopMarket ", log)
	return nil
}

// stop market 주문 (반대쪽 매매 )
func (ty *OnewayBot) sendStopMarket(symbol, positionside, openclose, price, amount string) error {
	// 이렇게 돌아가는 이유는 비상시에 여기다가 로그를 남기거나 에러를 남길 수 있는 장치를 넣을 수 있음.
	res, err := ty.mbinanceAccount.mBinanceAPI.SendOrderStopMarket(symbol, positionside, openclose, price, amount)
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendStopMarket ", log)
	return err
}

// 모든 주문 정리하기
func (ty *OnewayBot) removeAllOpenOrder() error {
	// 모든 openorder 정리하기
	// ty.mbinanceAccount.mBinanceAPI.
	// GetUserOpenOrders()
	res := ty.mbinanceAccount.GetUserOpenOrders()
	var orderids []int64
	for _, v := range res {
		orderids = append(orderids, v.OrderID)
	}

	ty.mbinanceAccount.mBinanceAPI.SendRemoveOpenOrder(ty.mConfigData.Symbol, orderids)
	return nil
}

// side 가져오기
func (ty *OnewayBot) GetSide(posside, openClose string) string {

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

// 현재가 가져오기
func (ty *OnewayBot) GetCurrentPrice() string {

	symbol := ty.mConfigData.Symbol
	res, err := binance.GetTicker(symbol)
	if err != nil {
		errors.Error("Error - OnewayBot GetCurrentPrice ", err)
		return ""
	}

	for _, v := range res {
		if v.Symbol == symbol {
			return v.LastPrice
		}
	}

	return ""
}

// positionSide 변경 하기
func (ty *OnewayBot) GetPositionSide(b bool) string {

	// 첫번째와 동일하면 true
	// 첫번째와 반대면 flase

	// 등록되어 있지 않고,
	if ty.FirstPosside == "" {
		// 랜덤으로 선택했다면
		if ty.mConfigData.FirstOrder == "RANDOM" {
			// 첫번째 주문은 랜덤으로 들어가세요.
			ty.FirstPosside = randomLongShort()
			return ty.FirstPosside

		} else {
			if ty.mConfigData.FirstOrder == "" {
				errors.Error("FirstOrder 설정이 되어 있지 않습니다. ")
				return ""
			}

			ty.FirstPosside = ty.mConfigData.FirstOrder
			return ty.FirstPosside
		}
	}

	side := strings.ToUpper(ty.FirstPosside)

	// 첫번째와 동일하다면
	if b {
		return side
	} else {
		if side == "LONG" {
			return "SHORT"
		}

		if side == "SHORT" {
			return "LONG"
		}
	}

	return ""
}

// 반대 price 계산기
func (ty *OnewayBot) GetOtherPrice(mEntryPrice, position string) string {
	// 현물 기준으로
	lss, err := utils.Float64(ty.mConfigData.LSSpacing)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	entryPrice, err := utils.Float64(mEntryPrice)
	if err != nil {
		return ""
	}

	var otherPrice float64
	if position == "LONG" {
		otherPrice = entryPrice * ((100 - lss) / 100)
	} else {
		otherPrice = entryPrice * ((100 + lss) / 100)
	}

	// 현물
	return utils.String(otherPrice)
}

// 구매가격 기준을 - 코인 수량 가져오기 ,  coin 소수점
func (ty *OnewayBot) GetAmount(mPrice string) string {

	price, _ := utils.Float64(mPrice)
	multiplier := ty.mConfigData.Leverage
	USDTArray := FirstAmount(ty.mConfigData.TotalUSDT, ty.mRound, multiplier)
	USDT := USDTArray[ty.mRound+1]
	decimal := ty.CoinDecimal
	if USDT == 0 {
		return ""
	}

	leverage := ty.mConfigData.Leverage
	// 총 증거금, 횟수, 배율 => 총 증거금 리스트

	// 현재가, 구매증거금, 레버리지, 코인 소수점
	amount := Amount(price, USDT, leverage, decimal)
	send := utils.String(amount)
	return send
}

func FirstAmount(totalMargin string, rounds int, multip string) []float64 {

	totalInvestment, _ := utils.Float64(totalMargin)
	multiplier, _ := utils.Float64(multip)

	// 첫 회차의 투자 금액 계산
	firstInvestment := totalInvestment / ((1 - math.Pow(multiplier, float64(rounds))) / (1 - multiplier))

	// 각 회차의 투자 금액 계산
	investments := make([]float64, rounds)
	for i := 0; i < rounds; i++ {
		investments[i] = firstInvestment * math.Pow(multiplier, float64(i))
	}
	return investments
}

// 증거금으로 수량 가져오기 // 수량  = 금액 / 현재 가격  3830, 6.5, 0.001,
func Amount(price float64, usdt float64, leverage string, decimal string) string {
	// 수량 = 증거금 * 레버리지 / 가격   -- 소수점
	tprice, _ := utils.Float64(price)
	tusdt, _ := utils.Float64(usdt)
	tleverage, _ := utils.Float64(leverage)
	tdecimal, _ := utils.Float64(decimal)
	qty := tusdt * tleverage / tprice

	scaledValue := qty / tdecimal
	truncatedScaledValue := math.Trunc(scaledValue)
	finalValue := truncatedScaledValue * tdecimal
	send := utils.String(finalValue)
	return send
}

// 익절값을 가져오는 func
func (ty *OnewayBot) GetTakePrice(mCurrentPrice string, longShort string) string {
	// 익절인데. 두개를 한다.
	// LONG 익절,
	// SHORT 익절,

	// 현재가 대비 - ticker
	// 마켓주문을 하고 그거대비로 하자 .
	var tmp float64
	currentPrice, err1 := utils.Float64(mCurrentPrice)
	if err1 != nil {
		errors.Error("GetTakePrice : ", err1)
		return ""
	}

	takeProfit, err2 := utils.Float64(ty.mConfigData.TakeProfit)
	if err2 != nil {
		errors.Error("GetTakePrice : ", err2)
		return ""
	}

	leverage, err3 := utils.Float64(ty.mConfigData.Leverage) //
	if err3 != nil {
		errors.Error("GetTakePrice : ", err3)
		return ""
	}

	if longShort == "LONG" {
		tmp = currentPrice * (1 + (takeProfit/100)*leverage)
	} else if longShort == "SHORT" {
		tmp = currentPrice * (1 - (takeProfit/100)*leverage)
	}

	send := utils.String(tmp)
	return send

}

// 로스컷 정리
func (ty *OnewayBot) GetStopLoss(mCurrentPrice string, longShort string) string {
	// 익절인데. 두개를 한다.
	// LONG 익절,
	// SHORT 익절,

	// 현재가 대비 - ticker
	// 마켓주문을 하고 그거대비로 하자 .
	var tmp float64

	currentPrice, err := utils.Float64(mCurrentPrice)
	if err != nil {
		errors.Error("GetStopLoss : ", err)
		return ""
	}

	tmpLoss, err2 := utils.Float64(ty.mConfigData.TakeProfit)
	if err2 != nil {
		errors.Error("GetTakePrice : ", err2)
		return ""
	}

	// 간격
	tmpPeroid, err3 := utils.Float64(ty.mConfigData.LSSpacing)
	if err3 != nil {
		errors.Error("LSSpacing : ", err3)
		return ""
	}

	// 로스컷 - takeProfit + 간견
	lossProfit := tmpLoss + tmpPeroid
	leverage, err4 := utils.Float64(ty.mConfigData.Leverage) //
	if err4 != nil {
		errors.Error("GetTakePrice : ", err4)
		return ""
	}

	if longShort == "LONG" {
		tmp = currentPrice * (1 - (lossProfit/100)*leverage)
	} else if longShort == "SHORT" {
		tmp = currentPrice * (1 + (lossProfit/100)*leverage)
	}

	send := utils.String(tmp)
	return send

}

// 주문 가격 정리 - 해당 GetPrice
func (ty *OnewayBot) GetPrice(side string) string {
	// 포지션이 있다면 그걸로 한다.
	// 1은 현재를 가져오기 때문에 괜찮음
	if ty.mRound == 1 {
		// 시장가를 가져옵니다.
		return ty.GetCurrentPrice()
	}

	res := ty.mbinanceAccount.GetUserPositions()
	key := ""

	if side == "LONG" {
		key = ty.mConfigData.Symbol + "_LONG"
	} else if side == "SHORT" {
		key = ty.mConfigData.Symbol + "_SHORT"
	}

	position := res[key]
	return position.EntryPrice

}

// 롱숏 랜덤
func randomLongShort() string {
	options := []string{"LONG", "SHORT"}
	rand.Seed(time.Now().UnixNano())
	return options[rand.Intn(len(options))]
}
