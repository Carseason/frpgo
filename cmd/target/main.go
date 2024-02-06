package main

import (
	"frpgo/pkgs/config"
	"frpgo/pkgs/sdk"
)

func main() {
	err := sdk.NewFrpVisitor(sdk.Config{
		ServerAddr: config.SERVER_ADDR,
		ServerPort: config.SERVER_PORT,
		AuthToken:  config.AUTH_TOKEN,
		User:       config.USER_NAME,
		UserToken:  config.USER_TOKEN,
	}, []sdk.VisitorConfig{
		{
			ServerName: config.DEVIVE_NAME,
			SecretKey:  config.DEVICE_SECRETKEY,
			LocalIP:    "0.0.0.0",
			LocalPort:  9002,
		},
		{
			ServerName: config.DEVIVE_NAME1,
			SecretKey:  config.DEVICE_SECRETKEY,
			LocalIP:    "0.0.0.0",
			LocalPort:  9003,
		},
	})
	if err != nil {
		panic(err)
	}
}
