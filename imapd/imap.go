package imapd

import (
    "encoding/base64"
    "fmt"
    "github.com/watsonserve/maild/compile/lexical"
    "github.com/watsonserve/maild/lib"
    "strings"
)

func (this *ImapContext) CAPABILITY(tag string) {
    abilities := "IMAP4rev1 LOGINDISABLED"

    if 0 == this.State {
        abilities += " AUTH=PLAIN"
    }

    this.Send(fmt.Sprintf("* CAPABILITY IMAP4%s\r\n%s OK CAPABILITY completed.\r\n", abilities, tag))
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

func (this *ImapContext) Select(that *Imapd, script *Mas) {
    this.mailBox = ""
    if "" == script.Parames || "\"\"" == script.Parames {
        this.Send(fmt.Sprintf("%s BAD %s FAILURE: arguments invalid.\r\n", script.Tag, script.Command))
        return
    }
    // mailBox := script.Parames[0]
    // sql, _ := that.StmtMap["select"]
    // row := sql.QueryRow(mailBox, this.User)
    // err := row.Scan(mailBox)
    // if nil != err {
    //     this.Send(fmt.Sprintf("%s NO SELECT FAILURE: no such mailbox.\r\n", script.Tag))
    //     return
    // }
    this.mailBox = script.Parames
    sum := 0
    recentCnt := 0
    unseen := 0
    validCode := 1
    nextUID := 1

    // read write flag
    this.rw = 1
    rw := "READ-ONLY"
    if "SELECT" == script.Command {
        this.rw = 3
        rw = "READ-WRITE"
    }
    // response text
    resp := fmt.Sprintf("* FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)\r\n")
    resp += fmt.Sprintf("* %d EXISTS\r\n", sum)
    resp += fmt.Sprintf("* %d RECENT\r\n", recentCnt)
    resp += fmt.Sprintf("* OK [UNSEEN %d] Message %d is first unseen\r\n", unseen, unseen)
    resp += fmt.Sprintf("* OK [UIDVALIDITY %d] UIDs valid\r\n", validCode)
    resp += fmt.Sprintf("* OK [UIDNEXT %d] Predicted next UID\r\n", nextUID)

    // A list of message flags that the client can change permanently.  If this is missing, the client should assume that all flags can be changed permanently.
    resp += fmt.Sprintf("* OK [PERMANENTFLAGS (\\Deleted \\Seen \\*)] Limited\r\n")
    this.Send(resp + fmt.Sprintf("%s OK [%s] %s completed.\r\n", script.Tag, rw, script.Command))
}

func (this *ImapContext) CLOSE(script *Mas) {
    this.mailBox = ""
    this.Send(fmt.Sprintf("%s OK CLOSE completed.\r\n", script.Tag))
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
    parames := lexical.Parse(script.Parames) // [ref, mailbox]
    if 2 != len(parames) {
        this.Send(fmt.Sprintf("%s NO LIST FAILURE: arguments invalid.\r\n", script.Tag))
        return
    }
    // 只要"" == mailbox就返回"/" ""
    if "\"\"" == parames[1].Value {
        this.Send(fmt.Sprintf("* LIST (\\Noselect \\HasChildren) \"/\" \"\"\r\n", script.Tag, script.Command))
        return
    }

    paths := []string{"", ""}
    
    for i := 0; i < 2; i++ {
        val := parames[i].Value
        if '"' == val[0] {
            val = val[1:-1]
        }
        paths[i] = val
    }
    // ["", "/"]  ["", "%"]  ["", "*"]  ["", "name"]
    path := strings.Join(paths, "/")
    if '/' == path[-1] {
        this.Send(fmt.Sprintf("* LIST (\\Noselect \\HasChildren) \"/\" \"\"\r\n", script.Tag, script.Command))
    }

    this.Send(fmt.Sprintf("* LIST ().\r\n", script.Tag, script.Command))
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
    // this.Send(fmt.Sprintf("* %d EXPUNGE\r\n", serial))
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}





func (this *ImapContext) ID(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}
