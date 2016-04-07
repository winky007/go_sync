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
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("-host etc: 127.0.0.1")
	fmt.Println("-port etc: 8989")
	fmt.Println("-local local folder path, etc: /var/log or /var/log/log.txt")
	fmt.Println("-dest remote server destination folder, etc: /tmp")
	fmt.Println("-interval, optional, sync's internal time, unit is second , default is 3 seconds, etc: 3")
	fmt.Println("-usereg optional, 1 or 0")
	fmt.Println("-regexp optional, use regexp to filter content, etc: \"start.*endstart\"")
}

func main() {
	if len(os.Args) < 5 {
		usage()
		return
	}
	localPtr := flag.String("local", "", "local folder path or file path; etc: /var/log/ or /var/log.txt")
	hostPtr := flag.String("host", "", "host; etc: 127.0.0.1")
	portPtr := flag.String("port", "", "port; etc: 8989")
	destPtr := flag.String("dest", "", "etc: /tmp")
	intervalPtr := flag.Int("interval", 3, "sync's internal time, unit is second , default is 3 seconds; etc: 3")
	useregPtr := flag.Int("usereg", 0, "optional, 1 or 0")
	regexpPtr := flag.String("regexp", "", "optional, use regexp to filter content, etc: \"start.*endstart\"")

	flag.Parse()
	local := *localPtr
	host := *hostPtr
	port := *portPtr
	dest := *destPtr
	interval := *intervalPtr
	usereg := *useregPtr
	regexpStr := *regexpPtr

	if local == "" {
		log.Fatal("invalid local folder path")
		os.Exit(1)
	}
	if !fileExist(local) {
		log.Fatal("local folder no exist...")
		os.Exit(1)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host+":"+port)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	echoProtocol := &echo.EchoProtocol{}

	for {
		fileList := []string{}
		err := filepath.Walk(local, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				fileList = append(fileList, path)
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
		local = strings.TrimRight(local, `/`)
		dir := strings.Split(local, `/`)
		parentDirS := dir[0 : len(dir)-1]
		parentDir := strings.Join(parentDirS, `/`) + `/`
		for _, file := range fileList {
			str, err := getFileContent(file)
			if err != nil {
				log.Fatal(err)
			}
			if usereg == 1 && regexpStr != "" {
				re := regexp.MustCompile(regexpStr)
				strSlice := re.FindAllString(str, -1)
				str = strings.Join(strSlice, "")
			}

			fmd5, err := md5String(str)
			if err != nil {
				log.Fatal(err)
			}
			destNew := strings.TrimRight(dest, "/")
			destNew = destNew + "/"
			destNew = destNew + strings.Replace(file, parentDir, "", 1)
			destNewS := strings.Split(destNew, "/")
			destNewS = destNewS[0 : len(destNewS)-1]
			destNew = strings.Join(destNewS, "/") + "/"
			data := file + "##%%^^##" + str + "##%%^^##" + fmd5 + "##%%^^##" + destNew
			conn.Write(echo.NewEchoPacket([]byte(data), false).Serialize())
			p, err := echoProtocol.ReadPacket(conn)
			if err == nil {
				echoPacket := p.(*echo.EchoPacket)
				fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}

	conn.Close()
}
