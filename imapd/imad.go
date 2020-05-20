package imapd

import (
    "errors"
    "fmt"
    "github.com/watsonserve/maild/lib"
    "net"
    "os"
    "regexp"
)

type Imapd struct {
	lib.TcpServer
    lib.ServerConfig
    re            *regexp.Regexp
    dbc           *DataAccessLayer
    outPermission map[string]bool
}

func New(dbc *DataAccessLayer, domain string, ip string) *Imapd {
    ret := &Imapd {
        dbc: dbc,
        re: regexp.MustCompile("<(.+)>"),
    }
    ret.ServerConfig = lib.ServerConfig {
        Domain: domain,
        Type: "IMAP",
        Name: "WS_IMAPD",
        Version: "1.0",
        Ip: ip,
    }

    ret.outPermission = map[string]bool {
        "AUTHENTICATE": true,
        "LOGIN": true,
        "LOGOUT": true,
        "CAPABILITY": true,
        "HELP": true,
        "NOOP": true,
        "QUIT": true,
        "RSET": true,
        "STARTTLS": true,
        "XCLIENT": true,
    }
    return ret
}

func (this *Imapd) Listen(port string) {
    err := this.TcpServer.Listen(port, this)
    if nil != err {
        fmt.Fprintln(os.Stderr, err)
    }
}

func (this *Imapd) TLSListen(port string, crt string, key string) {
    err := this.TcpServer.TLSListen(port, crt, key, this)
    if nil != err {
        fmt.Fprintln(os.Stderr, err)
    }
}

func (this *Imapd) Task(conn net.Conn) {
    ctx := InitImapContext(this.dbc, conn)
    fmt.Println("hello client")
    ctx.Send(this.Hola())

    for {
        err := ctx.Next(this)
        if nil != err {
            fmt.Fprintln(os.Stderr, "reading standard input:", err)
            break
        }
    }
}

func commandHash(this *Imapd, ctx *ImapContext, script *Mas) error {
    // 鉴权
    signIn := 0 < ctx.State
    need := this.needPermission(script.Command)
    if !signIn && need {
        ctx.Send(fmt.Sprintf("%s BAD Command received in Invalid state.", script.Tag))
        return nil
    }

    // 查找处理方法并执行处理
    switch script.Command {
        case "CAPABILITY":
            ctx.CAPABILITY(script.Tag)
        case "AUTHENTICATE":
            ctx.AUTHENTICATE(this, script)
        case "LOGIN":
            ctx.LOGIN(this, script)
        case "LOGOUT":
            ctx.LOGOUT(script.Tag)
        case "NOOP":
            ctx.NOOP(script.Tag)
        case "SELECT":
            fallthrough
        case "EXAMINE":
            ctx.Select(this, script)
        case "CLOSE":
            ctx.CLOSE(script)
        case "CREATE":
            ctx.CREATE(script)
        case "DELETE":
            ctx.DELETE(script)
        case "RENAME":
            ctx.RENAME(script)
        case "SUBSCRIBE":
            ctx.SUBSCRIBE(script)
        case "UNSUBSCRIBE":
            ctx.UNSUBSCRIBE(script)
        case "LIST":
            ctx.LIST(script)
        case "LSUB":
            ctx.LSUB(script)
        case "STATUS":
            ctx.STATUS(script)
        case "APPEND":
            ctx.APPEND(script)
        case "ID":
            ctx.ID(script)
        case "SEARCH":
            ctx.SEARCH(script)
        default:
            ctx.Send(fmt.Sprintf("%s BAD %s is not supported.\r\n", script.Tag, script.Command))
            return errors.New("method " + script.Command + " not valid")
    }
    return nil
}


func (this *Imapd) Hola() string {
    return "* OK " + this.Name + " IMAP4 service is ready.\r\n"
}

func (this *Imapd) needPermission(command string) bool {
    _, exist := this.outPermission[command]
    return !exist
}
