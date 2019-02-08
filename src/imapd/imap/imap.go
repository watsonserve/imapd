package imap

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"compile/lexical"
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
		lexical.Parse(client.Msg)
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
