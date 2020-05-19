// 词法分析
package lexical

/*
#cgo CFLAGS: -O3
#include "lexical.h"
*/
import "C"

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

// 解析入口
func Parse(str string) []Lexical_t {
    var val Lexical_t
    dest := make([]Lexical_t, 0)
    count := 0
    length := len(str)
    val.Value = ""
    cstr := C.CString(str)

    for i := 0; i < length; i++ {
        if ' ' == str[i] || '\t' == str[i] || '\r' == str[i] {
            continue
        }

        val.Type = getType(str[i])

        if 0 != val.Type {
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
        } else {
            // signal
            head := str[i : i + 2]
            operaLen := int(C.endOperationSignal(C.CString(head), 2))
            if 0 < operaLen {
                val.Type = OPERA
                end := i + operaLen
                val.Value = string(str[i:end])
                i = end - 1
            }
        }

        val.Cnt = count
        count++
        dest = append(dest, val)
    }

    return dest
}
