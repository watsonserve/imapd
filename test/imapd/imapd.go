package main

import (
    "database/sql"
    "os"
    "io"
    "log"
    "github.com/watsonserve/maild/lib"
    "github.com/watsonserve/maild/imapd"
    "fmt"
)

type Author struct {
    lib.Author
}

func (this *Author) Auth(username string, password string) string {
    fmt.Printf("%s %s\n", username, password)
    return "null"
}

func ConnPg(config map[string]string) *sql.DB {

    pgurl := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        config["DBUser"],
        config["DBPasswd"],
        config["DBHost"],
        config["DBPort"],
        config["DBName"],
    )
    db, err := sql.Open("postgres", pgurl)
    if nil != err {
        panic(err)
    }
    return db
}

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

    // db := ConnPg()
    // dbc := Imapd.NewDAL(db)
    imapServer := imapd.New(nil, "imap.watsonserve.com", "127.0.0.1")
    imapServer.Author = &Author{}

    fmt.Println("listen on port 993")
    imapServer.TLSListen(":993", "etc/imap.crt", "etc/imap.key")
}
