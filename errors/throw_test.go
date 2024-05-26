package errors

import "testing"

func TestLog(t *testing.T) {
	errmsg := "알릴 필요성이 있어서 알림 "
	Log(errmsg)
}

func TestError(t *testing.T) {
	errmsg := "에러가 발생했음"
	Error(errmsg)
}
