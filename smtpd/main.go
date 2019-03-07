package main

import (
    "os"
    "io"
    "log"
    "fmt"
    "github.com/watsonserve/maild/smtpd/smtp"
)

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

    smtpServer := smtp.New("watsonserve.com", "127.0.0.1")

    fmt.Println("listen on port 10025")
    smtpServer.Listen(":10025")
}
