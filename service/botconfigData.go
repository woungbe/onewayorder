package service

import (
	"fmt"
	"onewayorder/binance"
	"onewayorder/errors"
	util "onewayorder/util"
	"os"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/woungbe/utils"
)

/*
	주문을 하거단 계산을 할때 디폴트로 필요한 데이터 정보가 존재함.
	봇을 실행하면 잘 변경되지 않는 데이터
    : 심볼, 레버리지
	봇을 실행하면 잘 변경되는 데이터
	: side, positionSide, price, take profit, lose profit

	레버리지 버스켓도 여기서 관리해도 되겠네
	- 최대, 최소 관리
	- price ticker
	- USDT -> qty 로 변경하는 작업
	-

	// Init() 초기화 설정
	// 최소단위, 코인 소수점 가져오기
	// 미체결 가져오기
	// 포지션 가져오기
	// 현재가 가져오기
	// 평단가 가져오기 - (포지션 가져오는거에 중복)

	// 구매가격, 증거금 => 코인 수량 가져오기
	// 가격,Persent => 익절가격 계산하기
	// 가격,Persent => 손절가격 계산하기

	////////// 이벤트 영역 ///////////////

	// 주문하기
	// 마켓주문
	// TL주문
	// SL주문
	// Stop Market 주문
	// 미체결주문 취소
	// 미체결 모두 정리
	// 포지션 청산
*/

// 기본정보 가져오기
type BasicInfo struct {
	mbacc     *BinanceAccount // 바이낸스 기본정보
	symbol    string          // 심볼
	leverage  string          // 레버리지
	fleverage float64         // 레버리지 float64

	fNotionalPrice float64 // 최소 거래금액
	fMaxPrice      float64 // 최대 가격
	fMinPrice      float64 // 최소 가격
	fTickSizePrice float64 // 가격 tickSize
	fMaxQuantity   float64 // 최대 코인 수량
	fMinQuantity   float64 // 최소 코인 수량
	fStepSize      float64 // 코인 tickSize

}

