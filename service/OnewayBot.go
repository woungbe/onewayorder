package service

type InterfaceBot interface {
}

type OnewayBot struct {
	// AccessKey string
	// SecritKey string
	mBotName string

	mConfigData BotConfigData // 설정 정보 데이터
	mProcData   ProcData      // 봇 실행 정보
}

// 설정 정보 데이터
type BotConfigData struct {
}

// 봇 실행 정보
type ProcData struct {
}

// Init(userListenKey string)
