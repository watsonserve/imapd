package smtp

import (
	"fmt"
    "github.com/watsonserve/maild"
	"github.com/watsonserve/maild/server"
	"net"
	"bufio"
	"os"
	"errors"
	"strings"
)

type Smtpd struct {
	server.TcpServer
	maild.ServerConfig
	dict    map[string]func(*SmtpContext)
}

func New(domain string, ip string) *Smtpd {
	ret := &Smtpd{}

	ret.Domain = domain
	ret.Ip = ip
	ret.Type = "SMTP"
	ret.Name = "WS_SMTPD"
	ret.Version = "1.0"

	ret.dict = map[string]func(*SmtpContext) {
		"HELO": HELO,
		"EHLO": EHLO,
		"AUTH": AUTH,
		"QUIT": QUIT,
		"XCLIENT": XCLIENT,
		"STARTTLS": STARTTLS,
		"HELP": HELP,
		"NOOP": NOOP,
		"RSET": RSET,
		"MAIL": MAIL,
		"RCPT": RCPT,
		"DATA": DATA,
	}
	return ret
}

func (this *Smtpd) Task(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	ctx := InitSmtpContext(conn)
	ctx.CloneFrom(&this.ServerConfig)
	ctx.Hola()

	for scanner.Scan() {
		err := scanner.Err()
		if nil != err {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			break
		}
		// 行遍历
		msg := scanner.Text()
		ctx.Msg = msg
		// fmt.Println(ctx.Msg)
		switch ctx.Module {
		case MOD_COMMAND:
			err = commandHash(this, ctx)
			if nil != err {
				fmt.Fprintln(os.Stderr, err)
			}
		case MOD_HEAD:
			dataHead(ctx)
		case MOD_BODY:
			dataBody(ctx)
		}
	}
}

func (this *Smtpd) Listen(port string) {
	this.TcpServer.Listen(port, this)
}


func commandHash(this *Smtpd, ctx *SmtpContext) error {
	var key string
	// 截取第一个单词
	_, err := fmt.Sscanf(ctx.Msg, "%s", &key)
	if nil != err {
		return err
	}
	// 查找处理方法
	method, exist := this.dict[key]
	if !exist {
		ctx.Send("method " + key + " not valid\r\n")
		return errors.New("method " + key + " not valid")
	}
	// 执行处理
	method(ctx)
	return nil
}

func dataHead(ctx *SmtpContext) {
	if "" == ctx.Msg {
		ctx.Module = MOD_BODY
	} else if ' ' == ctx.Msg[0] || '\t' == ctx.Msg[0] {
		ctx.Email.Head[len(ctx.Email.Head)-1].Value += "\r\n" + ctx.Msg
	} else {
		attr := strings.Split(ctx.Msg, ": ")
		ele := &maild.KV {
			Name:  attr[0],
			Value: attr[1],
		}
		ctx.Email.Head = append(ctx.Email.Head, *ele)
	}
}

func dataBody(ctx *SmtpContext) {
	if "." == ctx.Msg {
		ctx.Module = MOD_COMMAND
		ctx.Send("250 2.6.0 Message received\r\n")
		ctx.TakeOff()
		fmt.Println("250 2.6.0 Message received")
		return
	}
	ctx.Email.MailContent += ctx.Msg + "\r\n"
}
