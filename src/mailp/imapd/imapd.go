package imapd

import (
    "bufio"
    "encoding/base64"
    "fmt"
    "net"
    "os"
    "reflect"
    "regexp"
    "strings"
    "time"
)
/*
const (
    BUFSIZ = 8192
)
*/

type Database struct {
}

func (this *Database) Auth(username string, password string) string {
    return ""
}

type Imapd struct {
    Domain string
    Type string
    Name string
    Version string
    Ip string
    re *regexp.Regexp
}

func Init(domain string, ip string) *Imapd {
    ret := &Imapd{}
    ret.Domain  = domain
    ret.Type    = "SMTP"
    ret.Name    = "WS_IMAPD"
    ret.Version = "1.0"
    ret.Ip      = ip
    ret.re      = regexp.MustCompile("<(.+)>")
    return ret
}

// 问候语
func (this *Imapd) Hola() string {
    return fmt.Sprintf("* OK %s IMAP4Server Ready\r\n", this.name)
}

// helo命令
func (this *Imapd) HELO(client *smtpContext.SmtpContext) bool {
    client.Module = smtpContext.MOD_COMMAND
    addr := client.Address
    name := client.Msg[5:]
    client.Send("250 " + this.Domain + " Hello " + name + " (" + addr + "[" + addr + "])\r\n")
    return false
}

// ehlo命令
func (this *Imapd) EHLO(client *smtpContext.SmtpContext) bool {
    client.Module = smtpContext.MOD_COMMAND
    addr := client.Address
    name := client.Msg[5:]
    client.Send("250-" + this.Domain + " Hello " + name + " (" + addr + "[" + addr + "])\r\n250-AUTH LOGIN PLAIN\r\n250-AUTH=LOGIN PLAIN\r\n250-PIPELINING\r\n250 ENHANCEDSTATUSCODES\r\n")
    return false
}

// 授权
func (this *Imapd) AUTH(client *smtpContext.SmtpContext) bool {
    content, err := base64.StdEncoding.DecodeString(client.Msg[11:])
    if nil == err {
        // debug
        fmt.Println("imap.go 75: " + string(content))
        //fmt.Fprintln(os.Stderr, "error:" + err)
        return false
    }

    for i := 0; i < len(content); i++ {
        if 0 == content[i] {
            content[i] = '\n'
        }
    }
    db := &Database{}
    userPassword := strings.Split(string(content), "\n")
    userId := db.Auth(userPassword[0], userPassword[1])
    buf := "535 Authentication Failed"
    if "" != userId {
        client.User = userPassword[0]
        buf = "235 Authentication Successful"
        client.Login = true
        fmt.Println("auth by self")
    }
    client.Send(buf)
    return false
}

// 
func (this *Imapd) QUIT(client *smtpContext.SmtpContext) bool {
    client.End("221 2.0.0 " + this.Domain + " Service closing transmission channel\r\n")
    fmt.Println(client.Head)
    fmt.Println(client.MailContent)
    return false
}

// 
func (this *Imapd) XCLIENT(client *smtpContext.SmtpContext) bool {
    fmt.Println("auth by agency")
    client.Login = true
    client.Send(this.Hola())
    return false
}

// 
func (this *Imapd) STARTTLS(client *smtpContext.SmtpContext) bool {
    client.Send("502 5.3.3 STARTTLS is not supported\r\n")
    return false
}

// 
func (this *Imapd) HELP(client *smtpContext.SmtpContext) bool {
    client.Send("502 5.3.3 HELP is not supported\r\n")
    return false
}

// 
func (this *Imapd) NOOP(client *smtpContext.SmtpContext) bool {
    client.Send("250 2.0.0 OK\r\n")
    return false
}

// 
func (this *Imapd) RSET(client *smtpContext.SmtpContext) bool {
    client.Send("250 2.0.0 OK\r\n")
    return false
}

