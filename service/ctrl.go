package service

import (
	"fmt"
	"onewayorder/errors"
	"sync"
)

// 인터페이스도 정리를 해야되는 것임
// 안쪽은 알필요도 없음 - ctrl 이랑 연결하고 나머지는 알아서 처리함.
/*
	신규 유저 등록
	신규 유저 삭제
	유저 봇 시작
	유저 봇 정지
	유저 봇 컨트롤 여러개 ( 반복정지 - 다음부터 진행 .. 등등 여러개 가능함 )
*/

type Controller struct {
	mUserMap map[string]*UserInfo
}

var (
	CtrlInstance *Controller
	once         sync.Once
)

func GetController() *Controller {
	once.Do(func() {
		// fmt.Println("Creating Singleton instance")
		// instance = &Singleton{Value: "My Singleton"}
		CtrlInstance = new(Controller)
	})
	return CtrlInstance
}

// 처음
func (ty *Controller) Init() {
	ty.mUserMap = make(map[string]*UserInfo)
}

// 유저 가져오기 - 상태 체크 - 로그 체크
func (ty *Controller) GetUserMap(userid string) *UserInfo {
	if val, ok := ty.mUserMap[userid]; ok {
		return val
	}
	return nil
}

// 유저저장하기 - 키는 등록해 놓는것
func (ty *Controller) SetUserMap(userid, acckey, seckey string) error {
	muserinfo := new(UserInfo)
	muserinfo.Init(acckey, seckey)
	if val, ok := ty.mUserMap[userid]; ok {
		return errors.ReturnError("이미 유저가 있습니다. ", userid, val.AccessKey)
	}
	ty.mUserMap[userid] = muserinfo
	return nil
}

// 시작하기 위한 봇 설정 정보
type BotConfigStruct struct {
	userID    string
	BotName   string
	BotConfig interface{}
}

// 유저 봇 시작 - 봇을 시작하기 위해서 어떤 값들이 필요할까 .
func (ty *Controller) StartBot(args BotConfigStruct) error {
	// 인터페이스
	// ty.mUserMap[]
	/*
		if v, ok := ty.mUserMap[args.userID]; ok {
			if args.BotName == "UpDown" {
				v.BotObject = args.BotConfig.(BotConfigData)
			}

		}
	*/

	return nil
}

// 유저 봇 정지 - 봇을 정지하기 위해서 어떤 값들이 필요할까?
func (ty *Controller) StopBot(userid string, aa interface{}) error {

	bot := ty.mUserMap[userid]

	for k, v := range bot.BotObject {
		fmt.Println(k, v)
	}

	return nil
}
