package lexical_test

import (
    "fmt"
    "github.com/watsonserve/imapd/lexical"
)

func TestLexical() {
    list := lexical.Parse("search UID 1:* (INTERNALDATE UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (date subject from content-type to cc bcc message-id in-reply-to references)])")
    for i := 0; i < len(list); i++ {
        el := list[i]
        fmt.Printf("%d %d %s|\n", el.Cnt, el.Type, el.Value)
    }
}
