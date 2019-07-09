package imapd

import (
    "fmt"
)

func (this *ImapContext) select(script *Mas) {
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
    this.mailBox = mailBox
    sum, recentCnt, unseen, validCode, nextUID

    // read write flag
    this.rw = 1
    rw := "READ-ONLY"
    if "SELECT" == script.Command {
        this.rw = 3
        rw := "READ-WRITE"
    }
    // response text
    resp := fmt.Sprintf("* %d FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)\r\n")
    resp += fmt.Sprintf("* %d EXISTS\r\n", sum)
    resp += fmt.Sprintf("* %d RECENT\r\n", recentCnt)
    resp += fmt.Sprintf("* OK [UNSEEN %d] Message %d is first unseen\r\n", unseen, unseen)
    resp += fmt.Sprintf("* OK [UIDVALIDITY %d] UIDs valid\r\n", validCode)
    resp += fmt.Sprintf("* OK [UIDNEXT %d] Predicted next UID\r\n", nextUID)

    // A list of message flags that the client can change permanently.  If this is missing, the client should assume that all flags can be changed permanently.
    resp += fmt.Sprintf("* OK [PERMANENTFLAGS (\\Deleted \\Seen \\*)] Limited\r\n")
    this.Send(resp + fmt.Sprintf("%s OK [%s] %s completed.\r\n", script.Tag, rw, script.Command))
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
