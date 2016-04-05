package main

import (
	"fmt"
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/examples/echo"
	//"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	echoPacket := p.(*echo.EchoPacket)
	str := string(echoPacket.GetBody())
	//fmt.Printf("OnMessage 1:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	s := strings.Split(str, `##%%^^##`)
	if l := len(s); l < 3 {
		fmt.Printf("OnMessage invalid data:[%v] [%v]\n", echoPacket.GetLength(), str)
		return true
	}

	file := s[0]
	content := s[1]
	md5 := s[2]
	Dest := s[3]
	Dest = strings.TrimRight(Dest, `/`)

	if _, err := os.Stat(Dest); os.IsNotExist(err) {
		err := os.MkdirAll(Dest, 0755)
		if err != nil {
			c.AsyncWritePacket(echo.NewEchoPacket([]byte("os.Stat error"), false), time.Second)
			return true
		}
	}

	fdir, err := os.Open(Dest)
	defer fdir.Close()
	if err != nil {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("os.Open error"), false), time.Second)
		return true
	}

	finfo, err := fdir.Stat()
	if err != nil {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("os.Stat error"), false), time.Second)
		return true
	}
	if finfo.Mode().IsRegular() {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("same file name already exit as dest param. please use other name"), false), time.Second)
		return true
	}
	f := strings.Split(file, `/`)
	flen := len(f)
	file = Dest + `/` + f[flen-1]
	if !fileExist(file) {
		status := writeFile(file, content)
		if status == true {
			c.AsyncWritePacket(echo.NewEchoPacket([]byte("success"), false), time.Second)
			return true
		} else {
			c.AsyncWritePacket(echo.NewEchoPacket([]byte("fail"), false), time.Second)
			return true
		}
	}

	fileContent, err := getFileContent(file)
	if err != nil {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("get file content error"), false), time.Second)
		return true
	}

	fmd5, err := md5String(fileContent)
	if err != nil {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("md5string error"), false), time.Second)
		return true
	}

	if fmd5 == md5 {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("success"), false), time.Second)
		return true
	}

	status := writeFile(file, content)
	if status == true {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("success"), false), time.Second)
		return true
	} else {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("write file error"), false), time.Second)
		return true
	}
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &Callback{}, &echo.EchoProtocol{})

	// starts service
	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}
