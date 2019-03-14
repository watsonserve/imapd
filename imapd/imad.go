package imapd

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/watsonserve/maild/compile/lexical"
    "github.com/watsonserve/maild/lib"
    "net"
    "os"
    "regexp"
)

type Imapd struct {
	lib.TcpServer
    lib.ServerConfig
    re            *regexp.Regexp
    outPermission map[string]bool
}

func New(db *sql.DB, domain string, ip string) *Imapd {
    ret := &Imapd{}
    ret.Domain = domain
    ret.Type = "IMAP"
    ret.Name = "WS_IMAPD"
    ret.Version = "1.0"
    ret.Ip = ip
    ret.re = regexp.MustCompile("<(.+)>")
    ret.outPermission = map[string]bool {
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
	this.TcpServer.Listen(port, this)
}

func (this *Imapd) Task(conn net.Conn) {
    scanner := bufio.NewScanner(conn)
    ctx := InitImapContext(conn)
	ctx.CloneFrom(&this.ServerConfig)
    ctx.Send(this.Hola())

    for scanner.Scan() {
        err := scanner.Err()
        if nil != err {
            fmt.Fprintln(os.Stderr, "reading standard input:", err)
            break
        }
        msg := scanner.Text()
        script := initMas(lexical.Parse(msg))

        err = commandHash(this, ctx, script)
        if nil != err {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}

func commandHash(this *Imapd, ctx *ImapContext, script *Mas) error {
    // 鉴权
    if !ctx.Login && this.needPermission(script.Command) {
        ctx.Send(fmt.Sprintf("%s BAD Command received in Invalid state.", script.Tag))
        return nil
    }

    // 查找处理方法并执行处理
    switch script.Command {
        case "CAPABILITY":
            ctx.CAPABILITY()
        case "XCLIENT":
            ctx.XCLIENT()
        case "LOGIN":
            ctx.LOGIN(script)
        case "LOGOUT":
            ctx.LOGOUT()
        case "HELO":
            ctx.HELO()
        case "EHLO":
            ctx.EHLO()
        case "SELECT":
            fallthrough
        case "EXAMINE":
            ctx.select(script)
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
        default:
            ctx.Send(fmt.Sprintf("%s BAD %s is not supported.\r\n", script.Tag, script.Command))
            return errors.New("method " + script.Command + " not valid")
    }
    return nil
}


func (this *Imapd) Hola() string {
    return "OK " + this.Name + " IMAP4 service is ready.\r\n"
}

func (this *Imapd) needPermission(command string) bool {
    _, exist := this.outPermission[command]
    return !exist
}
