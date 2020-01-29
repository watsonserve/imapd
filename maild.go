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

type SmtpServerConfigure interface {
	GetConfig() *ServerConfig
	Auth(username string, password string) string
	TakeOff(email *Mail)
}
