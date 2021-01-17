// Package main is the entrypoint for fetch-hls
package main

import "go.uber.org/zap"

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
