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

func usage() {
	fmt.Println("Usage: \n")
	fmt.Println("-host etc: 127.0.0.1\n")
	fmt.Println("-port etc: 8989\n")
	fmt.Println("-file file path, etc: /var/log.txt\n")
	fmt.Println("-dest remote server destination folder, etc: /tmp\n")
}

func main() {
	if len(os.Args) != 5 {
		usage()
		return
	}
	filePtr := flag.String("file", "", "file path; etc: /var/log.txt")
	hostPtr := flag.String("host", "", "host; etc: 127.0.0.1")
	portPtr := flag.String("port", "", "port; etc: 8989")
	destPtr := flag.String("dest", "", "etc: /tmp")
	flag.Parse()
	file := *filePtr
	host := *hostPtr
	port := *portPtr
	dest := *destPtr

	if file == "" {
		log.Fatal("invalid file path")
		os.Exit(1)
	}
	if !fileExist(file) {
		log.Fatal("file no exist...")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+port)
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
		data := file + "##%%^^##" + str + "##%%^^##" + fmd5 + "##%%^^##" + dest
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
