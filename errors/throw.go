package errors

import (
	"fmt"
	"os"
	"time"
)

// 일반 알림 로그
func Log(aa ...interface{}) {
	send := fmt.Sprint(aa...)
	err := fmt.Errorf(send)
	logErrorToFile("[LOG] ", err)
}

// Error 알림 로그
func Error(aa ...interface{}) {
	send := fmt.Sprint(aa...)
	err := fmt.Errorf(send)
	logErrorToFile("[Error] ", err)
}

// Error 알림 로그
func ReturnError(aa ...interface{}) error {
	Error(aa...)
	tmp := fmt.Sprintln(aa...) // tmp
	return fmt.Errorf(tmp)     // Errorf
}

func logErrorToFile(tag string, err error) {
	// 현재 날짜와 시간을 포맷하여 문자열로 변환
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")

	// 에러 메시지를 포맷
	logMessage := fmt.Sprintf("[%s] [%s] %s\n", currentDateTime, tag, err.Error())

	// 로그 파일 이름을 날짜 기반으로 생성
	logFileName := "./logs/error_log" + time.Now().Format("2006-01-02") + ".txt"

	// 파일을 append 모드로 열기, 파일이 없으면 생성
	file, fileErr := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		fmt.Printf("Failed to open log file: %s\n", fileErr.Error())
		return
	}
	defer file.Close()

	// 에러 메시지를 파일에 기록
	if _, writeErr := file.WriteString(logMessage); writeErr != nil {
		fmt.Printf("Failed to write to log file: %s\n", writeErr.Error())
	}
}
