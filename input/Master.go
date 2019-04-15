package input

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Verbose    bool
	OutputFile string
	MaxTimeout int
	Work       int
	IP         string
	MinPort    int
	MaxPort    int
}

func Master() Data {
	var m string

	verbose := flag.Bool("v", false, "显示详细信息")
	output := flag.String("o", "output.txt", "输出文件文件名")
	maxTimeout := flag.Int("m", 10, "最大超时时间，单位为秒")
	ports := flag.String("p", "0-65535", "扫描端口范围")
	work := flag.Int("w", 500, "每秒请求数,最多不能超过 1000000")

	if *work > 1000000 {
		fmt.Println("请求数超过容量")
		os.Exit(0)
	}

	flag.Usage = usage

	flag.Parse()

	in := flag.Args()
	switch len(in) {
	case 0:
		flag.Usage()
		os.Exit(0)
	case 1:
		context := in[0]
		m = option(context)
	default:
		fmt.Println("参数过多")
		os.Exit(0)
	}

	foo := strings.Split(*ports, "-")
	var minPort, maxPort int
	var err error
	if len(foo) == 2 {
		a := foo[0]
		b := foo[1]
		minPort, err = strconv.Atoi(a)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		maxPort, err = strconv.Atoi(b)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		if minPort > maxPort {
			fmt.Println("起始位置，比结束位置大，程序退出")
			os.Exit(0)
		}
	} else {
		fmt.Println("端口参数错误")
		os.Exit(0)
	}

	temp := Data{
		Verbose:    *verbose,
		OutputFile: *output,
		MaxTimeout: *maxTimeout,
		Work:       *work,
		IP:         m,
		MinPort:    minPort,
		MaxPort:    maxPort,
	}
	return temp

}

func option(in string) string {
	var s string
	switch in {
	case "version":
		version()
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	default:
		a := net.ParseIP(in)
		if a != nil {
			in = in + "/32"
		}
		_, b, c := net.ParseCIDR(in)
		if c != nil {
			fmt.Println("参数错误")
			os.Exit(0)
		}
		s = b.String()
	}
	return s
}
