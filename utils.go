// Package main is the entrypoint for fetch-hls
package main

import (
	"net"

	"go.uber.org/zap"
)

// Log is the initialized logger for the app.
var Log *zap.SugaredLogger

func initLogger(level string) {
	config := zap.NewProductionConfig()

	if level != "" {
		err := config.Level.UnmarshalText([]byte(level))
		if err != nil {
			panic(err)
		}
	}

	initLogger, _ := config.Build()
	Log = initLogger.Sugar()
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	//nolint
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
