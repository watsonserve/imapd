package imapd

import (
    "database/sql"
    "fmt"
    "net"
    "github.com/watsonserve/maild/lib"
)

type ImapContext struct {
    lib.SentStream
    lib.ServerConfig
    Login   bool

    rw      int
    mailBox string
    db      *sql.DB
    StmtMap map[string]*sql.Stmt

    User    string
    Email   *lib.Mail
}


func InitImapContext(sock net.Conn) *ImapContext {
    ret := &ImapContext{Login: false}

    ret.SentStream = *lib.InitSentStream(sock)
    ret.Email = &lib.Mail{}
    ret.db = db
    ret.StmtMap = make(map[string]*sql.Stmt)
    ret.Prepare("select", "select box_id where name=? and user_id=?")
    return ret
}

func (this *ImapContext) Prepare(index string, query string) {
    stmt ,err := this.db.Prepare(query)
    if nil != err {
        panic(err)
    }
    this.StmtMap[index] = stmt
}
