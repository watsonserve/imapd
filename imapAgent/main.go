package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/watsonserve/goutils"
	"github.com/watsonserve/imapd/lib"
)

type AgentConfig struct{}

func (this *AgentConfig) Auth(username string, password string) string {
	fmt.Printf("%s %s\n", username, password)
	return "session_id"
}

func (this *AgentConfig) Read() *UpResult {
	return &UpResult{
		Sess:   "session_id",
		Result: "",
	}
}

func (this *AgentConfig) Send(sess string, spt *lib.Mas) {
	fmt.Printf(`{"sess": "%s", "tag": "%s", "cmd": "%s", "params": "%s"}\n`, sess, spt.Tag, spt.Command, spt.Parames)
}

func main() {
	fp := os.Stderr
	log.SetOutput(io.Writer(fp))
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	ln, err := goutils.TLSSocket(":993", "etc/imap.crt", "etc/imap.key")
	if nil != err {
		log.Println(err)
		return
	}

	fmt.Println("listen on port 993")
	Service("WS_IMAPD", ln, &AgentConfig{})
}
