package smtp

import (
	"github.com/watsonserve/maild/auth"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

// helo命令
func HELO(ctx *SmtpContext) {
	ctx.Module = MOD_COMMAND
	addr := ctx.Address
	name := ctx.Msg[5:]

	ctx.Send(fmt.Sprintf("250 %s Hello %s (%s[%s])\r\n", ctx.Domain, name, addr, addr))
}

// ehlo命令
func EHLO(ctx *SmtpContext) {
	ctx.Module = MOD_COMMAND
	addr := ctx.Address
	name := ctx.Msg[5:]
	msg := fmt.Sprintf(
		"250-%s Hello %s (%s[%s])\r\n%s\r\n%s\r\n%s\r\n%s\r\n",
		ctx.Domain, name, addr, addr,
		"250-AUTH LOGIN PLAIN",
		"250-AUTH=LOGIN PLAIN",
		"250-PIPELINING",
		"250 ENHANCEDSTATUSCODES",
	)
	ctx.Send(msg)
}

// 授权
func AUTH(ctx *SmtpContext) {
	content, err := base64.StdEncoding.DecodeString(ctx.Msg[11:])
	if nil != err {
		fmt.Fprintln(os.Stderr, "error: "+err.Error())
		return
	}

	for i := 0; i < len(content); i++ {
		if 0 == content[i] {
			content[i] = '\n'
		}
	}
	author := auth.New()
	userPassword := strings.Split(string(content), "\n")
	userId := author.Auth(userPassword[0], userPassword[1])
	buf := "535 Authentication Failed\r\n"
	if "" != userId {
		ctx.User = userPassword[0]
		buf = "235 Authentication Successful\r\n"
		ctx.Login = true
		fmt.Println("auth by self")
	}
	ctx.Send(buf)
}

//
func QUIT(ctx *SmtpContext) {
	ctx.End("221 2.0.0 " + ctx.Domain + " Service closing transmission channel\r\n")
}

//
func XCLIENT(ctx *SmtpContext) {
	fmt.Println("auth by agency")
	ctx.Login = true
	ctx.Hola()
}

//
func STARTTLS(ctx *SmtpContext) {
	ctx.Send("502 5.3.3 STARTTLS is not supported\r\n")
	fmt.Println("startTTS")
}

//
func HELP(ctx *SmtpContext) {
	ctx.Send("502 5.3.3 HELP is not supported\r\n")
}

//
func NOOP(ctx *SmtpContext) {
	ctx.Send("250 2.0.0 OK\r\n")
	fmt.Println("noop")
}

//
func RSET(ctx *SmtpContext) {
	ctx.Send("250 2.0.0 OK\r\n")
	fmt.Println("rset")
}

//
func MAIL(ctx *SmtpContext) {
	ctx.Email.Sender = ctx.re.FindStringSubmatch(ctx.Msg)[1]
	clientDomain := strings.Split(ctx.Email.Sender, "@")[1]
	if (clientDomain == ctx.Domain) != (!ctx.Login) { // 本域已登录 or 外域未登录
		ctx.Send("250 2.1.0 Sender <" + ctx.Email.Sender + "> OK\r\n")
		return
	}
	ctx.Send("530 5.7.1 Authentication Required\r\n")
}

//
func RCPT(ctx *SmtpContext) {
	recver := ctx.re.FindStringSubmatch(ctx.Msg)[1]
	if strings.Split(recver, "@")[1] != ctx.Domain && !ctx.Login { // 非登录用户 to 外域
		ctx.Send("530 5.7.1 Authentication Required\r\n")
		return
	}
	//fmt.Println(strings.Split(recver, "@")[1] + " ", !ctx.Login)
	ctx.Email.Recver.PushBack(recver)
	ctx.Send("250 2.1.5 Recipient <" + recver + "> OK\r\n")
}

//
func DATA(ctx *SmtpContext) {
	format := "from %s ([%s]) by %s over TLS secured channel with %s(%s)\r\n\t%d"
	ctx.Module = MOD_HEAD
	ele := &KV {
		Name:  "Received",
		Value: fmt.Sprintf(format, ctx.Domain, ctx.Ip, ctx.Domain, ctx.Name, ctx.Version, time.Now().Unix()),
	}
	ctx.Email.Head = append(ctx.Email.Head, *ele)

	ctx.Send("354 Ok Send data ending with <CRLF>.<CRLF>\r\n")
}
