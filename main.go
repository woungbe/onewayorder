package main

import (
	"onewayorder/errors"
	"onewayorder/service"
	util "onewayorder/util"
)

func main() {
	util.RandInitSeed()

	service.GetController().Init()

	done := make(chan bool)
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
