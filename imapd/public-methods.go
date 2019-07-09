package imapd

import (
	"fmt"
)

func (this *ImapContext) NOOP(tag string) {
	// check state
	this.Send(fmt.Sprintf("%s OK NOOP completed.\r\n", tag))
}

func (this *ImapContext) RSET(script *Mas) {
	this.Send(fmt.Sprintf("%s OK RSET completed.\r\n", script.Tag))
}

func (this *ImapContext) CAPABILITY(tag string) {
    abilities := "IMAP4rev1 LOGINDISABLED"

    if 0 == this.State {
        abilities += " AUTH=PLAIN"
    }

    this.Send(fmt.Sprintf("* CAPABILITY IMAP4%s\r\n%s OK CAPABILITY completed.\r\n", abilities, tag))
}
