package main

import (
    "os"
    "io"
    "log"
    "fmt"
    "github.com/watsonserve/maild/lib"
    "github.com/watsonserve/maild/smtpd"
)

type Author struct {
    lib.Author
}

func (this *Author) Auth(username string, password string) string {
    fmt.Printf("%s %s\n", username, password)
    return "null"
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

    smtpServer := smtpd.New("watsonserve.com", "127.0.0.1")
    smtpServer.Author = &Author{}

    fmt.Println("listen on port 10025")
    smtpServer.Listen(":10025")
}
