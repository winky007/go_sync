package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gansidui/gotcp/examples/echo"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type jsonParam struct {
	File    string `json:"file"`
	Md5     string `json:"md5"`
	Content string `json:"content"`
}

type jsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int64  `json:"code"`
}

func jsonReturn(j *jsonResponse) string {
	res2B, _ := json.Marshal(j)
	return string(res2B)
}

func main() {
	filePtr := flag.String("file", "", "file path; etc: /var/log.txt")
	flag.Parse()
	file := *filePtr

	if file == "" {
		log.Fatal("invalid file path")
		os.Exit(1)
	}
	if !fileExist(file) {
		log.Fatal("file no exist...")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	echoProtocol := &echo.EchoProtocol{}

	for {
		fmd5, err := md5File(file)
		if err != nil {
			log.Fatal(err)
		}

		cByte, err := ioutil.ReadFile(file)
		content := string(cByte)
		jsResp := jsonParam{File: file, Md5: fmd5, Content: content}
		res2B, _ := json.Marshal(jsResp)
		js := string(res2B)
		conn.Write(echo.NewEchoPacket([]byte(js), false).Serialize())
		p, err := echoProtocol.ReadPacket(conn)
		if err == nil {
			echoPacket := p.(*echo.EchoPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
		}

		time.Sleep(1 * time.Second)
	}

	conn.Close()
}
