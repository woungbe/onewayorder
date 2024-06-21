package service

import (
	"fmt"
	"math/rand"
	"onewayorder/binance"
	"onewayorder/errors"
	"onewayorder/utils"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

type callBackFunc_onStatusChange func(nStatus int)     //상태 변경 콜백 함수
type callBackFunc_onAccountErrorClose func(msg string) //상태 오류로 인한 계좌정보 서비스 정지

type callBackFunc_onLeverage func(Symbol string, Leverage int)                                                                                                                 //레버리지 변경 이벤트
type callBackFunc_onMarginTypeChange func(Symbol string, Side string, nIsolated int /*1=격리, 0= 크로스*/)                                                                          //마진 타입 변경
type callBackFunc_onPositionClose func(Symbol string, Side string)                                                                                                             //포지션 청산
type callBackFunc_onPositionUpdate func(Symbol string, Side string, MarginType bool, Leverage int, AvgPrice string, IsolatedWallet string, Amount string, changeAvgPrice bool) //포지션 정보 업데이트
type callBackFunc_onMarginTransfer func(Symbol string, Side string, margin float64)                                                                                            //포지션 마진 변경

type callBackFunc_onAddOpenOrder func(Symbol string, PositionSide string, Side string, OrderType string, ClientOrderID string)  //신규 주문 추가 (시장가 강제청산 주문 제외)
type callBackFunc_onCanceledOrder func(Symbol string, PositionSide string, Side string, OrderType string, ClientOrderID string) //주뭄 취소
type callBackFunc_onExpiredOrder func(Symbol string, PositionSide string, Side string, OrderType string, ClientOrderID string)  //주문 만료(만료 주문)

// = TradeType = PARTIALLY_FILLED : 부분체결, FILLED : 전체체결완료 , LIQUIDATION : 강제청산 , closeOrder = true=청산주문,false 오픈 주문 , MarginType= true : 격리
type callBackFunc_onTradeOrder func(ClientOrderID string, Symbol string, PositionSide string, Side string, MarginType bool, Leverage int, TradeType string, closeOrder bool, AvgPrice string, EntryPrice string, tradeQty string, Commission string, RealizedPnL string, TradeTime int64) //체결 이벤트

type BinanceAccount struct {
	mApikey      string // 유저 API 키
	mApiSeckey   string // 유저 시크릿키
	mListenkey   string // 웹소켓 리슨키
	ExchangeInfo []futures.Symbol

	mUserBalance      string                          // USDT 발란스
	mUserOpenOrders   map[int64]UserOpenOrder         // 오픈정보
	mUserPositions    map[string]UserPositionInfo     // 포지션
	mUserLeverageInfo map[string]UserPairLeverageInfo // 레버리지 및 마진타입 정보

	bReconnecting      bool  // 자동 재접속중인가.
	bAutoReconnectMode bool  // true = 자동 연결 사용, false= 자동연결 미사용
	mWSReConnectTime   int64 // 웹소켓 재접속 시간 (ms)
	mBreakListenChk    bool  // 리스킨 체크
	mLastLintenKeyTime int64 // 마지막 웹소켓 리슨키 가져온시간 (ms)

	mBinanceUserWS *BinanceUserWSObject // 바이낸스유저 웹소켓

	mProcStatus int // 섹션 상태 정보 ( 0= 사용정지, 1=정상 )

	// mBinanceAPI *futures.Client // 바이낸스 API 인스턴스
	mBinanceAPI *binance.BinanceUser

	mCB_ChangeStatus      callBackFunc_onStatusChange      //상태 변경 콜백 함수
	mCB_AccountErrorClose callBackFunc_onAccountErrorClose // 계좌정보 로드및 상태 오류로 인한 웹소켓 정지
	mCD_Leverage          callBackFunc_onLeverage          //레버리지 변경 이벤트
	mCD_MarginTypeChange  callBackFunc_onMarginTypeChange  //마진 타입 변경
	mCD_PositionClose     callBackFunc_onPositionClose     //포지션 청산
	mCD_PositionUpdate    callBackFunc_onPositionUpdate    //포지션 정보 업뎃
	mCD_MarginTransfer    callBackFunc_onMarginTransfer    //포지션 마진 변경
	mCD_AddOpenOrder      callBackFunc_onAddOpenOrder      //주문 추가
	mCD_CanceledOrder     callBackFunc_onCanceledOrder     //주문 취소(완료)
	mCD_ExpiredOrder      callBackFunc_onExpiredOrder      //주문 만료
	mCD_TradeOrder        callBackFunc_onTradeOrder        //체결이벤트

	mBotObject []InterfaceBot // 봇 리스트 정리 - 롱봇 숏봇 상관없지 ?...
}

func (ty *BinanceAccount) Init(ApiKey, ApiSeceryKey string) {
	ty.mApikey = ApiKey
	ty.mApiSeckey = ApiSeceryKey
	ty.mListenkey = ""

	ty.mUserOpenOrders = make(map[int64]UserOpenOrder)
	ty.mUserPositions = make(map[string]UserPositionInfo)
	ty.mUserLeverageInfo = make(map[string]UserPairLeverageInfo)

	ty.bReconnecting = true      // 자동 재접속중
	ty.bAutoReconnectMode = true // 자동 연결 사용
	ty.mWSReConnectTime = 3000   // 웹소켓 재접속 시간
	ty.mLastLintenKeyTime = 5000 // 마지막 웹소켓 리슨키 가져오는

	// 일단 신청해서 mBinanceAPI 하자 !!
	ty.mBinanceAPI = new(binance.BinanceUser)
	ty.mBinanceUserWS = new(BinanceUserWSObject)

	// 여기서 해야되는 것들 정리
	// err := ty.GetWSKey()

	b := ty.LoadAccountInfo()
	if !b {
		errors.Error("Crit Panic", "Init - ty.LoadAccountInfo  ", b)
		return
	}
}

// 벨런스 가져오기
func (ty *BinanceAccount) GetUserBalance() string {
	return ty.mUserBalance
}

// 오픈오더 가져오기
func (ty *BinanceAccount) GetUserOpenOrders() map[int64]UserOpenOrder {
	return ty.mUserOpenOrders
}

// 포지션 가져오기
func (ty *BinanceAccount) GetUserPositions() map[string]UserPositionInfo {
	return ty.mUserPositions
}

// 레버리지 가져오기
func (ty *BinanceAccount) GetUserLeverageInfo() map[string]UserPairLeverageInfo {
	return ty.mUserLeverageInfo
}

func (ty *BinanceAccount) SetBotList(args map[string]interface{}) {
	if val, ok := args["BotName"].(string); ok {
		if val == "OnewayBot" {
			tmp := new(OnewayBot)
			b, err := tmp.SetConfigData(args)
			if err != nil {
				fmt.Println(err)
			}
			if b {
				ty.mBotObject = append(ty.mBotObject, tmp)
			}
		}
	}
}

// getWSkey
func (ty *BinanceAccount) GetWSKey() bool {
	listenkey, err := ty.mBinanceAPI.GetStartUserStreamService()
	if err != nil {
		errors.Error()
		return true
	}
	ty.mListenkey = listenkey
	return false
}

// GetWS - Connection
func (ty *BinanceAccount) GetWSConnection() {
	ty.mBinanceUserWS.Init(ty.mListenkey)
}

// 바이낸스 유저 웹소켓 최소 연결
func (ty *BinanceAccount) StrartUserWS() bool {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit Panic", "BinanceAccount.StrartUserWS ", err)
		}
	}()

	ty.bAutoReconnectMode = true
	ty.bReconnecting = false
	ty.mBinanceUserWS = nil

	newObj := new(BinanceUserWSObject)
	newObj.Init(ty.mListenkey)
	newObj.SetCallbackFunc(ty.onConnected, ty.onUnconnected, ty.messagePasering)

	b := newObj.ClientConnect()
	if !b {
		ty.setStatus(-1)
		return false
	}

	//연결 성공시 처리부
	ty.mWSReConnectTime = 0
	ty.mBreakListenChk = true
	ty.mBinanceUserWS = newObj

	go ty.checkListenkey()
	ty.setStatus(1)

	return true
}

