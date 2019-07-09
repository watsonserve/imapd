package imapd

import (
    "fmt"
)

func (this *ImapContext) XCLIENT(script *Mas) {
	this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) LOGIN(script *Mas) {
    for 2 == len(mas.Parames) {
        userId := that.Author.Auth(mas.Parames[0], mas.Parames[1])
        if "" == userId {
            break
        }
        this.User = userId
        this.Login = true
        this.Send(fmt.Sprintf("%s OK LOGIN completed.\r\n", mas.Tag))
        return
    }
    this.Send(fmt.Sprintf("%s NO LOGIN FAILURE: username or password rejected.\r\n", mas.Tag))
}

func (this *ImapContext) AUTHENTICATE(script *Mas) {
	this.Send(fmt.Sprintf("%s NO %s FAILURE.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) LOGOUT(script *Mas) {
    this.Login = false
	this.End(fmt.Sprintf("* BYE IMAP4rev1 Server logging out\r\n%s OK LOGOUT completed.\r\n", script.Tag))
}
