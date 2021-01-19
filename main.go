// Package main is the entrypoint for fetch-hls
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version = "v0.0.0"
)

// Fetch contains the body of a fetch call copied from a browsers web developer tools.
type Fetch struct {
	Credentials string `json:"credentials"`
	Referrer    string `json:"referrer"`
	Method      string `json:"method"`
	Mode        string `json:"mode"`
	Headers     map[string]string
}

func proxyRequest(reqURL string, fetch *Fetch, c echo.Context) error {
	reqURLParsed, err := url.Parse(reqURL)
	if err != nil {
		return c.String(503, err.Error())
	}

	proxyReq, err := http.NewRequest(c.Request().Method, reqURL, c.Request().Body)
	if err != nil {
		return c.String(503, err.Error())
	}

	for key, val := range c.Request().Header {
		proxyReq.Header.Set(strings.ToLower(key), val[0])
	}
	for key, val := range fetch.Headers {
		proxyReq.Header.Set(strings.ToLower(key), val)
	}
	if fetch.Referrer != "" {
		proxyReq.Header.Set("referer", fetch.Referrer)
	}

	proxyReq.URL.RawQuery = c.Request().URL.RawQuery
	if reqURLParsed.RawQuery != "" {
		proxyReq.URL.RawQuery = reqURLParsed.RawQuery
	}

	httpClient := http.Client{}
	proxyResp, err := httpClient.Do(proxyReq)
	if err != nil {
		return c.String(503, err.Error())
	}

	c.Response().Status = proxyResp.StatusCode

	for key, val := range proxyResp.Header {
		c.Response().Header().Set(key, val[0])
	}

	_, err = io.Copy(c.Response().Writer, proxyResp.Body)
	if err != nil {
		return c.String(503, err.Error())
	}
	return proxyResp.Body.Close()
}

func run(rootCmd *cobra.Command, args []string) {
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	initLogger(viper.GetString("log-level"))
	if err != nil {
		Log.Debug(err)
	}

	stdinBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Log.Fatal(err)
	}

	var fetch Fetch
	fetchStr := string(stdinBytes)
	Log.Debug(fetchStr)
	err = json.Unmarshal([]byte(fetchStr[strings.Index(fetchStr, "{"):strings.LastIndex(fetchStr, "}")+1]), &fetch)
	if err != nil {
		Log.Fatal(err)
	}

	m3u8URL, err := url.Parse(strings.ReplaceAll(regexp.MustCompile(`"(.*?)"`).FindString(fetchStr), `"`, ""))
	if err != nil {
		Log.Fatal(err)
	}

	localStreamURL := "http://" + viper.GetString("ip") + ":" + viper.GetString("port") + m3u8URL.Path
	streamHost := m3u8URL.Scheme + "://" + m3u8URL.Hostname()
	if m3u8URL.Port() != "80" && m3u8URL.Port() != "443" {
		streamHost += ":" + m3u8URL.Port()
	}

	playlistPath := m3u8URL.Path
	if m3u8URL.RawQuery != "" {
		playlistPath += "?" + m3u8URL.RawQuery
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Up")
	})

	e.GET(m3u8URL.Path, func(c echo.Context) error {
		return proxyRequest(streamHost+playlistPath, &fetch, c)
	})

	e.GET("/*", func(c echo.Context) error {
		return proxyRequest(streamHost+c.Request().URL.Path, &fetch, c)
	})

	go func() {
		fmt.Println("Stream available at " + localStreamURL)
		e.Logger.Fatal(e.Start(":" + viper.GetString("port")))
	}()

	if viper.GetString("player") == "chromecast" {
		playChromecast(localStreamURL)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for range c {
		os.Exit(0)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "fetch-hls",
		Version: version,
		Example: `  pbpaste | fetch-hls
  cat fetch.js | fetch-hls --player chromecast`,
		Run:   run,
		Short: "A quick and lazy solution to proxy HLS streams to external players (Chromecast, VLC) when the stream itself has some odd authentications through either query parameters or HTTP headers, which by some external players will ignore or not have access to.",
	}

	flags := rootCmd.PersistentFlags()
	flags.StringP("player", "p", "http", "Player to use. Accepts 'http' and 'chromecast'.")
	flags.StringP("ip", "i", getLocalIP(), "Local IP address for HTTP server.")
	flags.String("port", "8899", "Port for HTTP server.")
	flags.String("log-level", "info", "Log level.")

	viper.SetConfigName("fetch-hls")
	viper.SetEnvPrefix("fetch-hls")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetConfigType("json")

	flags.VisitAll(func(f *pflag.Flag) {
		err := viper.BindPFlag(f.Name, flags.Lookup(f.Name))
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
