package service

type UserInfo struct {
	// 계좌정보 - Market
	// 봇정보 - OnewayBot
	UserID    string // 유저 아이디
	AccessKey string // key
	SecritKey string // key

	BotObject map[string]InterfaceBot // key 종목+포지션
	// 웹소켓 연결 하기
	/*
		체결시 broadcase로 각각 bot 에 전달하기
	*/
	WSObject BinanceUserWSObject
}

// 초기화
func (ty *UserInfo) Init(AccessKey, SecritKey string) {
	ty.AccessKey = AccessKey
	ty.SecritKey = SecritKey
	ty.BotObject = make(map[string]InterfaceBot)
}

// 웹소켓 연결하기
func (ty *UserInfo) WSConnection() {

	ty.onConnect()
	ty.onUnConnect()
	ty.onMessage()

}

// 연결시
func (ty *UserInfo) onConnect() {

}

// 해제시
func (ty *UserInfo) onUnConnect() {

}

// 메시지 전달시
func (ty *UserInfo) onMessage() {

}
