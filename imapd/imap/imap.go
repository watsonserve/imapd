package imap

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/watsonserve/maild/compile/lexical"
    "net"
    "os"
    "reflect"
    "regexp"
)

/*
const (
    BUFSIZ = 8192
)
*/


func (this *Imapd) Hola() string {
    return "OK " + this.Name + " IMAP4 service is ready.\r\n"
}

func (this *Imapd) needPermission(command string) bool {
    _, exist := this.outPermission[command]
    return !exist
}

func (this *Imapd) CAPABILITY(ctx *ImapContext, mas *Mas) {
    abilities := ""
    length := len(this.capability)

    for i := 0; i < length; i++ {
        item := this.capability[i]
        if !ctx.Login && item.Permission {
            continue
        }
        abilities += " " + item.Ability
    }

    ctx.Send(fmt.Sprintf(
        "* CAPABILITY IMAP4%s\r\n%d OK CAPABILITY completed.\r\n",
        abilities, mas.Count,
    ))
}

func (this *Imapd) LOGIN(ctx *ImapContext, mas *Mas) {
    ctx.Login = true
    ctx.Send(fmt.Sprintf("%d OK LOGIN completed.\r\n", mas.Count))
}
