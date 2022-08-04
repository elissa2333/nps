package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/elissa2333/nps/check"
	"github.com/elissa2333/nps/input"

	"golang.org/x/sync/semaphore"
)

func main() {
	data := input.Master()

	coroutine := data.MaxTimeoutSecond * data.Work
	ctx := context.Background()
	sem := semaphore.NewWeighted(int64(coroutine))

	tiout := 1000000 / data.Work

	ch := make(chan string, 1000)

	rec := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()

		f, err := os.Create(data.OutputFile)
		if err != nil {
			fmt.Println("写文件错误： ", err)
		}
		defer f.Close()

		sum := 0

	loop:
		for {
			select {
			case <-rec:
				break loop
			case x := <-ch:
				if sum == 0 {
					f.WriteString(x)
					sum++
				} else {
					f.WriteString("\n" + x)
				}
			}
		}
	}()

	a, b, _ := net.ParseCIDR(data.IP)

	minPort := data.MinPort
	maxPort := data.MaxPort

	for a := a.Mask(b.Mask); b.Contains(a); forIP(a) {
		for i := minPort; i < maxPort+1; i++ {
			time.Sleep(time.Duration(tiout) * time.Microsecond) //时间除以进程数

			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("start 无法获取信号量: %v\n", err)
				break
			}

			go work(a.String(), uint16(i), sem, data, ch)

		}
	}

	if err := sem.Acquire(ctx, int64(coroutine)); err != nil { //退出循环后阻塞主进程以等待最后的 goroutine 执行完成
		log.Printf("exit 无法获取信号量: %v\n", err)
	}

	rec <- struct{}{}

	<-done
}

func work(host string, port uint16, sem *semaphore.Weighted, data input.Data, stream chan<- string) {
	defer sem.Release(1)

	var wg sync.WaitGroup

	if data.Verbose {
		fmt.Printf("正在扫描%v:%d\n", host, port)
	}

	timeout := time.Duration(data.MaxTimeoutSecond) * time.Second
	isOpen := check.TCPPort(host, port, timeout)

	if isOpen {
		wg.Add(2)

		address := net.JoinHostPort(host, strconv.Itoa(int(port)))
		fmt.Printf("检测到 %v 端口开放，正在进行下一步检测\n", address)
		go func() {
			httpAddress := "http://" + address
			defer wg.Done()
			isHttp := check.IsProxy(httpAddress, timeout)
			if isHttp {
				stream <- httpAddress

				fmt.Printf("%v 为HTTP代理\n", address)
			}
		}()

		go func() {
			defer wg.Done()

			socks5Address := "socks5://" + address
			isSocks5 := check.IsProxy(socks5Address, timeout)
			if isSocks5 {
				stream <- socks5Address
				fmt.Printf("%v 为socks5代理\n", address)
			}
		}()

		wg.Wait()
	}

}

func forIP(ip net.IP) {

	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
