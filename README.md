# fetch-hls

A *quick and lazy* solution to proxy HLS streams to external players (Chromecast, VLC) when the stream itself has some odd authentications through either query parameters or HTTP headers, which by some external players will ignore or not have access to. Proxy is initialized by copying a successful request from your Browsers web developer tools, and passing to the application through stdin.

## Example

Using Apple's HLS examples:

1. Open URL in browser: https://developer.apple.com/streaming/examples/basic-stream-osx-ios4-3.html
2. Open web developer tools and find the initial `*.m3u8` request.
3. Right click on the request and `Copy as Fetch` ([Screenshot](https://i.imgur.com/FYl2Ovx.png))
4. `pbpaste | fetch-hls`

## Installation


**homebrew** (OSX / Linux):

```sh
brew install dustinblackman/tap/fetch-hls
```

**scoop** (Windows):

```sh
$ scoop bucket add dustinblackman https://github.com/dustinblackman/scoop-bucket.git
$ scoop install fetch-hls
```

**deb/rpm** (Linux):

Download the `.deb` or `.rpm` from the [releases page](https://github.com/dustinblackman/fetch-hls/releases) and
install with `dpkg -i` and `rpm -i` respectively.


**manually**:

Download the pre-compiled binaries from the [releases page](https://github.com/dustinblackman/fetch-hls/releases) and
copy to the desired location.

**go/master branch:**

```
go get -u github.com/dustinblackman/fetch-hls
```

## Usage

```
HLS proxy that extracts m3u8 playlist context from Google/Firefox's web dev tools 'Get as Fetch' function on an m3u8 request, guaranteeing all request information provided in browser to allow the request to succeed will be accessible from external devices like a Chromecast.

Usage:
  fetch-hls [flags]

Examples:
  pbpaste | fetch-hls
  cat fetch.js | fetch-hls --player chromecast

Flags:
  -h, --help               help for fetch-hls
  -i, --ip string          Local IP address for HTTP server. (default "192.168.1.10")
      --log-level string   Log level. (default "info")
  -p, --player string      Player to use. Accepts 'http' and 'chromecast'. (default "http")
      --port string        Port for HTTP server. (default "8899")
```

## [License](./LICENSE)

MIT
