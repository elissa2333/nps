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

	master()
}

func master() {
	data := input.Master()
	coroutine := data.MaxTimeout * data.Work
	ctx := context.TODO()
	sem := semaphore.NewWeighted(int64(coroutine))

	tiout := 1000000 / data.Work

	ch := make(chan string, 1000)
	rec := make(chan bool)
	done := make(chan bool)
	go func() {
		f, err := os.Create(data.OutputFile)
		if err != nil {
			fmt.Println("写文件错误： ", err)
		}
		defer f.Close()

		sum := 0
		for {
			select {
			case <-rec:
				goto dones
			case x := <-ch:
				if sum == 0 {
					f.WriteString(x)
					sum++
				} else {
					f.WriteString("\n" + x)
				}
			}
		}

	dones:
		done <- true

	}()

	a, b, _ := net.ParseCIDR(data.IP)

	minPort := data.MinPort
	maxPort := data.MaxPort

	for ip := a.Mask(b.Mask); b.Contains(a); forIP(a) {
		for i := minPort; i < maxPort+1; i++ {
			f := strconv.Itoa(i)

			time.Sleep(time.Duration(tiout) * time.Microsecond) //时间除以进程数

			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("start 无法获取信号量: %v\n", err)
				break
			}

			go ck(ip.String(), f, sem, data, ch)

		}
	}

	if err := sem.Acquire(ctx, int64(coroutine)); err != nil { //退出循环后阻塞主进程以等待最后的 goroutine 执行完成
		log.Printf("exit 无法获取信号量: %v\n", err)
	}

	time.Sleep(1 * time.Second)
	rec <- true

	<-done

}

func ck(in, in2 string, sem *semaphore.Weighted, data input.Data, stream chan string) {
	defer sem.Release(1)

	var wg sync.WaitGroup

	if data.Verbose == true {
		fmt.Printf("正在扫描%v:%v\n", in, in2)
	}

	foo := check.TCPPort(in, in2, data.MaxTimeout)
	if foo == true {
		wg.Add(2)

		ber := in + ":" + in2
		fmt.Printf("检测到 %v 端口开放，正在进行下一步检测\n", ber)
		go func(kk string) {
			defer wg.Done()
			http := check.ProxyOfhttp(kk, data.MaxTimeout)
			if http == true {
				kk2 := "http://" + kk
				stream <- kk2

				fmt.Printf("%v 为HTTP代理\n", kk)
			}
		}(ber)

		go func(kk string) {
			defer wg.Done()
			socks := check.ProxyOfsocks5(kk, data.MaxTimeout)
			if socks == true {
				socks5 := "socks5://" + kk
				stream <- socks5
				fmt.Printf("%v 为socks5代理\n", kk)
			}
		}(ber)

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
