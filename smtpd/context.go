package smtpd

import (
    "fmt"
    "net"
    "github.com/watsonserve/maild/lib"
    "regexp"
    "time"
)

const (
    MOD_COMMAND = 1
    MOD_HEAD = 2
    MOD_BODY = 4
)

type SmtpContext struct {
    lib.SentStream
    lib.ServerConfig
    Module int
    Login  bool
	re     *regexp.Regexp
    Msg    string
    User   string
    Email  *lib.Mail
}

func InitSmtpContext(sock net.Conn) *SmtpContext {
    ret := &SmtpContext{}

    ret.SentStream = *lib.InitSentStream(sock)
    ret.Module = MOD_COMMAND
    ret.Login = false
    ret.re = regexp.MustCompile("<(.+)>")
    ret.Email = &lib.Mail{}
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