func (ty *BinanceAccount) onConnected() {

}
func (ty *BinanceAccount) onUnconnected() {

}

func (ty *BinanceAccount) messagePasering(msg []byte) {

}

// exchangeInfo
func (ty *BinanceAccount) GetExchangeInfo() error {
	// ExchangeInfo []futures.Symbol

	exchangeInfo, err := binance.GetExchangeInfo()
	if err != nil {
		return err
	}

	ty.ExchangeInfo = exchangeInfo.Symbols
	return nil
}

// 객체 상태 정보
func (ty *BinanceAccount) setStatus(nStatus int) {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit defer", "BinanceAccount.setStatus - ", nStatus)
		}
	}()

	if ty.mProcStatus != nStatus {
		ty.mProcStatus = nStatus

		if ty.mCB_ChangeStatus != nil {
			ty.mCB_ChangeStatus(ty.mProcStatus)
		}
		//== 이벤트 발생
	}
}

// 웹소켓 연결된 상태에서는 10분마다 연장 체크를 하도록 한다.
func (ty *BinanceAccount) checkListenkey() {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit defer", "BinanceAccount.checkListenkey - ", err)
		}
	}()

	var chktime int64
	rnd := rand.New(utils.GetRandSeed())
	rndnum := int64(rnd.Intn(600) + 600)

	chktime = rndnum
	for {
		if !ty.mBreakListenChk {
			return
		}

		if ty.mBinanceUserWS != nil && ty.mBinanceUserWS.IsConnect() {
			curTime := utils.GetCurrentTimestamp()
			if curTime > chktime+ty.mLastLintenKeyTime {

				rnd := rand.New(utils.GetRandSeed())
				rndnum := int64(rnd.Intn(600) + 600)
				chktime = rndnum
				ty.procListenKey()
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// 리슨키 가져오기
func (ty *BinanceAccount) procListenKey() {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit Panic", "BinanceAccount.procListenKey  ", err)
		}
	}()

	start := time.Now()
	listenKey, er := ty.mBinanceAPI.GetStartUserStreamService()
	if er != nil {
		//오류
		endtime := time.Since(start)
		et := endtime.Milliseconds()
		mm := fmt.Sprintf("BinanceApi error(KeepAliveListenKey) 경과시간=%d", et)
		errors.Error("Crit Panic", "BinanceUserWSObject.Init - ClientConnect ", mm)
		return
	}

	ty.mLastLintenKeyTime = utils.GetCurrentTimestamp()
	ty.mListenkey = listenKey

	if ty.mListenkey != listenKey {
		ty.mListenkey = listenKey
		//이전 리슨키와 다르다면 !@!!  이전 웹소켓 닫고 새로 연결 처리
		if ty.mBinanceUserWS != nil {
			ty.mBinanceUserWS.SetCallbackFunc(nil, nil, nil)
			ty.mBinanceUserWS.Close()
		}
		ty.mBinanceUserWS = nil

		ty.mWSReConnectTime = utils.MakeTimestamp()
		ty.mWSReConnectTime = ty.mWSReConnectTime - 10000

		ty.reconnectingUserWS()
	}
}

// reconnectingUserWS 바이넨스 유저웹소켓 재접속처리부
func (ty *BinanceAccount) reconnectingUserWS() {
	defer func() {
		ty.bReconnecting = false
		ty.mWSReConnectTime = 0
		if err := recover(); err != nil {
			errors.Error("Crit defer", "BinanceUserWSObject.Init - ", err)
		}

	}()
	ty.bReconnecting = true

	//== 서버가 연결이 해제된 경우이기때문에 바이넨스 웹소켓에 연결된후 발란스,포지션,미체결 정보를 다시 리셋해야한다.
	// msg := fmt.Sprintf("-b- Try Reconnect Binance UnConnected = %s", ty.mApikey)
	newObj := new(BinanceUserWSObject)
	newObj.Init(ty.mListenkey)
	newObj.SetCallbackFunc(ty.onConnected, ty.onUnconnected, ty.messagePasering)
	ty.mBinanceUserWS = newObj

	for {
		if !ty.bAutoReconnectMode {
			ty.setStatus(0)
			return
		}

		if ty.mBinanceUserWS.ClientConnect() {
			curTime := utils.MakeTimestamp()
			if ty.mWSReConnectTime == 0 || (curTime-ty.mWSReConnectTime) >= 3000 {
				bAcc := ty.getAccountInfo()
				if !bAcc {
					if ty.mCB_AccountErrorClose != nil {
						ty.mCB_AccountErrorClose("계좌정보 로드 오류")
					}
					ty.setStatus(-2)
					return
				}
				bOp := ty.GetOpenOrderList()
				if !bOp {
					if ty.mCB_AccountErrorClose != nil {
						ty.mCB_AccountErrorClose("미체결정보 로드 오류")
					}
					ty.setStatus(-2)
					return
				}
			}
			ty.setStatus(1)
			return
		}
		time.Sleep(time.Duration(time.Millisecond * 300))
	}
}

// 사용자 계좌 정보(포지션,미체결,계좌 발란스및 상세 정보 로드)
func (ty *BinanceAccount) LoadAccountInfo() bool {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit defer", "BinanceUserWSObject.Init - ", err)
		}
	}()

	//-- 레버리지 버킷 리스트 정보가 있는가 체크
	bL := ty.GetWSKey() //리슨키를 못가져온경우 (계좌 초기화 오류 처리 )
	if !bL {
		ty.setStatus(-1)
		return false
	}
	bAcc := ty.getAccountInfo()
	if !bAcc {
		ty.setStatus(-1)
		return false
	}
	bOp := ty.GetOpenOrderList()
	if !bOp {
		ty.setStatus(-1)
		return false
	}
	ty.setStatus(100)

	return true
}

