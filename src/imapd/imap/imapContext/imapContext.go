package imapContext

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

type ImapContext struct {
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

func InitImapContext(sock net.Conn) *ImapContext {
    ret := &ImapContext{}

    ret.SentStream = *server.InitSentStream(sock)
    ret.Module = MOD_COMMAND
    ret.MailContent = ""
    ret.Login = false
    return ret
}

func (this *ImapContext) TakeOff() {
    fmt.Println(this.Head)
    fmt.Println(this.MailContent)
}
