package compile

import (
	"container/list"
)

const (
    CONST_STRING = 1
	CONST_NUMBER = 2
	VAR = 4
    OPERA = 8
)

type Lexical_t {
	Value string
	Type  int
	Cnt   int
}

func endNumberConst(str []byte) int {
	var i int
	length := len(str)

    /* 10进制 */
	if '0' != str[0] {
		for i = 0; i < length && (('0' <= str[i] && str[i] <= '9') || '.' == str[i]); i++ {}
		return i
	}

    /* hex */
	if 'x' == str[1] || 'X' == str[1] {
		for i = 2; i < length && (('0' <= str[i] && str[i] <= '9') || ('A' <= str[i] && str[i] <= 'F')); i++ {}
		return i
	}

    /* oct */
	for i = 0; i < length && ('0' <= str[i] && str[i] < '8'); i++ {}
	return i
}

func endStringConst(str string, length int, endCh byte, offset int) int {
	var i int

	for i = offset; i < length; i++ {
		if(endCh == str[i]) {
			break
		}
		if('\\' == str[i]) {
			i++
		}
	}
	return i
}

func endVariable(str []byte, length int) int {
	var i int

	if 0 == length {
		length = len(str)
	}

	for i = 0;
		(
			('A' <= str[i] && str[i] <= 'Z') ||
			('a' <= str[i] && str[i] <= 'z') ||
			('0' <= str[i] && str[i] <= '9') ||
			'$' == str[i] ||
			'_' == str[i] ||
			0x80 <= str[i]);
		i ++ {}

	return i;
}

func isOperationSignal(ch byte) bool {
	lst := ",(){}[]+-*/%=&|~^!><;?:@#` \t\r\n"
    length := len(lst)

    for i := 0; i < length; i++ {
        if (ch == lst[i]) {
			return true
		}
    }

	return false
}

/**
 * @ 词法分析
 * @return 成功 0，失败 -1
 */
func Lexical(str string) *list.List {
	var val Lexical_t
	dest := list.New()
	count := 0
	end := 0
	length := len(str)
	val.Value = ""

	for i := 0; i < length; i++ {
		switch str[i] {
			case '"':    //const string
				fallthrough
			case '\'':    //const string
				end = endStringConst(str, length, str[i], i + 1)
				string ret(p + i + 1, length)
				i += length + 1
				val.Value = str[i:end]
				val.Type = CONST_STRING
			case '0' <= p[i] && p[i] <= '9':  //const number
				length = endNumberConst(p[i:])
				string ret(p + i, length)
				i += length - 1
				val.Value = ret
				val.Type = CONST_NUMBER
			case 'A' <= p[i] && p[i] <= 'Z':
				fallthrough
			case 'a' <= p[i] && p[i] <= 'z':
				fallthrough
			case '_':
				fallthrough
			case '$':
				fallthrough
			case 0x80 <= p[i]:    //word
				length = endVariable(p + i, 0)
				string ret(p + i, length)
				i += length - 1
				val.Value = ret
				val.Type = VAR
			case isOperationSignal(p[i]):    //signal
				val.Value = p[i]
				val.Type = OPERA
			default:
				return -1
		}

		val.Cnt = count++
		dest.PushBack(val)
	}

	return dest
}
