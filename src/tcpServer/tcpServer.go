package tcpServer

import (
    "fmt"
    "net"
    "log"
    "dispatcher"
)

type TcpServer struct {
    disp dispatcher.Dispatcher
}

type SentStream struct {
    Address string
    Sock net.Conn
}

func Init() *TcpServer {
    return &TcpServer{}
}

func (this *TcpServer) SetDispatcher(disp dispatcher.Dispatcher) {
    this.disp = disp
}

func (this *TcpServer) Listen(port string) int {
    // port = ":465"
    ln, err := net.Listen("tcp", port)
    if err != nil {
        return -1
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println("a connect exception")
        }
        go this.disp.Task(conn)
    }
}

func InitSentStream(sock net.Conn) *SentStream {
    ret := &SentStream {}
    ret.Address = sock.RemoteAddr().String()
    ret.Sock = sock
    return ret
}

// 发送
func (this *SentStream) Send(content string) {
    //fmt.Print(content)
    fmt.Fprint(this.Sock, content)
}

// 发送并关闭
func (this *SentStream) End(content string) {
    fmt.Fprint(this.Sock, content)
    this.Sock.Close()
}
