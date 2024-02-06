package main

import (
	"frpgo/pkgs/config"
	"frpgo/pkgs/sdk"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 启动两个客户端
	// device
	go func() {
		err := sdk.NewFrpProxy(sdk.Config{
			ServerAddr: config.SERVER_ADDR,
			ServerPort: config.SERVER_PORT,
			AuthToken:  config.AUTH_TOKEN,
			User:       config.USER_NAME,
			UserToken:  config.USER_TOKEN,
		}, []sdk.ProxyConfig{
			{
				ServerName: config.DEVIVE_NAME,
				SecretKey:  config.DEVICE_SECRETKEY,
				LocalIP:    "192.168.100.127",
				LocalPort:  8301,
			},
		})
		if err != nil {
			panic(err)
		}
	}()
	// device1
	go func() {
		err := sdk.NewFrpProxy(sdk.Config{
			ServerAddr: config.SERVER_ADDR,
			ServerPort: config.SERVER_PORT,
			AuthToken:  config.AUTH_TOKEN,
			User:       config.USER_NAME,
			UserToken:  config.USER_TOKEN,
		}, []sdk.ProxyConfig{
			{
				ServerName: config.DEVIVE_NAME1,
				SecretKey:  config.DEVICE_SECRETKEY,
				LocalIP:    "192.168.100.127",
				LocalPort:  9910,
			},
		})
		if err != nil {
			panic(err)
		}
	}()
	wg.Wait()

}
