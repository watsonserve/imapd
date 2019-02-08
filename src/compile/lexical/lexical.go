package lexical

/*
#cgo CFLAGS: -O3
#include "lexical.h"
*/
import "C"
import (
    "container/list"
)

const (
    CONST_STRING = 1
    CONST_NUMBER = 2
    VAR = 4
    OPERA = 8
)

type Lexical_t struct {
    Value string
    Type  int
    Cnt   int
}

func getType(ch byte) int {
    // signal
    if 1 == int(C.isOperationSignal(C.char(ch))) {
        return OPERA
    }

    // const string
    if '"' == ch || '\'' == ch {
        return CONST_STRING
    }

    // const number
    if '0' <= ch && ch <= '9' {
        return CONST_NUMBER
    }

    // word
    if '_' == ch || '$' == ch || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || 0x80 <= ch {
        return VAR
    }

    return 0
}

/**
 * @ 词法分析
 * @return 成功 0，失败 -1
 */
func Parse(str string) *list.List {
    var val Lexical_t
    dest := list.New()
    count := 0
    length := len(str)
    val.Value = ""
    cstr := C.CString(str)

    for i := 0; i < length; i++ {
        if ' ' == str[i] || '\t' == str[i] || '\r' == str[i] {
            continue
        }

        val.Type = getType(str[i])

        if OPERA == val.Type {
            val.Value = string(str[i])
        } else {
            cstring := C.movPointer(cstr, C.int(i))
            limit := C.int(length - i)
            increment := 0
            wordSize := 0

            switch val.Type {
            case CONST_STRING:
                increment = int(C.endStringConst(C.char(str[i]), cstring, limit))
                wordSize = 1 + increment
            case CONST_NUMBER:
                wordSize = int(C.endNumberConst(cstring, limit))
                increment = wordSize - 1
            case VAR:
                wordSize = int(C.endVariable(cstring, limit))
                increment = wordSize - 1
            default:
                return nil
            }

            end := i + wordSize
            if length < end {
                end -= 1
            }
            val.Value = str[i:end]
            i += increment
        }

        val.Cnt = count
        count++
        dest.PushBack(val)
    }

    return dest
}

