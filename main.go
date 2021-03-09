package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"sync"
	"time"
)

var (
	help bool
	inputfile string
	outputfile string
	names []string
	ips []string
)

func init(){
	flag.BoolVar(&help,"h, --help",false,"help, 帮助命令")
	flag.StringVar(&inputfile,"i","","要读取的文件")
	flag.StringVar(&outputfile,"o","","要输出的文件")
	flag.Usage = usage
	flag.Parse()

}

func usage(){
	fmt.Fprintf(color.Output,color.HiCyanString(`Ameng编写的Go语言域名解析工具
Options:
`))
	flag.PrintDefaults()
}


func Getip(inputfile string, outputfile string)(count1 int, count2 int){

	infile, err := os.OpenFile(inputfile,os.O_RDONLY,1)
	outfile, err := os.OpenFile(outputfile,os.O_WRONLY,2)
	if err != nil {
		fmt.Fprintln(color.Output,time.Now().Format("2006/01/02 12:22:48"),color.RedString("[ERROR]")+":",err)
	}
	hostnames := bufio.NewScanner(infile)
	var wg sync.WaitGroup
	for hostnames.Scan() {
		count1++
		hostname := hostnames.Text()
		names = append(names, hostname)
		go func() {
			wg.Add(1)
			defer wg.Done()
			ip, err := net.ResolveIPAddr("ip",hostname)
			if err != nil {
				fmt.Fprintln(color.Output,time.Now().Format("2006/01/02 12:22:48"),color.YellowString("[WARNING]")+":",err)
			}
			if ip != nil {
				ips = append(ips, ip.String())
			}
		}()
	}
	wg.Wait()


	result := Removesamesip(ips)
	for _,v := range result{
		outfile.WriteString(v+"\n")
	}
	return count1,len(result)
}

func Removesamesip(ips [] string)(result []string){
	result = make([]string, 0)
	tempMap := make(map[string]bool, len(ips))
	for _, e := range ips{
		if tempMap[e] == false{
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}


func main(){
	if inputfile != "" && outputfile != "" {
		start := time.Now()
		count1,count2 := Getip(inputfile,outputfile)
		end := time.Since(start)
		defer fmt.Fprint(color.Output,color.GreenString("域名解析已完成,总用时:"),end,color.GreenString(" 解析域名数:"),count1,color.GreenString(" 获得ip数:"),count2)
	}

}
