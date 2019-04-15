package input

import (
	"flag"
	"fmt"
	"os"
)

func usage() { //help info
	fmt.Fprintf(os.Stderr, `
Author: niconiconi

由 go 编写的一款高性能的网络代理扫描器，可同时扫描 socks5 代理和 http 代理，每秒可发起到达一百万次的请求。

本软件可以对 单 IP 或 IP 段进行扫描

例如：

	nps 8.8.8.8

或

	nps 8.8.8.0/24

命令行参数：
   help      查看帮助信息
   version   查看版本信息

`)
	flag.PrintDefaults()
}
