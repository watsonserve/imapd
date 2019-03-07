package main

import (
    "os"
    "io"
    "log"
    "github.com/watsonserve/maild/imapd/imap"
    "github.com/watsonserve/maild/server"
    "fmt"
)

func main() {

    fp := os.Stderr
    /*
    fp, err := os.OpenFile("/var/log/mail_auth.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
    if nil != err {
        log.Fatal(err)
        return
    }
    */
    log.SetOutput(io.Writer(fp))
    log.SetFlags(log.Ldate|log.Ltime|log.Lmicroseconds)

    imapd := imap.New("imap.watsonserve.com", "127.0.0.1")

    fmt.Println("listen on port 10143")
    imapd.Listen(":10143")
}
