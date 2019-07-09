package imapd

import (
    "strings"
)

// mail access structor
type Mas struct {
    Tag string
    Command string
    Parames string
}

func initMas(msg string) *Mas {
    raw := strings.SplitN(msg, " ", 3)
    length := len(raw)
    if length < 2 {
        return nil
    }
    ret := &Mas {
        Tag: raw[0],
        Command: raw[1],
        Parames: "",
    }
    if 2 < length {
        ret.Parames = raw[2]
    }

    return ret
}
