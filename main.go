package main

import "onewayorder/errors"

func main() {

	done := make(chan bool)

	errmsg := "알릴 필요성이 있어서 알림 "
	errors.Log(errmsg)

	errmsg2 := "에러가 발생했음"
	errors.Error(errmsg2)

	<-done

}

func CMDPaser(strCMD string) {
	defer func() {
		if err := recover(); err != nil {
			errors.Error(err)
		}
	}()
	if strCMD == "" {
		return
	}

}
