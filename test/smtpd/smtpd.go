package main

import (
    "os"
    "io"
    "log"
    "fmt"
    "github.com/watsonserve/maild"
    "github.com/watsonserve/maild/smtpd"
)

type SmtpConf struct {}

func (this *SmtpConf) GetConfig() *maild.ServerConfig {
    return &maild.ServerConfig{
        Domain: "watsonserve.com",
        Ip: "127.0.0.1",
        Type: "SMTP",
        Name: "WS_SMTPD",
        Version: "1.0",
    }
}

func (this *SmtpConf) Auth(username string, password string) string {
    fmt.Printf("%s %s\n", username, password)
    return "null"
}

func (this *SmtpConf) TakeOff(email *maild.Mail) {
    fmt.Println(email.Head)
    fmt.Println(email.MailContent)
}

func main() {
    //*/
    fp := os.Stderr
    /*/
    fp, err := os.OpenFile("/var/log/mail_auth.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
    if nil != err {
        log.Fatal(err)
        return
    }
    //*/
    log.SetOutput(io.Writer(fp))
    log.SetFlags(log.Ldate|log.Ltime|log.Lmicroseconds)

    smtpServer := smtpd.New(&SmtpConf {})

    fmt.Println("listen on port 10025")
    smtpServer.Listen(":10025")
}
