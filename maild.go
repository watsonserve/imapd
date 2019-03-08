package maild

import (
	"container/list"
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

func (this *ServerConfig) CloneFrom(ctx *ServerConfig) {
	this.Domain = ctx.Domain
	this.Ip = ctx.Ip
	this.Name = ctx.Name
	this.Type = ctx.Type
	this.Version = ctx.Version
}
