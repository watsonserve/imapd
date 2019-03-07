package smtp

import (
    "fmt"
    "net"
    "container/list"
    "github.com/watsonserve/maild/server"
    "regexp"
    "time"
)

type KV struct {
    Name string
    Value string
}

type Mail struct {
    Sender string
    Recver list.List
    Head []KV
    MailContent string
}

type SmtpdConfig struct {
	Domain  string
	Ip      string
	Name    string
	Type    string
	Version string
}

func (this *SmtpdConfig) CloneFrom(ctx *SmtpdConfig) {
	this.Domain = ctx.Domain
	this.Ip = ctx.Ip
	this.Name = ctx.Name
	this.Type = ctx.Type
	this.Version = ctx.Version
}

const (
    MOD_COMMAND = 1
    MOD_HEAD = 2
    MOD_BODY = 4
)

type SmtpContext struct {
    server.SentStream
    SmtpdConfig
    Module int
    Login  bool
	re     *regexp.Regexp
    Msg    string
    User   string
    Email  Mail
}

func InitSmtpContext(sock net.Conn) *SmtpContext {
    ret := &SmtpContext{}

    ret.SentStream = *server.InitSentStream(sock)
    ret.Module = MOD_COMMAND
    ret.Login = false
	ret.re = regexp.MustCompile("<(.+)>")
    return ret
}

// 问候语
func (this *SmtpContext) Hola() {
	this.Send(fmt.Sprintf("220 %s %s Server (%s %s Server %s) ready %d\r\n",
		this.Domain, this.Type, this.Name, this.Type, this.Version, time.Now().Unix(),
	))
}

func (this *SmtpContext) TakeOff() {
    fmt.Println(this.Email.Head)
    fmt.Println(this.Email.MailContent)
}