// 어카운트 가져오기
func (ty *BinanceAccount) getAccountInfo() bool {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit Panic", "getAccountInfo  ", err)
		}
	}()

	balance, err := ty.mBinanceAPI.GetAvailableBalance("USDT")
	if err != nil {
		errors.Error("Crit Panic", "BinanceAccount.getAccountInfo  ", err)
		return true
	}
	ty.mUserBalance = balance // 이용 가능 자산 !!
	return false
}

// 오픈 리스트 가져오기
func (ty *BinanceAccount) GetOpenOrderList() bool {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit Panic", "getOpenOrderList  ", err)
		}
	}()

	// 오픈오더 가져오기
	res, err := ty.mBinanceAPI.GetListOpenOrdersService("") // 전체 가져오기
	if err != nil {
		errors.Error("Crit Panic", "BinanceAccount.getOpenOrderList  ", err)
		return false
	}

	// 초기화 시키고 등록 - set하는 것때문에 이렇게 함. mutax 등록하는 편이 좋음.
	ty.mUserOpenOrders = make(map[int64]UserOpenOrder)
	for _, v := range res {
		tmpOrder := new(UserOpenOrder)
		ty.mUserOpenOrders[v.OrderID] = tmpOrder.SetOpenOrder(*v)
	}

	return false
}

