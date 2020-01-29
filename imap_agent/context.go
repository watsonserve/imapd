package imap_agent

import (
    "encoding/base64"
    "fmt"
    "net"
    "github.com/watsonserve/maild"
    "strings"
)

type ImapAgentContext struct {
    maild.ReadStream
    maild.SentStream
    Sess string
}

func split(raw []byte, sp byte) []string {
	var ret []string
	length := len(raw)
	dest := make([]byte, length)
	for i := 0; i < length; i++ {
		ch := raw[i]
		if 0 == ch {
			ch = '\n'
		}
		dest[i] = ch
	}
	list := strings.Split(string(dest), "\n")
	length = len(list)

	for i := 0; i < length; i++ {
		if "" == list[i] {
			continue
		}
		ret = append(ret, list[i])
	}
	return ret
}

func InitImapAgentContext(sock net.Conn) *ImapAgentContext {
    return &ImapAgentContext {
        ReadStream: *maild.InitReadStream(sock),
        SentStream: *maild.InitSentStream(sock),
        Sess: "",
    }
}

func (this *ImapAgentContext) Checked() bool {
    return "" == this.Sess
}

func (this *ImapAgentContext) AUTHENTICATE(script *Mas) []string {
    if "PLAIN" != script.Parames {
	    this.Send(fmt.Sprintf("%s NO AUTHENTICATE FAILURE.\r\n", script.Tag))
        return nil
    }
    this.Send(fmt.Sprintf("+\r\n"))
    msg, err := this.ReadLine()
    if nil != err {
        return nil
    }
    if "*" == msg {
        this.Send(fmt.Sprintf("%s NO The AUTH protocol exchange was canceled by the client.\r\n", script.Tag))
        return nil
    }

    decodeContent, err := base64.StdEncoding.DecodeString(msg)
    if nil != err {
        this.Send(fmt.Sprintf("%s NO arguments invalid.\r\n", script.Tag))
        return nil
    }
    return split(decodeContent, '\x00')
}

func (this *ImapAgentContext) LOGOUT(tag string) {
    this.Sess = ""
	this.End(fmt.Sprintf("* BYE IMAP4rev1 Server logging out\r\n%s OK LOGOUT completed.\r\n", tag))
}

func (this *ImapAgentContext) NOOP(tag string) {
	this.Send(fmt.Sprintf("%s OK NOOP completed.\r\n", tag))
}

func (this *ImapAgentContext) RSET(script *Mas) {
	this.Send(fmt.Sprintf("%s OK RSET completed.\r\n", script.Tag))
}

func (this *ImapAgentContext) CAPABILITY(tag string) {
    abilities := "IMAP4rev1 LOGINDISABLED"

    if "" == this.Sess {
        abilities += " AUTH=PLAIN"
    }

    this.Send(fmt.Sprintf("* CAPABILITY IMAP4%s\r\n%s OK CAPABILITY completed.\r\n", abilities, tag))
}