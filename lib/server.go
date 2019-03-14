package lib

import (
    "fmt"
    "net"
    "log"
)

type Dispatcher interface {
    Task(conn net.Conn)
}

type TcpServer struct {
    Dispatcher
}

func InitTCPServer() *TcpServer {
    return &TcpServer{}
}

/*
 * 这里使用的是每个链接启动一个新的go程的模型
 * 高并发的话，性能取决于go语言的协程能力
 */
func (this *TcpServer) Listen(port string, that Dispatcher) int {
    // port = ":465"
    ln, err := net.Listen("tcp", port)
    if err != nil {
        return -1
    }
    defer ln.Close()
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println("a connect exception")
        }
        defer conn.Close()
        go that.Task(conn)
    }
}


type SentStream struct {
    Address string
    Sock net.Conn
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