// 초기화 설정
func (ty *BasicInfo) Init(m *BinanceAccount, symbol string, Leverage string) {
	ty.mbacc = m
	ty.symbol = symbol
	ty.leverage = Leverage

	fleverage, err := utils.Float64(Leverage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	ty.fleverage = fleverage
}

// 심볼 변경
func (ty *BasicInfo) SetSymbol(symbol string) {
	ty.symbol = symbol
}

// 레버리지 변경
func (ty *BasicInfo) SetLeverage(Leverage string) {
	ty.leverage = Leverage
	fleverage, err := utils.Float64(Leverage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	ty.fleverage = fleverage
}

// 최소단위, 코인 소수점 가져오기
func (ty *BasicInfo) SetExChangeInfo() error {

	if len(ty.mbacc.ExchangeInfo) == 0 {
		// 없으면 가져오세요 ~
		ty.mbacc.GetExchangeInfo()
	}

	for _, v := range ty.mbacc.ExchangeInfo {
		if v.Symbol == ty.symbol {
			// ty.NotionalPrice = v.MinNotionalFilter().Notional // 최소 거래금액
			NotionalPrice, err := utils.Float64(v.MinNotionalFilter().Notional) // 최소 거래금액
			if err != nil {
				return err
			}

			FMaxPrice, err := utils.Float64(v.PriceFilter().MaxPrice)
			if err != nil {
				return err
			}
			FMinPrice, err := utils.Float64(v.PriceFilter().MinPrice)
			if err != nil {
				return err
			}
			FTickSizePrice, err := utils.Float64(v.PriceFilter().TickSize)
			if err != nil {
				return err
			}
			FMaxQuantity, err := utils.Float64(v.LotSizeFilter().MaxQuantity)
			if err != nil {
				return err
			}
			FMinQuantity, err := utils.Float64(v.LotSizeFilter().MinQuantity)
			if err != nil {
				return err
			}
			FStepSize, err := utils.Float64(v.LotSizeFilter().StepSize)
			if err != nil {
				return err
			}

			ty.fNotionalPrice = NotionalPrice  // 최소 거래금액
			ty.fMaxPrice = FMaxPrice           // 최대 가격
			ty.fMinPrice = FMinPrice           // 최소 가격
			ty.fTickSizePrice = FTickSizePrice // 가격 tickSize
			ty.fMaxQuantity = FMaxQuantity     // 최대 코인 수량
			ty.fMinQuantity = FMinQuantity     // 최소 코인 수량
			ty.fStepSize = FStepSize           // 코인 tickSize

			return nil
		}
	}
	return fmt.Errorf("Error:BasicInfo.SetExchangeInfo() msg:not found symbol ")
}

// 미체결 리스트 리턴 - 걸려있는 것들만.
func (ty *BasicInfo) GetOpenOrder() map[int64]UserOpenOrder {
	// 요청하고,
	b := ty.mbacc.GetOpenOrderList()
	if b {
		errors.Error("have a openorder")
		return nil
	}

	// 로드하고,
	send := make(map[int64]UserOpenOrder)
	miorderList := ty.mbacc.GetUserOpenOrders()
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
	b := ty.mbacc.GetPositionList()
	if b {
		emsg := "CheckPosition :  GetPositionList"
		errors.Error(emsg)
		return nil, fmt.Errorf(emsg)
	}
	// 포지션 데이터 가져와서 확인
	send := make(map[string]UserPositionInfo)
	res := ty.mbacc.GetUserPositions()
	for k, v := range res {
		if v.Symbol == ty.symbol {
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
	b := ty.mbacc.GetPositionList()
	if b {
		msg := fmt.Sprintf("GetEntryPrice : GetPositionList")
		errors.Error(msg)
		return "", fmt.Errorf("Error BasicInfo.GetEntryPrice ")
	}

	res := ty.mbacc.GetUserPositions()
	for _, v := range res {
		if v.PositionSide == longshort && v.Symbol == ty.symbol {
			return v.EntryPrice, nil
		}
	}
	return "", fmt.Errorf("Error BasicInfo.GetEntryPrice msg=not found Symbol")
}

// 구매가격,증거금 기준으로 코인 수량 가져오기
func (ty *BasicInfo) GetCoinQtyForPrice(price, uset string) (string, error) {
	// 코인 수량을 구하는건데.
	// 레버리지는 알아서 포함하기

	send := ""
	// 구매가격
	fPrice, err := utils.Float64(price)
	if err != nil {
		return send, err
	}
	// 금액으로
	fUsdt, err := utils.Float64(uset)
	if err != nil {
		return send, err
	}

	// 증거금*레버리지 = 가격 * 수량
	amount := fUsdt * ty.fleverage / fPrice
	send = utils.String(amount)
	return send, nil
}

// 가격으로 익절가격 가져오기 - 20%, 30%, ...
func (ty *BasicInfo) GetTPPrice(price, positionSide, takePersent string) (string, error) {
	// 가격, 익절가격
	// 20% , 30% 라고 했을때...
	fpri, err := utils.Float64(price)
	if err != nil {
		return "", err
	}

	tp, err := utils.Float64(takePersent)
	if err != nil {
		return "", err
	}

	var tmp float64
	if tp > 1 {
		if positionSide == "LONG" {
			tmp = fpri * (1 + tp/100/ty.fleverage)
		} else if positionSide == "SHORT" {
			tmp = fpri * (1 - tp/100/ty.fleverage)
		}
	} else {
		if positionSide == "LONG" {
			tmp = fpri * (1 + tp/ty.fleverage)
		} else if positionSide == "SHORT" {
			tmp = fpri * (1 - tp/ty.fleverage)
		}
	}

	send := utils.String(tmp)
	return send, nil
}

// 가격으로 손절가격 가져오기
func (ty *BasicInfo) GetSLPrice(price, positionSide, takePersent string) (string, error) {
	fpri, err := utils.Float64(price)
	if err != nil {
		return "", err
	}

	tp, err := utils.Float64(takePersent)
	if err != nil {
		return "", err
	}

	var tmp float64
	if tp > 1 {
		if positionSide == "LONG" {
			tmp = fpri * (1 - tp/100/ty.fleverage)
		} else if positionSide == "SHORT" {
			tmp = fpri * (1 + tp/100/ty.fleverage)
		}
	} else {
		if positionSide == "LONG" {
			tmp = fpri * (1 - tp/ty.fleverage)
		} else if positionSide == "SHORT" {
			tmp = fpri * (1 + tp/ty.fleverage)
		}
	}

	send := utils.String(tmp)
	return send, nil
}

//////// 이벤트 영역 //////////

// 가격 금액 체크
func (ty *BasicInfo) CheckMinTraderPirce(price string, Amount string) (bool, error) {
	// price, Amount
	fprice, err := utils.Float64(price)
	if err != nil {
		return false, err
	}

	famount, err := utils.Float64(Amount)
	if err != nil {
		return false, err
	}

	tmp := fprice * famount * ty.fleverage
	if tmp < ty.fNotionalPrice {
		return false, fmt.Errorf("low Price ")
	}

	return true, nil
}

// 가격 최소 최대
func (ty *BasicInfo) CheckPriceMinMax(price string) (bool, error) {
	fprice, err := utils.Float64(price)
	if err != nil {
		return false, err
	}

	if ty.fMaxPrice < fprice {
		return false, fmt.Errorf("it is Over Price")
	}

	if ty.fMinPrice > fprice {
		return false, fmt.Errorf("it is Low Price")
	}

	return true, nil
}

// 가격을 넣으면 step 가격 틱사이즈로 자르기 - 가격으로 할때 소수점을 자르는 작업
func (ty *BasicInfo) ReturnPriceForSize(price string) (string, error) {

	decimalStr := utils.String(ty.fTickSizePrice)     //
	price, err := util.FormatPrice(price, decimalStr) // 소수점 처리 작업
	if err != nil {
		return "", err
	}

	return price, nil
}

// 최소 코인 수량
func (ty *BasicInfo) CheckMinAmount(Amount string) (bool, error) {
	famount, err := utils.Float64(Amount)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// 최소 수량
	if ty.fMinQuantity > famount {
		return false, fmt.Errorf("it is Low Amount")
	}

	return true, nil
}

// 최대 코인 수량
func (ty *BasicInfo) CheckMaxAmount(Amount string) (bool, error) {

	famount, err := utils.Float64(Amount)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// 최소 수량
	if ty.fMaxQuantity < famount {
		return false, fmt.Errorf("it is Over Amount")
	}

	return true, nil
}

// 주문 전 체크
func (ty *BasicInfo) BeforeOrder(currentPrice, amout string) (bool, error) {
	// currentPrice, Amount
	if currentPrice == "" {
		// 현재가를 기준으로 한다.

	}

	if amout == "" {
		// amount 제거
	}

	return true, nil

}

// 미체결 주문하기
// 심볼, 포지션, side, price, amount
/*
	symobl - 있으니까 패스
	position - LONG, SHORT
	side - buy,sell 대신, OPEN, CLOSE 로 변경
	price -
	amount -
*/
func (ty *BasicInfo) SendOpenOrder(position, side, price, amount string) (*futures.CreateOrderResponse, error) {
	// 사용하지는 않는데 만들어는 놓치뭐
	res, err := ty.mbacc.GetBinanceUser().SendOpenOrder(ty.symbol, position, side, price, amount)
	if err != nil {
		return nil, err
	}
	// 필요하다면 response 를 가공해도 됨 !!
	// (*futures.CreateOrderResponse, error)
	return res, nil
}

// 마켓주문
func (ty *BasicInfo) SendMarketOrder(symbol, position, openclose, price string) (*futures.CreateOrderResponse, error) {
	res, err := ty.mbacc.GetBinanceUser().SendOrderMarket(symbol, position, openclose, price)
	if err != nil {
		errors.Error("OnewayBot : SendOrderMarket", err)
		return nil, err
	}
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendTakeProfit ", log)
	return res, nil
}

// TL 주문
func (ty *BasicInfo) SendTakeProfit(symbol, position, openclose, price string) (*futures.CreateOrderResponse, error) {
	errors.Log("sendTakeProfit : ", symbol, position, openclose, price)
	res, err := ty.mbacc.GetBinanceUser().SendOrderTakeProfit(symbol, position, openclose, price)
	if err != nil {
		errors.Error("OnewayBot : SendOrderTakeProfit", err)
		return nil, err
	}
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendTakeProfit ", log)
	return res, nil
}

// SL 주문
func (ty *BasicInfo) SendStopLoss(symbol, position, openclose, price string) error {
	res, err := ty.mbacc.GetBinanceUser().SendOrderStopLoss(symbol, position, openclose, price)
	if err != nil {
		errors.Error("OnewayBot : sendStopLoss", err)
		return err
	}
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendStopMarket ", log)
	return nil
}

// Stop Market 주문
func (ty *BasicInfo) SendStopMarket(symbol, positionside, openclose, price, amount string) error {
	// 이렇게 돌아가는 이유는 비상시에 여기다가 로그를 남기거나 에러를 남길 수 있는 장치를 넣을 수 있음.
	res, err := ty.mbacc.GetBinanceUser().SendOrderStopMarket(symbol, positionside, openclose, price, amount)
	log := fmt.Sprintf("%+v\n", res)
	errors.Log("sendStopMarket ", log)
	return err
}

// 미체결 주문 취소
func (ty *BasicInfo) CancelOrder(symbol, ClientOrderID string) (*futures.CancelOrderResponse, error) {
	res, err := ty.mbacc.GetBinanceUser().CancelOrder(symbol, ClientOrderID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 미체결 모두 취소 - 실패한게 어떤건지 확인하기 !!
func (ty *BasicInfo) SendRemoveOpenOrder() []string {
	// 해당 심볼의 미체결리스트 가져와서
	// 싹다 취소 한번 돌리기 !!
	res := ty.mbacc.GetBinanceUser().SendRemoveOpenOrderForSymbol(ty.symbol)
	if len(res) != 0 {
		var send []string
		return send
	}
	return res
}

// 포지션 청산 - 2개, POSITION,
func (ty *BasicInfo) SendClosePosition(position string) error {
	// 포지션을 가져와서 !!
	res, err := ty.mbacc.GetBinanceUser().GetPositionRiskService(ty.symbol) // 전체 가져오기
	if err != nil {
		errors.Error("Crit Panic", "BinanceAccount.getOpenOrderList  ", err)
		return err
	}

	for _, v := range res {
		if v.Symbol == ty.symbol && v.PositionSide == position {
			// res, err := SendOrderMarket(symbol, position, openclose, amount)
			res, err := ty.mbacc.GetBinanceUser().SendOrderMarket(ty.symbol, position, "CLOSE", v.PositionAmt)
			if err != nil {
				return err
			}

			if res != nil {
				return fmt.Errorf("SendClosePosition response error")
			}
		}
	}
	return nil
}
