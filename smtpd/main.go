package main

import (
    "os"
    "io"
    "log"
    "fmt"
    "github.com/watsonserve/maild/smtpd/smtp"
    "github.com/watsonserve/maild/server"
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

    dispatcher := smtp.Init("watsonserve.com", "127.0.0.1")
    server := server.InitTCPServer()
    server.SetDispatcher(dispatcher)

    fmt.Println("listen on port 10025")
    server.Listen(":10025")
}