// 포지션
func (ty *BinanceAccount) GetPositionList() bool {
	defer func() {
		if err := recover(); err != nil {
			errors.Error("Crit Panic", "getOpenOrderList  ", err)
		}
	}()

	res, err := ty.mBinanceAPI.GetPositionRiskService("") // 전체 가져오기
	if err != nil {
		errors.Error("Crit Panic", "BinanceAccount.getOpenOrderList  ", err)
		return false
	}

	for _, v := range res {
		tmp := new(UserPositionInfo)
		key, posinfo := tmp.SetPosition(v)
		// symbol + "_" + posside  ex) NOTUSDT_SHORT
		ty.mUserPositions[key] = posinfo
	}

	return false
}

/*
	콜백 함수 설정
*/

func (ty *BinanceAccount) SetCallbackFunc_StatusChange(cb callBackFunc_onStatusChange) {
	ty.mCB_ChangeStatus = cb
}

func (ty *BinanceAccount) SetCallbackFunc_AccountErrorClose(cb callBackFunc_onAccountErrorClose) {
	ty.mCB_AccountErrorClose = cb
}

func (ty *BinanceAccount) SetCallbackFunc_Leverage(cb callBackFunc_onLeverage) {
	ty.mCD_Leverage = cb
}

// 마진 타입 변경
func (ty *BinanceAccount) SetCallbackFunc_MarginTypeChange(cb callBackFunc_onMarginTypeChange) {
	ty.mCD_MarginTypeChange = cb
}

// 포지션 청산
func (ty *BinanceAccount) SetCallbackFunc_PositionClose(cb callBackFunc_onPositionClose) {
	ty.mCD_PositionClose = cb
}

// 포지션 정보 업뎃
func (ty *BinanceAccount) SetCallbackFunc_PositionUpdate(cb callBackFunc_onPositionUpdate) {
	ty.mCD_PositionUpdate = cb
}

// 포지션 마진 변경
func (ty *BinanceAccount) SetCallbackFunc_MarginTransfer(cb callBackFunc_onMarginTransfer) {
	ty.mCD_MarginTransfer = cb
}

// 주문 추가
func (ty *BinanceAccount) SetCallbackFunc_AddOpenOrder(cb callBackFunc_onAddOpenOrder) {
	ty.mCD_AddOpenOrder = cb
}

// 주문 취소(완료)
func (ty *BinanceAccount) SetCallbackFunc_CanceledOrder(cb callBackFunc_onCanceledOrder) {
	ty.mCD_CanceledOrder = cb
}

// 주문 만료
func (ty *BinanceAccount) SetCallbackFunc_ExpiredOrder(cb callBackFunc_onExpiredOrder) {
	ty.mCD_ExpiredOrder = cb
}

// 체결이벤트
func (ty *BinanceAccount) SetCallbackFunc_TradeOrder(cb callBackFunc_onTradeOrder) {
	ty.mCD_TradeOrder = cb
}
