package imapd

import (
    "fmt"
)

func (this *ImapContext) XCLIENT(script *Mas) {
	this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) signIn(that *Imapd, tag string, parames []string) {
    for 2 == len(parames) {
        userId := that.Author.Auth(parames[0], parames[1])
        if "" == userId {
            break
        }
        this.User = userId
        this.State = 1
        this.Send(fmt.Sprintf("%s OK LOGIN completed.\r\n", tag))
        return
    }
    this.Send(fmt.Sprintf("%s NO LOGIN FAILURE: username or password rejected.\r\n", tag))
}

func (this *ImapContext) LOGIN(that *Imapd, script *Mas) {
    parames := strings.Split(script.Parames, " ")
    this.signIn(that, script.Tag, parames)
}

func (this *ImapContext) AUTHENTICATE(that *Imapd, script *Mas) {
    if "PLAIN" != script.Parames {
	    this.Send(fmt.Sprintf("%s NO AUTHENTICATE FAILURE.\r\n", script.Tag))
        return
    }
    this.Send(fmt.Sprintf("+\r\n"))
    msg, err := this.ReadLine()
    if nil != err {
        return
    }
    if "*" == msg {
        this.Send(fmt.Sprintf("%s NO The AUTH protocol exchange was canceled by the client.\r\n", script.Tag))
        return
    }

    decodeContent, err := base64.StdEncoding.DecodeString(msg)
    if nil != err {
        this.Send(fmt.Sprintf("%s NO arguments invalid.\r\n", script.Tag))
        return
    }
    up := lib.Split(decodeContent, '\x00')
    this.signIn(that, script.Tag, up)
}

func (this *ImapContext) LOGOUT(tag string) {
    this.State = 0
	this.End(fmt.Sprintf("* BYE IMAP4rev1 Server logging out\r\n%s OK LOGOUT completed.\r\n", tag))
}
