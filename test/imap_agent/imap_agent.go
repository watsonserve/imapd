package main

import (
    "os"
    "io"
    "log"
    "github.com/watsonserve/maild/imap_agent"
    "fmt"
)

type AgentConfig struct {}

func (this *AgentConfig) Auth(username string, password string) string {
    fmt.Printf("%s %s\n", username, password)
    return "session_id"
}

func (this *AgentConfig) Read() *imap_agent.UpResult {
    return &imap_agent.UpResult {
        Sess: "session_id",
        Result: "",
    }
}

func (this *AgentConfig) Send(sess string, spt *imap_agent.Mas) {
    fmt.Printf(`{"sess": "%s", "tag": "%s", "cmd": "%s", "params": "%s"}\n`, sess, spt.Tag, spt.Command, spt.Parames)
}

func main() {
    fp := os.Stderr
    log.SetOutput(io.Writer(fp))
    log.SetFlags(log.Ldate|log.Ltime|log.Lmicroseconds)

    imapServer := imap_agent.New("WS_IMAPD", &AgentConfig{})

    fmt.Println("listen on port 993")
    imapServer.TLSListen(":993", "etc/imap.crt", "etc/imap.key")
}
