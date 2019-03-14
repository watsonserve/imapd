package imapd

import (
    "strconv"
    "github.com/watsonserve/maild/compile/lexical"
)

// mail access structor
type Mas struct {
    Tag string
    Command string
    Parames []string
}

func initMas(raw []lexical.Lexical_t) *Mas {
    length := len(raw)
    ret := &Mas {
        Tag: raw[0].Value,
        Command: "",
        Parames: nil,
    }
    if 1 < length {
        ret.Command = raw[1].Value

        if 2 < length {
            length -= 2
            for i:= 2; i < length; i++ {
                ret.Parames = append(ret.Parames, raw[i].Value)
            }
        }
    }

    return ret
}
