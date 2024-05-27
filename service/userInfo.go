package service

type UserInfo struct {
	// 계좌정보 - Market
	// 봇정보 - OnewayBot
	UserID    string // 유저 아이디
	AccessKey string // key
	SecritKey string // key

	BotObject map[string]InterfaceBot // key 종목+포지션 =
}

// 초기화
func (ty *UserInfo) Init(AccessKey, SecritKey string) {
	ty.AccessKey = AccessKey
	ty.SecritKey = SecritKey
	ty.BotObject = make(map[string]InterfaceBot)
}

///////// 여기에는 각각 필요한 사람 /////////
///////// 뭐가 필요할지 생각하기 ///////
