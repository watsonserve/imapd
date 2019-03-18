package lib

import (
	"container/list"
	"strings"
)

type Author interface {
    Auth(username string, password string) string
}

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
    Author  Author
	Domain  string
	Ip      string
	Name    string
	Type    string
	Version string
}

func (this *ServerConfig) CloneFrom(ctx *ServerConfig) {
	this.Author = ctx.Author
	this.Domain = ctx.Domain
	this.Ip = ctx.Ip
	this.Name = ctx.Name
	this.Type = ctx.Type
	this.Version = ctx.Version
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