// 
func (this *Imapd) MAIL(client *smtpContext.SmtpContext) bool {
    client.Sender = this.re.FindStringSubmatch(client.Msg)[1]
    clientDomain := strings.Split(client.Sender, "@")[1]
    if (clientDomain == this.Domain) != (!client.Login) { // 本域已登录 or 外域未登录
        client.Send("250 2.1.0 Sender <" + client.Sender + "> OK\r\n")
        return false
    }
    client.Send("530 5.7.1 Authentication Required\r\n")
    return false
}

// 
func (this *Imapd) RCPT(client *smtpContext.SmtpContext) bool {
    recver := this.re.FindStringSubmatch(client.Msg)[1]
    if strings.Split(recver, "@")[1] != this.Domain && !client.Login { // 非登录用户 to 外域
        client.Send("530 5.7.1 Authentication Required\r\n")
        return false
    }
    //fmt.Println(strings.Split(recver, "@")[1] + " ", !client.Login)
    client.Recver.PushBack(recver)
    client.Send("250 2.1.5 Recipient <" + recver + "> OK\r\n")
    return false
}

// 
func (this *Imapd) DATA(client *smtpContext.SmtpContext) bool {
    format := "from %s ([%s]) by %s over TLS secured channel with %s(%s)\r\n\t%d"
    client.Module = smtpContext.MOD_HEAD
    ele := &smtpContext.KV {
        Name: "Received",
        Value: fmt.Sprintf(format, this.Domain, this.Ip, this.Domain, this.Name, this.Version, time.Now().Unix()),
    }
    client.Head = append(client.Head, *ele)

    client.Send("354 Ok Send data ending with <CRLF>.<CRLF>\r\n")
    return false
}

func (this *Imapd) DataHead(client *smtpContext.SmtpContext) {
    if "" == client.Msg {
        client.Module = smtpContext.MOD_BODY
    } else if ' ' == client.Msg[0] || '\t' == client.Msg[0] {
        client.Head[len(client.Head) - 1].Value += "\r\n" + client.Msg
    } else {
        attr := strings.Split(client.Msg, ": ")
        ele := &smtpContext.KV {
            Name: attr[0],
            Value: attr[1],
        }
        client.Head = append(client.Head, *ele)
    }
}

func (this *Imapd) DataBody(client *smtpContext.SmtpContext) {
    if "." == client.Msg {
        client.Module = smtpContext.MOD_COMMAND
        client.Send("250 2.6.0 Message received\r\n")
        fmt.Println("250 2.6.0 Message received")
        return
    }
    client.MailContent += client.Msg + "\r\n"
}


func (this *Imapd) CommandHash(client *smtpContext.SmtpContext) bool {
    var key string
    // 截取第一个单词
    _, err := fmt.Sscanf(client.Msg, "%s", &key)
    if nil != err {
        fmt.Fprintln(os.Stderr, err)
        return true
    }
    // 查找处理方法
    that := reflect.ValueOf(this)
    method := that.MethodByName(key)
    if !method.IsValid() {
        fmt.Fprintln(os.Stderr, "method " + key + "not valid")
        return true
    }
    // 执行处理
    clientValue := reflect.ValueOf(client)
    inArgs := []reflect.Value{clientValue}
    return method.Call(inArgs)[0].Bool()
}

func (this *Imapd) Task(conn net.Conn) {
    scanner := bufio.NewScanner(conn)
    client := smtpContext.InitSmtpContext(conn)
    client.Send(this.Hola())

    for scanner.Scan() {
        err := scanner.Err()
        if nil != err {
            fmt.Fprintln(os.Stderr, "reading standard input:", err)
            break
        }
        // 行遍历
        client.Msg = scanner.Text()
        //if "QUIT" == client.Msg {
        //}
        fmt.Println(client.Msg)
        switch client.Module {
            case smtpContext.MOD_COMMAND:
                if this.CommandHash(client) {
                }
            case smtpContext.MOD_HEAD:
                this.DataHead(client)
            case smtpContext.MOD_BODY:
                this.DataBody(client)
        }
    }
}
