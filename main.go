package main

import (
	"onewayorder/errors"
	"onewayorder/utils"
)

func main() {

	utils.RandInitSeed()

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
