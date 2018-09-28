package smtpContext

import (
    "net"
    "container/list"
    "tcpServer"
)

type KV struct {
    Name string
    Value string
}

type SmtpContext struct {
    tcpServer.SentStream
    Module int
    Address string
    Msg string
    MailContent string
    Login bool
    User string
    Sender string
    Recver list.List
    Head []KV
    Sock net.Conn
}

const (
    MOD_COMMAND = 1
    MOD_HEAD = 2
    MOD_BODY = 4
)

func InitSmtpContext(sock net.Conn) *SmtpContext {
    ret := &SmtpContext{}

    ret.MailContent = ""
    ret.Module = MOD_COMMAND
    ret.Login = false
    ret.Address = sock.RemoteAddr().String()
    ret.Sock = sock
    return ret
}
