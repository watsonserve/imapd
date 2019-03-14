package imapd

import (
    "fmt"
)

func (this *ImapContext) CAPABILITY(script *Mas) {
    abilities := "IMAP4rev1 LOGINDISABLED "

    if this.Login {
        abilities += "SASL-IR UIDPLUS MOVE ID UNSELECT CLIENTACCESSRULE CLIENTNETWORKPRESENCELOCATIO BACKENDAUTHENTICAT CHILDREN IDLE NAMESPACE LITERAL+"
    } else {
        abilities += "AUTH=PLAIN AUTH=XOAUTH2"
    }

    this.Send(fmt.Sprintf("* CAPABILITY IMAP4%s\r\n%s OK CAPABILITY completed.\r\n", abilities, mas.Tag))
}

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

func (this *ImapContext) HELO(script *Mas) {
	this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) EHLO(script *Mas) {
	this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) SELECT(script *Mas) {
    this.MailBox = ""
    if 1 != len(script.Parames) {
        this.Send(fmt.Sprintf("%s BAD SELECT FAILURE: arguments invalid.\r\n", script.Tag))
        return
    }
    mailBox := script.Parames[0]
    sql, _ := that.StmtMap["select"]
    row := sql.QueryRow(mailBox, this.User)
    err := row.Scan(mailBox)
    if nil != err {
        this.Send(fmt.Sprintf("%s NO SELECT FAILURE: no such mailbox.\r\n", script.Tag))
        return
    }
    this.MailBox = mailBox
    sum, recentCnt, unseen, validCode, nextUID
    fmt.Sprintf("* %d FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)\r\n")
    fmt.Sprintf("* %d EXISTS\r\n", sum)
    fmt.Sprintf("* %d RECENT\r\n", recentCnt)
    fmt.Sprintf("* OK [UNSEEN %d] Message %d is first unseen\r\n", unseen, unseen)
    fmt.Sprintf("* OK [UIDVALIDITY %d] UIDs valid\r\n", validCode)
    fmt.Sprintf("* OK [UIDNEXT %d] Predicted next UID\r\n", nextUID)

    // A list of message flags that the client can change permanently.  If this is missing, the client should assume that all flags can be changed permanently.
    fmt.Sprintf("* OK [PERMANENTFLAGS (\\Deleted \\Seen \\*)] Limited\r\n")
    this.Send(fmt.Sprintf("%s OK [READ-WRITE] SELECT completed.\r\n", script.Tag))
}

func (this *ImapContext) EXAMINE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) CREATE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) DELETE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) RENAME(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) SUBSCRIBE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) UNSUBSCRIBE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) LIST(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) LSUB(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) STATUS(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) APPEND(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) EXPUNGE(script *Mas) {
    this.Send(fmt.Sprintf("* %d EXPUNGE\r\n", serial))
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}
