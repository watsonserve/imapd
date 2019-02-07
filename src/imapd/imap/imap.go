package imap

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"imapd/imap/imapContext"
	"net"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

/*
const (
    BUFSIZ = 8192
)
*/

type Database struct {
}

func (this *Database) Auth(username string, password string) string {
	return ""
}

type Imapd struct {
	Domain  string
	Type    string
	Name    string
	Version string
	Ip      string
	re      *regexp.Regexp
}

func Init(domain string, ip string) *Imapd {
	ret := &Imapd{}
	ret.Domain = domain
	ret.Type = "SMTP"
	ret.Name = "WS_IMAPD"
	ret.Version = "1.0"
	ret.Ip = ip
	ret.re = regexp.MustCompile("<(.+)>")
	return ret
}

// 问候语
func (this *Imapd) Hola() string {
	return fmt.Sprintf("* OK %s IMAP4Server Ready\r\n", this.Name)
}

// helo命令
func (this *Imapd) HELO(client *imapContext.ImapContext) {
	client.Module = imapContext.MOD_COMMAND
	addr := client.Address
	name := client.Msg[5:]
	client.Send("250 " + this.Domain + " Hello " + name + " (" + addr + "[" + addr + "])\r\n")
}

// ehlo命令
func (this *Imapd) EHLO(client *imapContext.ImapContext) {
	client.Module = imapContext.MOD_COMMAND
	addr := client.Address
	name := client.Msg[5:]
	client.Send("250-" + this.Domain + " Hello " + name + " (" + addr + "[" + addr + "])\r\n250-AUTH LOGIN PLAIN\r\n250-AUTH=LOGIN PLAIN\r\n250-PIPELINING\r\n250 ENHANCEDSTATUSCODES\r\n")
}

//
func (this *Imapd) QUIT(client *imapContext.ImapContext) {
	client.End("221 2.0.0 " + this.Domain + " Service closing transmission channel\r\n")
	fmt.Println(client.Head)
	fmt.Println(client.MailContent)
}

//
func (this *Imapd) XCLIENT(client *imapContext.ImapContext) {
	fmt.Println("auth by agency")
	client.Login = true
	client.Send(this.Hola())
}

//
func (this *Imapd) STARTTLS(client *imapContext.ImapContext) {
	client.Send("502 5.3.3 STARTTLS is not supported\r\n")
}

//
func (this *Imapd) HELP(client *imapContext.ImapContext) {
	client.Send("502 5.3.3 HELP is not supported\r\n")
}

//
func (this *Imapd) NOOP(client *imapContext.ImapContext) {
	client.Send("250 2.0.0 OK\r\n")
}

//
func (this *Imapd) RSET(client *imapContext.ImapContext) {
	client.Send("250 2.0.0 OK\r\n")
}

//
func (this *Imapd) DATA(client *imapContext.ImapContext) bool {
	format := "from %s ([%s]) by %s over TLS secured channel with %s(%s)\r\n\t%d"
	client.Module = imapContext.MOD_HEAD
	ele := &imapContext.KV{
		Name:  "Received",
		Value: fmt.Sprintf(format, this.Domain, this.Ip, this.Domain, this.Name, this.Version, time.Now().Unix()),
	}
	client.Head = append(client.Head, *ele)

	client.Send("354 Ok Send data ending with <CRLF>.<CRLF>\r\n")
	return false
}

func (this *Imapd) CommandHash(client *imapContext.ImapContext) error {
	var key string
	// 截取第一个单词
	_, err := fmt.Sscanf(client.Msg, "%d %s", &key)
	if nil != err {
		return err
	}
	// 查找处理方法
	that := reflect.ValueOf(this)
	method := that.MethodByName(key)
	if !method.IsValid() {
		return errors.New("method " + key + "not valid")
	}
	// 执行处理
	clientValue := reflect.ValueOf(client)
	inArgs := []reflect.Value{clientValue}
	method.Call(inArgs)[0]
	return nil
}

func (this *Imapd) Task(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	client := imapContext.InitImapContext(conn)
	client.Send(this.Hola())

	for scanner.Scan() {
		err := scanner.Err()
		if nil != err {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			break
		}
		// 行遍历
		client.Msg = scanner.Text()
		//if "QUIT" == client.Msg {
		//}
		fmt.Println(client.Msg)
		Lexical(client.Msg)
		switch client.Module {
		case imapContext.MOD_COMMAND:
			err = this.CommandHash(client)
			if nil != err {
				fmt.Fprintln(os.Stderr, err)
			}
		case imapContext.MOD_HEAD:
			this.DataHead(client)
		case imapContext.MOD_BODY:
			this.DataBody(client)
		}
	}
}
