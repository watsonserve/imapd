// IMAP agent service
package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/watsonserve/imapd/lib"
)

type UpResult struct {
	Sess   string
	Result string
}

type ImapAgentFace interface {
	Auth(username string, password string) string
	Read() *UpResult
	Send(sess string, spt *lib.Mas)
}

type imap_agent_t struct {
	iface         ImapAgentFace
	name          string
	re            *regexp.Regexp
	outPermission map[string]bool // 不需要登录态
	Uppop         map[string]bool // 需要上游服务器处理
	sessMap       map[string]*imap_agent_context_t
}

func (this *imap_agent_t) upHandle() {
	for {
		// 从服务端读取数据
		data := this.iface.Read()
		if nil == data {
			continue
		}
		ctx, ex := this.sessMap[data.Sess]
		if !ex {
			continue
		}
		ctx.Send(data.Result)
	}
}

func (this *imap_agent_t) task(conn net.Conn) {
	ctx := initImapAgentContext(conn)
	fmt.Println("hello client")
	ctx.Send("* OK " + this.name + " IMAP4 service is ready.\r\n")

	for {
		msg, err := ctx.ReadLine()
		if nil == err {
			if "" == msg {
				err = errors.New("EOF")
			}
		}
		if nil != err {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			break
		}
		script := lib.InitMas(msg)
		if nil == script {
			ctx.Send(fmt.Sprintf("%s BAD Command Error.\r\n", msg))
			continue
		}

		err = this.commandHash(ctx, script)
		if nil != err {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (this *imap_agent_t) commandHash(ctx *imap_agent_context_t, script *lib.Mas) error {
	// 鉴权
	signIn := ctx.Checked()
	need := this.needPermission(script.Command)
	if !signIn && need {
		ctx.Send(fmt.Sprintf("%s BAD Command received in Invalid state.", script.Tag))
		return nil
	}

	// 查找处理方法并执行处理
	switch script.Command {
	case "CAPABILITY":
		ctx.CAPABILITY(script.Tag)
	case "AUTHENTICATE":
		parames := ctx.AUTHENTICATE(script)
		this.signIn(ctx, script.Tag, parames)
	case "LOGIN":
		parames := strings.Split(script.Parames, " ")
		this.signIn(ctx, script.Tag, parames)
	case "LOGOUT":
		ctx.LOGOUT(script.Tag)
	case "NOOP":
		ctx.NOOP(script.Tag)
	default:
		_, exist := this.Uppop[script.Command]
		if exist {
			this.iface.Send(ctx.Sess, script)
			break
		}
		ctx.Send(fmt.Sprintf("%s BAD %s is not supported.\r\n", script.Tag, script.Command))
		return errors.New("method " + script.Command + " not valid")
	}
	return nil
}

func (this *imap_agent_t) needPermission(command string) bool {
	_, exist := this.outPermission[command]
	return !exist
}

func (this *imap_agent_t) signIn(ctx *imap_agent_context_t, tag string, parames []string) {
	for 2 == len(parames) {
		sessionId := this.iface.Auth(parames[0], parames[1])
		if "" == sessionId {
			break
		}
		ctx.Sess = sessionId
		this.sessMap[sessionId] = ctx
		ctx.Send(fmt.Sprintf("%s OK LOGIN completed.\r\n", tag))
		return
	}
	ctx.Send(fmt.Sprintf("%s NO LOGIN FAILURE: username or password rejected.\r\n", tag))
}

func Service(name string, ln net.Listener, iface ImapAgentFace) {
	self := &imap_agent_t{
		iface: iface,
		name:  name,
		re:    regexp.MustCompile("<(.+)>"),
		outPermission: map[string]bool{
			"AUTHENTICATE": true,
			"LOGIN":        true,
			"LOGOUT":       true,
			"CAPABILITY":   true,
			"HELP":         true,
			"NOOP":         true,
			"QUIT":         true,
			"RSET":         true,
			"STARTTLS":     true,
			"XCLIENT":      true,
		},
		Uppop: map[string]bool{
			"SELECT":      true,
			"EXAMINE":     true,
			"CLOSE":       true,
			"CREATE":      true,
			"DELETE":      true,
			"RENAME":      true,
			"SUBSCRIBE":   true,
			"UNSUBSCRIBE": true,
			"LIST":        true,
			"LSUB":        true,
			"STATUS":      true,
			"APPEND":      true,
			"ID":          true,
			"UID":         true,
		},
	}

	go self.upHandle()
	for {
		conn, err := ln.Accept()
		if nil != err {
			log.Println("a connect exception", err)
			continue
		}
		defer conn.Close()
		go self.task(conn)
	}
}
