package main

import (
	"encoding/json"
	"fmt"
	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/examples/echo"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	echoPacket := p.(*echo.EchoPacket)
	info := string(echoPacket.GetBody())
	jsByte := []byte(info)
	var param jsonParam
	if err := json.Unmarshal(jsByte, &param); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OnMessage 1:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	if !fileExist(param.File) {
		writeFile(param.File, param.Content)
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("success"), false), time.Second)
	}

	fmd5, err := md5File(param.File)
	if err != nil {
		log.Fatal(err)
	}
	if fmd5 == param.Md5 {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("success"), false), time.Second)
		return true
	}

	status := writeFile(param.File, param.Content)
	if status == true {
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("success"), false), time.Second)
		return true
	} else {
		fmt.Println(status)
		c.AsyncWritePacket(echo.NewEchoPacket([]byte("fail"), false), time.Second)
	}
	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
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
