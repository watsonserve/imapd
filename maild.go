package maild

import (
	"container/list"
	"strings"
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

type ServerConfig struct {
	Domain  string
	Ip      string
	Name    string
	Type    string
	Version string
}

type ImapAgentConfigure interface {
	GetConfig() *ServerConfig
	Auth(username string, password string) string
	Trans(email *Mail)
}

type SmtpServerConfigure interface {
	GetConfig() *ServerConfig
	Auth(username string, password string) string
	TakeOff(email *Mail)
}

func Split(raw []byte, sp byte) []string {
	var ret []string
	length := len(raw)
	dest := make([]byte, length)
	for i := 0; i < length; i++ {
		ch := raw[i]
		if 0 == ch {
			ch = '\n'
		}
		dest[i] = ch
	}
	list := strings.Split(string(dest), "\n")
	length = len(list)

	for i := 0; i < length; i++ {
		if "" == list[i] {
			continue
		}
		ret = append(ret, list[i])
	}
	return ret
}
