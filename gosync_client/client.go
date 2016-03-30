package main

import (
	//	"encoding/json"
	"flag"
	"fmt"
	"github.com/gansidui/gotcp/examples/echo"
	//"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

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
		str, err := getFileContent(file)
		if err != nil {
			log.Fatal(err)
		}

		fmd5, err := md5String(str)
		if err != nil {
			log.Fatal(err)
		}
		data := file + "##%%^^##" + str + "##%%^^##" + fmd5
		conn.Write(echo.NewEchoPacket([]byte(data), false).Serialize())
		p, err := echoProtocol.ReadPacket(conn)
		if err == nil {
			echoPacket := p.(*echo.EchoPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
		}

		time.Sleep(1 * time.Second)
	}

	conn.Close()
}
