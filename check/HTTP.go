package check

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func ProxyOfhttp(in string, in2 int) bool {

	proxy := "http://" + in
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Println("parse url error: ", err)
		return false
	}

	client := &http.Client{
		Timeout: time.Duration(in2) * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
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
