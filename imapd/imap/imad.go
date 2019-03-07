package imap

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/watsonserve/maild/compile/lexical"
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
    Domain        string
    Type          string
    Name          string
    Version       string
    Ip            string
    re            *regexp.Regexp
    capability    []cap_t
    outPermission map[string]bool
}

func Init() *Imapd {
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

func (this *Imapd) Hola() string {
    return "OK " + this.Name + " IMAP4 service is ready.\r\n"
}

func (this *Imapd) needPermission(command string) bool {
    _, exist := this.outPermission[command]
    return !exist
}




type Imapd struct {
    imapdCreator func() * Imapd
}

func New(domain string, ip string) *Imapd {
    ret := &Imapd{}
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
        fmt.Println(msg)

        err = commandHash(this, ctx, msg)
        if nil != err {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}

func commandHash(this *Imapd, ctx *ImapContext, msg string) error {
    script := initMas(lexical.Parse(msg))

    // 鉴权
    if !ctx.Login && this.needPermission(script.Command) {
        ctx.Send(fmt.Sprintf("%d BAD Command received in Invalid state.", script.Count))
        return nil
    }

    // 查找处理方法
    that := reflect.ValueOf(this)
    method := that.MethodByName(script.Command)
    if !method.IsValid() {
        return errors.New("method " + script.Command + "not valid")
    }

    // 执行处理
    clientValue := reflect.ValueOf(ctx)
    scriptValue := reflect.ValueOf(script)
    inArgs := []reflect.Value{clientValue, scriptValue}
    method.Call(inArgs)
    return nil
}
