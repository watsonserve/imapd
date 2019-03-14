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
