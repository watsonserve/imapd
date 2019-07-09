package imapd

import (
	"fmt"
)

func (this *ImapContext) NOOP(script *Mas) {
	// check state
	this.Send(fmt.Sprintf("%s OK NOOP completed.\r\n", script.Tag))
}

func (this *ImapContext) RSET(script *Mas) {
	this.Send(fmt.Sprintf("%s OK RSET completed.\r\n", script.Tag))
}

func (this *ImapContext) CAPABILITY(script *Mas) {
    abilities := "IMAP4rev1 LOGINDISABLED "

    if this.Login {
        abilities += "SASL-IR UIDPLUS MOVE ID UNSELECT CLIENTACCESSRULE CLIENTNETWORKPRESENCELOCATIO BACKENDAUTHENTICAT CHILDREN IDLE NAMESPACE LITERAL+"
    } else {
        abilities += "AUTH=PLAIN AUTH=XOAUTH2"
    }

    this.Send(fmt.Sprintf("* CAPABILITY IMAP4%s\r\n%s OK CAPABILITY completed.\r\n", abilities, mas.Tag))
}

func (this *ImapContext) HELO(script *Mas) {
	this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) EHLO(script *Mas) {
	this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}
