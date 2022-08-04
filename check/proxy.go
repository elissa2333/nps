package check

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func IsProxy(proxy string, timeout time.Duration) bool {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Println("parse proxy url error: ", err)
		return false
	}

	client := &http.Client{
		Timeout: timeout,
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

	body, _ := io.ReadAll(res.Body)

	res.Body.Close()

	ipS := strings.TrimSpace(string(body))

	ip := net.ParseIP(ipS)
	return len(ip) != 0
}
