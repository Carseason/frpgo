package main

import (
	"context"
	"fmt"
	"routergo/pkgs/config"
	frpplugins "routergo/pkgs/frpPlugins"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/server"
)

func main() {
	if err := runApp(); err != nil {
		panic(err)
	}
}
func runApp() error {
	tcpMux := true
	detailedErrorsToClient := true
	config := &v1.ServerConfig{
		Log: v1.LogConfig{
			// 写入文件,console为 stdout
			To:      "console",
			Level:   "info",
			MaxDays: 3,
		},
		Auth: v1.AuthServerConfig{
			Method: v1.AuthMethodToken,
			Token:  config.AUTH_TOKEN,
		},
		BindAddr:         config.SERVER_ADDR,
		BindPort:         config.SERVER_PORT,
		ProxyBindAddr:    "0.0.0.0",
		VhostHTTPTimeout: 60,
		Transport: v1.ServerTransportConfig{
			TCPMux:                  &tcpMux,
			TCPMuxKeepaliveInterval: 60,
			TCPKeepAlive:            7200,
			MaxPoolCount:            5,
			HeartbeatTimeout:        90,
			QUIC: v1.QUICOptions{
				KeepalivePeriod:    10,
				MaxIdleTimeout:     30,
				MaxIncomingStreams: 100000,
			},
		},
		DetailedErrorsToClient:          &detailedErrorsToClient,
		UserConnTimeout:                 10,
		UDPPacketSize:                   1500,
		NatHoleAnalysisDataReserveHours: 168,
		HTTPPlugins:                     frpplugins.Plugins,
	}
	frpsServer, err := server.NewService(config)
	if err != nil {
		return err
	}
	frpsServer.Run(context.Background())
	fmt.Println("frp server started")
	return nil
}
