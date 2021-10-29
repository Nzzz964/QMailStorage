package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"qmailstorage/client"
	"qmailstorage/utils"
	"strings"
)

var c string
var f string

func init() {
	flag.StringVar(&c, "c", "config.json", "指定配置文件")
	flag.StringVar(&f, "f", "", "指定上传文件")
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
	data, err := ioutil.ReadFile(c)
	if err != nil {
		log.Panic(err)
	}

	config := &Config{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Panic(err)
	}

	smtpClient, err := client.NewClient(config.Server, config.Username, config.Password)
	if err != nil {
		log.Panic(err)
	}
	defer smtpClient.Close()

	fileInfo, err := os.Stat(f)
	if os.IsNotExist(err) {
		log.Panic(err)
	}

	basename := fileInfo.Name()
	sha1str, err := utils.Sha1File(f)
	if err != nil {
		log.Println("计算文件 sha1 失败")
	}

	chunks, err := utils.MakeChunk(f, config.Chunksize)
	if err != nil {
		if len(chunks) > 0 {
			utils.RemoveChunk(chunks)
		}
		log.Panic(err)
	}

	defer utils.RemoveChunk(chunks)

	mime := client.NewMail(config.From, config.To)
	log.Println("总分片数: " + fmt.Sprint(len(chunks)))

	for k, v := range chunks {
		chunkBaseName := path.Base(v)
		mime.Subject("文件: " + basename + " 分片: " + fmt.Sprintf("%d", k))
		mime.Text(strings.Join([]string{
			"sha1: " + sha1str,
			"file: " + path.Base(v),
		}, " \n"))
		mime.Attach(v, chunkBaseName)

		if err = smtpClient.Mail(config.From); err != nil {
			log.Panic(err)
		}
		if err = smtpClient.Rcpt(config.To); err != nil {
			log.Panic(err)
		}

		writer, err := smtpClient.Data()
		if err != nil {
			log.Panic(err)
		}

		log.Println("发送分片:", k)
		mime.WriteTo(writer)

		mime.Reset()
		writer.Close()
	}

}
