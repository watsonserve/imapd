package imap

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/watsonserve/maild"
    "github.com/watsonserve/maild/compile/lexical"
	"github.com/watsonserve/maild/server"
    "net"
    "os"
    "reflect"
    "regexp"
)

/*
const (
    BUFSIZ = 8192
)
*/

type cap_t struct {
    Ability    string
    Permission bool
}

type Imapd struct {
	server.TcpServer
    maild.ServerConfig
    re            *regexp.Regexp
    capability    []cap_t
    outPermission map[string]bool
}

func New(domain string, ip string) *Imapd {
    ret := &Imapd{}
    ret.Domain = domain
    ret.Type = "IMAP"
    ret.Name = "WS_IMAPD"
    ret.Version = "1.0"
    ret.Ip = ip
    ret.re = regexp.MustCompile("<(.+)>")
    ret.capability = []cap_t {
        {Ability: "IMAP4rev1", Permission: false},
        {Ability: "AUTH=PLAIN", Permission: false},
        {Ability: "AUTH=XOAUTH2", Permission: false},
        {Ability: "SASL-IR", Permission: false},
        {Ability: "UIDPLUS", Permission: false},
        {Ability: "MOVE", Permission: false},
        {Ability: "ID", Permission: false},
        {Ability: "UNSELECT", Permission: false},
        {Ability: "CLIENTACCESSRULES",Permission: true},
        {Ability: "CLIENTNETWORKPRESENCELOCATION",Permission: true},
        {Ability: "BACKENDAUTHENTICATE",Permission: true},
        {Ability: "CHILDREN", Permission: false},
        {Ability: "IDLE", Permission: false},
        {Ability: "NAMESPACE", Permission: false},
        {Ability: "LITERAL+", Permission: false},
    }
    ret.outPermission = map[string]bool {
        "LOGIN": true,
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

func (this *Imapd) Task(conn net.Conn) {
    scanner := bufio.NewScanner(conn)
    ctx := InitImapContext(conn)
    imapd.Hola()

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
        ctx.Send(fmt.Sprintf("%d BAD Command received in Invalid state.", script.Count))
        return nil
    }

    // 查找处理方法
    that := reflect.ValueOf(this)
    method, exist := this.dict[script.Command]
    if !exist {
        return errors.New("method " + script.Command + "not valid")
    }

    // 执行处理
    method(ctx, script)
    return nil
}


func (this *Imapd) Hola() string {
    return "OK " + this.Name + " IMAP4 service is ready.\r\n"
}

func (this *Imapd) needPermission(command string) bool {
    _, exist := this.outPermission[command]
    return !exist
}
