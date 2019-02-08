package smtpContext

import (
    "fmt"
    "net"
    "container/list"
    "server"
)

type KV struct {
    Name string
    Value string
}

type SmtpContext struct {
    server.SentStream
    Module int
    MailContent string
    Login bool

    Msg string
    User string
    Sender string
    Recver list.List
    Head []KV
}

const (
    MOD_COMMAND = 1
    MOD_HEAD = 2
    MOD_BODY = 4
)

func InitSmtpContext(sock net.Conn) *SmtpContext {
    ret := &SmtpContext{}

    ret.SentStream = *server.InitSentStream(sock)
    ret.Module = MOD_COMMAND
    ret.MailContent = ""
    ret.Login = false
    return ret
}

func (this *SmtpContext) TakeOff() {
    fmt.Println(this.Head)
    fmt.Println(this.MailContent)
}