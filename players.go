// Package main is the entrypoint for fetch-hls
package main

import (
	"context"
	"strings"

	cast "github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/oleksandr/bonjour"
	"github.com/spf13/viper"
)

func playChromecast(url string) {
	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		Log.Fatal("Failed to initialize resolver", err)
	}

	var cc *bonjour.ServiceEntry
	results := make(chan *bonjour.ServiceEntry)

	err = resolver.Browse("_googlecast._tcp", "local.", results)
	if err != nil {
		Log.Fatal("Failed to browse", err)
	}

	playerName := strings.ToLower(viper.GetString("player-name"))
	for entry := range results {
		if playerName == "" {
			cc = entry
			break
		}

		deviceName := strings.ToLower(strings.Join(entry.Text, " "))
		if strings.Contains(deviceName, playerName) {
			cc = entry
			break
		}
	}

	ctx := context.Background()
	client := cast.NewClient(cc.AddrIPv4, 8009)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	media, err := client.Media(ctx)
	if err != nil {
		panic(err)
	}
	item := controllers.MediaItem{
		ContentId:   url,
		StreamType:  "BUFFERED",
		ContentType: "application/x-mpegurl",
	}

	_, err = media.LoadMedia(ctx, item, 0, true, map[string]interface{}{})
	if err != nil {
		Log.Warn(err)
	}
}
