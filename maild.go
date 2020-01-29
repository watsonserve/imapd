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
