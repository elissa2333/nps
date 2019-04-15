package check

import (
	"bytes"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func ProxyOfsocks5(in string, in2 int) bool {
	tou := time.Duration(in2)

	dialer, err := proxy.SOCKS5("tcp", in, nil,
		&net.Dialer{
			Timeout:   tou * time.Second,
			KeepAlive: tou * time.Second,
		},
	)
	if err != nil {

		return false
	}

	httpTransport := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{
		Transport: httpTransport,
		Timeout:   tou * time.Second,
	}

	req, _ := http.NewRequest("GET", "http://ip.sb", nil)
	req.Header.Set("User-Agent", "curl/7.52.1")

	res, err := client.Do(req)
	if err != nil {

		return false
	}
	body, _ := ioutil.ReadAll(res.Body)

	res.Body.Close()

	var foo []byte

	cache := bytes.Replace(body, []byte("\n"), foo, 1)

	cache2 := net.ParseIP(string(cache))
	if cache2 != nil {
		return true
	}

	return false
}
