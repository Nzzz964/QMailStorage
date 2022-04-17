package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"qmailstorage/client"
)

var c string
var f string
var d string

func init() {
	flag.StringVar(&c, "c", "config.json", "指定配置文件")
	flag.StringVar(&f, "f", "", "指定上传文件")
	flag.StringVar(&d, "d", "", "文件描述")
	flag.Parse()
	checkArgs()
}

func checkArgs() {
	if f == "" {
		fmt.Fprintln(os.Stderr, "请指定上传文件 : qmailstorage.exe -f example.txt")
		os.Exit(1)
	}
}

func main() {
	// 读取配置文件
	data, err := ioutil.ReadFile(c)
	if err != nil {
		log.Panic(err)
	}

	config := &Config{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Panic(err)
	}

	// 初始化 smtp 客户端
	smtpClient, err := client.NewClient(config.Server, config.Username, config.Password)
	if err != nil {
		log.Panic(err)
	}
	defer smtpClient.Close()

	sender, err := client.NewQMailSender(smtpClient, f, config.Chunksize)
	if err != nil {
		log.Panic(err)
	}
	defer sender.Destory()

	sender.Description(d)
	err = sender.Send(config.From, config.To)
	if err != nil {
		log.Panic(err)
	}
}
