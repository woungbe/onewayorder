package service

import (
	"fmt"
	"onewayorder/binance"
	"onewayorder/errors"
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

	MinNotion string // 최소 금액 USDT
	MinSize   string // 코인 최소 수량 ex)
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

// 최소단위, 코인 소수점 가져오기
func (ty *BasicInfo) SetExChangeInfo() (string, string, error) {
	for _, v := range ty.mbacc.ExchangeInfo {
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

// 가격으로 익절가격 가져오기
func (ty *BasicInfo) GetTPPrice(price, takePersent string) (string, error) {
	// 가격, 익절가격

	return "", nil
}

// 가격으로 손절가격 가져오기
func (ty *BasicInfo) GetSLPrice(price, takePersent string) (string, error) {

	return "", nil
}

//////// 이벤트 영역 //////////

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
		// tmp := new(UserPositionInfo)
		// key, posinfo := tmp.SetPosition(v)
		// symbol + "_" + posside  ex) NOTUSDT_SHORT
		// ty.mUserPositions[key] = posinfo
		if v.Symbol == ty.symbol && v.PositionSide == position {
			// res, err := SendOrderMarket(symbol, position, openclose, amount)
			res, err := ty.mbacc.GetBinanceUser().SendOrderMarket(ty.symbol, position, "CLOSE", v.PositionAmt)
			if err != nil {
				return err
			}

			if res != nil {
				return fmt.Errorf("SendClosePosition response error")
			}

			// 뭔가 response 에서 에러를 줄 것 같은 느낌인데...
		}
	}
	return nil
}
