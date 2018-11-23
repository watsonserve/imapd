package lexical

func endNumberConst(str string, length int) int {
	var i int

	if 0 == length {
		length = len(str)
	}

    /* 10进制 */
	if '0' != *str {
		for i = 0; i < length && (('0' <= str[i] && str[i] <= '9') || '.' == str[i]); i ++ {}
		return i
	}

    /* hex */
	if 'x' == str[1] {
		for i = 2; i < length && (('0' <= str[i] && str[i] <= '9') || ('A' <= str[i] && str[i] <= 'F')); i ++ {}
		return i
	}

    /* oct */
	for i = 0; i < length && ('0' <= str[i] && str[i] < '8'); i ++ {}
	return i
}

func endStringConst(endCh int, str string, length int) int {
	var i int

	if 0 == length {
		length = len(str)
	}

	for i = 0; i < length; i ++ {
		if(endCh == str[i]) {
			break
		}
		if('\\' == str[i]) {
			i++
		}
	}
	return i;
}

func endVariable(str string, length int) int {
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

func isOperationSignal(ch char) bool {
	lst := ",(){}[]+-*/%=&|~^!><;?:@#` \t\r\n"
    length := len(lst)

    for i := 0; i < length; i ++ {
        if (ch == lst[i]) {
			return true
		}
    }

	return false
}

/**
 * @return 成功 0，失败 -1
 */
func lexical(dest LexNodeList, p string, length int) int {  //词法分析
	var val node_t
	var i, length, j int

	j = 0;
	val.value = "";

	if(0 == length) {
		length = len(p)
	}

	for i = 0; i < length; i ++ {
		if '"' == p[i] || '\'' == p[i] {    //const string
			length = endStringConst(p[i], p + i + 1, 0)
			string ret(p + i + 1, length)
			i += length + 1
			val.value = ret
			val.type = CONST_STRING
		}
		elseif '0' <= p[i] && p[i] <= '9' {  //const number

			//printf("debug: %c ", p[i]);
			length = endNumberConst(p + i, 0);
			string ret(p + i, length);
			i += length - 1;
			val.value = ret;
			val.type = CONST_NUMBER;
		}
		else if ('A' <= p[i] && p[i] <= 'Z') || ('a' <= p[i] && p[i] <= 'z') || '_' == p[i] || '$' == p[i] || 0x80 <= p[i] {    //word
			length = endVariable(p + i, 0);
			//printf("%d ", length);
			string ret(p + i, length);
			i += length - 1;
			val.value = ret;
			val.type = VAR;
        }
        else if(isOperationSignal(p[i]))    //signal
        {
            val.value = p[i];
            val.type = OPERA;
        }
        else
        {
            return -1;
        }

		val.val = j ++;
		dest.push_back(val);
	}

	return 0;
}
