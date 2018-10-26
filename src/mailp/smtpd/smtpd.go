package smtpd

import (
    "auth"
    "bufio"
    "encoding/base64"
    "errors"
    "fmt"
    "net"
    "os"
    "reflect"
    "regexp"
    "strings"
    "time"
    "mailp/smtpd/smtpContext"
)
/*
const (
    BUFSIZ = 8192
)
*/

type Smtpd struct {
    Domain string
    Type string
    Name string
    Version string
    Ip string
    re *regexp.Regexp
}

func Init(domain string, ip string) *Smtpd {
    ret := &Smtpd{}
    ret.Domain  = domain
    ret.Type    = "SMTP"
    ret.Name    = "WS_SMTPD"
    ret.Version = "1.0"
    ret.Ip      = ip
    ret.re      = regexp.MustCompile("<(.+)>")
    return ret
}

// 问候语
func (this *Smtpd) Hola() string {
    return fmt.Sprintf("220 %s %s Server (%s %s Server %s) ready %d\r\n",
        this.Domain, this.Type, this.Name, this.Type, this.Version, time.Now().Unix(),
    )
}

// helo命令
func (this *Smtpd) HELO(client *smtpContext.SmtpContext) {
    client.Module = smtpContext.MOD_COMMAND
    addr := client.Address
    name := client.Msg[5:]
    client.Send("250 " + this.Domain + " Hello " + name + " (" + addr + "[" + addr + "])\r\n")
}

// ehlo命令
func (this *Smtpd) EHLO(client *smtpContext.SmtpContext) {
    client.Module = smtpContext.MOD_COMMAND
    addr := client.Address
    name := client.Msg[5:]
    client.Send("250-" + this.Domain + " Hello " + name + " (" + addr + "[" + addr + "])\r\n250-AUTH LOGIN PLAIN\r\n250-AUTH=LOGIN PLAIN\r\n250-PIPELINING\r\n250 ENHANCEDSTATUSCODES\r\n")
}

// 授权
func (this *Smtpd) AUTH(client *smtpContext.SmtpContext) {
    content, err := base64.StdEncoding.DecodeString(client.Msg[11:])
    if nil != err {
        fmt.Fprintln(os.Stderr, "error: " + err.Error())
        return
    }

    for i := 0; i < len(content); i++ {
        if 0 == content[i] {
            content[i] = '\n'
        }
    }
    author := auth.New()
    userPassword := strings.Split(string(content), "\n")
    userId := author.Auth(userPassword[0], userPassword[1])
    buf := "535 Authentication Failed\r\n"
    if "" != userId {
        client.User = userPassword[0]
        buf = "235 Authentication Successful\r\n"
        client.Login = true
        fmt.Println("auth by self")
    }
    client.Send(buf)
}

// 
func (this *Smtpd) QUIT(client *smtpContext.SmtpContext) {
    client.End("221 2.0.0 " + this.Domain + " Service closing transmission channel\r\n")
}

// 
func (this *Smtpd) XCLIENT(client *smtpContext.SmtpContext) {
    fmt.Println("auth by agency")
    client.Login = true
    client.Send(this.Hola())
}

// 
func (this *Smtpd) STARTTLS(client *smtpContext.SmtpContext) {
    client.Send("502 5.3.3 STARTTLS is not supported\r\n")
    fmt.Println("startTTS")
}

// 
func (this *Smtpd) HELP(client *smtpContext.SmtpContext) {
    client.Send("502 5.3.3 HELP is not supported\r\n")
}

// 
func (this *Smtpd) NOOP(client *smtpContext.SmtpContext) {
    client.Send("250 2.0.0 OK\r\n")
    fmt.Println("noop")
}

// 
func (this *Smtpd) RSET(client *smtpContext.SmtpContext) {
    client.Send("250 2.0.0 OK\r\n")
    fmt.Println("rset")
}

// 
func (this *Smtpd) MAIL(client *smtpContext.SmtpContext) {
    client.Sender = this.re.FindStringSubmatch(client.Msg)[1]
    clientDomain := strings.Split(client.Sender, "@")[1]
    if (clientDomain == this.Domain) != (!client.Login) { // 本域已登录 or 外域未登录
        client.Send("250 2.1.0 Sender <" + client.Sender + "> OK\r\n")
        return
    }
    client.Send("530 5.7.1 Authentication Required\r\n")
}

// 
func (this *Smtpd) RCPT(client *smtpContext.SmtpContext) {
    recver := this.re.FindStringSubmatch(client.Msg)[1]
    if strings.Split(recver, "@")[1] != this.Domain && !client.Login { // 非登录用户 to 外域
        client.Send("530 5.7.1 Authentication Required\r\n")
        return
    }
    //fmt.Println(strings.Split(recver, "@")[1] + " ", !client.Login)
    client.Recver.PushBack(recver)
    client.Send("250 2.1.5 Recipient <" + recver + "> OK\r\n")
}

// 
func (this *Smtpd) DATA(client *smtpContext.SmtpContext) {
    format := "from %s ([%s]) by %s over TLS secured channel with %s(%s)\r\n\t%d"
    client.Module = smtpContext.MOD_HEAD
    ele := &smtpContext.KV {
        Name: "Received",
        Value: fmt.Sprintf(format, this.Domain, this.Ip, this.Domain, this.Name, this.Version, time.Now().Unix()),
    }
    client.Head = append(client.Head, *ele)

    client.Send("354 Ok Send data ending with <CRLF>.<CRLF>\r\n")
}

func (this *Smtpd) DataHead(client *smtpContext.SmtpContext) {
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

func (this *Smtpd) DataBody(client *smtpContext.SmtpContext) {
    if "." == client.Msg {
        client.Module = smtpContext.MOD_COMMAND
        client.Send("250 2.6.0 Message received\r\n")
        client.TakeOff()
        fmt.Println("250 2.6.0 Message received")
        return
    }
    client.MailContent += client.Msg + "\r\n"
}


func (this *Smtpd) CommandHash(client *smtpContext.SmtpContext) error {
    var key string
    // 截取第一个单词
    _, err := fmt.Sscanf(client.Msg, "%s", &key)
    if nil != err {
        return err
    }
    // 查找处理方法
    that := reflect.ValueOf(this)
    method := that.MethodByName(key)
    if !method.IsValid() {
        return errors.New("method " + key + " not valid")
    }
    // 执行处理
    clientValue := reflect.ValueOf(client)
    inArgs := []reflect.Value{clientValue}
    method.Call(inArgs)
    return nil
}

func (this *Smtpd) Task(conn net.Conn) {
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
                err = this.CommandHash(client)
                if nil != err {
                    fmt.Fprintln(os.Stderr, err)
                }
            case smtpContext.MOD_HEAD:
                this.DataHead(client)
            case smtpContext.MOD_BODY:
                this.DataBody(client)
        }
    }
}
