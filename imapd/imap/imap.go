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

func Init(domain string, ip string) *Imapd {
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

func (this *Imapd) CommandHash(client *ImapContext, msg string) error {
    script := initMas(lexical.Parse(msg))

    if !client.Login && this.needPermission(script.Command) {
        client.Send(fmt.Sprintf("%d BAD Command received in Invalid state.", script.Count))
        return nil
    }
    // 查找处理方法
    that := reflect.ValueOf(this)
    method := that.MethodByName(script.Command)
    if !method.IsValid() {
        return errors.New("method " + script.Command + "not valid")
    }
    // 执行处理
    clientValue := reflect.ValueOf(client)
    scriptValue := reflect.ValueOf(script)
    inArgs := []reflect.Value{clientValue, scriptValue}
    method.Call(inArgs)
    return nil
}

func (this *Imapd) Task(conn net.Conn) {
    scanner := bufio.NewScanner(conn)
    client := InitImapContext(conn)
    client.Send(this.Hola())

    for scanner.Scan() {
        err := scanner.Err()
        if nil != err {
            fmt.Fprintln(os.Stderr, "reading standard input:", err)
            break
        }
        msg := scanner.Text()
        fmt.Println(msg)

        err = this.CommandHash(client, msg)
        if nil != err {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}

func (this *Imapd) CAPABILITY(client *ImapContext, mas *Mas) {
    abilities := ""
    length := len(this.capability)

    for i := 0; i < length; i++ {
        item := this.capability[i]
        if !client.Login && item.Permission {
            continue
        }
        abilities += " " + item.Ability
    }

    client.Send(fmt.Sprintf(
        "* CAPABILITY IMAP4%s\r\n%d OK CAPABILITY completed.\r\n",
        abilities, mas.Count,
    ))
}

func (this *Imapd) LOGIN(client *ImapContext, mas *Mas) {
    client.Login = true
    client.Send(fmt.Sprintf("%d OK LOGIN completed.\r\n", mas.Count))
}
