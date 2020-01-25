package smtpd

import (
    "fmt"
    "net"
    "github.com/watsonserve/maild"
    "regexp"
    "time"
)

const (
    MOD_COMMAND = 1
    MOD_HEAD = 2
    MOD_BODY = 4
)

type SmtpContext struct {
    maild.SentStream
    handlers maild.SmtpServerConfigure
    conf     *maild.ServerConfig
    Module   int
    Login    bool
	re       *regexp.Regexp
    Email    *maild.Mail
    // 其他
    Msg      string
    User     string
}

func InitSmtpContext(sock net.Conn, config maild.SmtpServerConfigure) *SmtpContext {
    this := &SmtpContext{
        handlers: config,
        conf: config.GetConfig(),
        Module: MOD_COMMAND,
        Login: false,
        re: regexp.MustCompile("<(.+)>"),
        Email: &maild.Mail{},
    }

    this.SentStream = *maild.InitSentStream(sock)

    return this
}

// 问候语
func (this *SmtpContext) Hola() {
    config := this.conf
	this.Send(fmt.Sprintf("220 %s %s Server (%s %s Server %s) ready %d\r\n",
        config.Domain, config.Type, config.Name, config.Type, config.Version, time.Now().Unix(),
	))
}

func (this *SmtpContext) TakeOff() {
    this.handlers.TakeOff(this.Email)
}
