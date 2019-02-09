package imap

import (
    "strconv"
    "github.com/watsonserve/maild/compile/lexical"
)

// mail access structor
type Mas struct {
    Count int
    Command string
    Parames []string
}

func initMas(raw []lexical.Lexical_t) *Mas {
    length := len(raw)
    count, err := strconv.ParseInt(raw[0].Value, 10, 0)

    ret := &Mas {
        Count: 0,
        Command: "",
        Parames: nil,
    }

    if nil == err {
        ret.Count = int(count)
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
