package imapd

import (
    "errors"
    "fmt"
    "net"
    "os"
    "github.com/watsonserve/maild/lib"
)

type ImapContext struct {
    lib.ReadStream
    lib.SentStream
    lib.ServerConfig
    rw      int
    User    string
    State   int
    mailBox string
    dal     *DataAccessLayer
    Email   *lib.Mail
}

func InitImapContext(dbc *DataAccessLayer, sock net.Conn) *ImapContext {
    ret := &ImapContext{
        rw: 0,
        User: "",
        State: 0,
        mailBox: "",
        dal: dbc,
        Email: &lib.Mail{},
    }

    ret.ReadStream = *lib.InitReadStream(sock)
    ret.SentStream = *lib.InitSentStream(sock)
    return ret
}

func (this *ImapContext) Next(imapd *Imapd) error {
    msg, err := this.ReadLine()
    if nil != err {
        return err
    }
    if "" == msg {
        return errors.New("EOF")
    }
    script := initMas(msg)
    if nil == script {
        this.Send(fmt.Sprintf("%s BAD Command Error.\r\n", msg))
        return nil
    }

    err = commandHash(imapd, this, script)
    if nil != err {
        fmt.Fprintln(os.Stderr, err)
    }
    return nil
}
