package maild

import (
    "bufio"
    "crypto/tls"
    "fmt"
    "io"
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

func (this *TcpServer) listen(ln net.Listener, that Dispatcher) error {
    defer ln.Close()
    for {
        conn, err := ln.Accept()
        if nil != err {
            log.Println("a connect exception")
        }
        defer conn.Close()
        go that.Task(conn)
    }
    return nil
}

/*
 * 这里使用的是每个链接启动一个新的go程的模型
 * 高并发的话，性能取决于go语言的协程能力
 */
func (this *TcpServer) TLSListen(port string, crt string, key string, that Dispatcher) error {
    cert, err := tls.LoadX509KeyPair(crt, key)
    if nil != err {
        return err
    }
    ln, err := tls.Listen("tcp", port, &tls.Config {
        Certificates: []tls.Certificate{cert},
        CipherSuites: []uint16 {
          tls.TLS_RSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
          tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        },
        PreferServerCipherSuites: true,
    })
    if nil != err {
        return err
    }
    return this.listen(ln, that)
}

/*
 * 这里使用的是每个链接启动一个新的go程的模型
 * 高并发的话，性能取决于go语言的协程能力
 */
func (this *TcpServer) Listen(port string, that Dispatcher) error {
    ln, err := net.Listen("tcp", port)
    if nil != err {
        return err
    }
    return this.listen(ln, that)
}


type ReadStream struct {
    scanner *bufio.Scanner
}

func InitReadStream(sock io.Reader) *ReadStream {
    return &ReadStream {
        scanner: bufio.NewScanner(sock),
    }
}

func (this *ReadStream) ReadLine() (string, error) {
    this.scanner.Scan()
    err := this.scanner.Err()
    if nil != err {
        return "", err
    }
    msg := this.scanner.Text()
    fmt.Printf("c: %s\n", msg)
    return msg, nil
}


type SentStream struct {
    Sock io.WriteCloser
}

func InitSentStream(sock io.WriteCloser) *SentStream {
    return &SentStream {
        Sock: sock,
    }
}

// 发送
func (this *SentStream) Send(content string) {
    fmt.Printf("s: %s\n", content)
    fmt.Fprint(this.Sock, content)
}

// 发送并关闭
func (this *SentStream) End(content string) {
    fmt.Fprint(this.Sock, content)
    this.Sock.Close()
}
