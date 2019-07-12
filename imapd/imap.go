package imapd

import (
    "encoding/base64"
    "fmt"
    "github.com/watsonserve/maild/compile/lexical"
    "github.com/watsonserve/maild/lib"
    "strings"
)

func (this *ImapContext) Select(that *Imapd, script *Mas) {
    this.mailBox = ""
    if "" == script.Parames || "\"\"" == script.Parames {
        this.Send(fmt.Sprintf("%s BAD %s FAILURE: arguments invalid.\r\n", script.Tag, script.Command))
        return
    }
    // mailBox := script.Parames[0]
    // sql, _ := that.StmtMap["select"]
    // row := sql.QueryRow(mailBox, this.User)
    // err := row.Scan(mailBox)
    // if nil != err {
    //     this.Send(fmt.Sprintf("%s NO SELECT FAILURE: no such mailbox.\r\n", script.Tag))
    //     return
    // }
    this.mailBox = script.Parames
    sum := 0
    recentCnt := 0
    unseen := 0
    validCode := 1
    nextUID := 1

    // read write flag
    this.rw = 1
    rw := "READ-ONLY"
    if "SELECT" == script.Command {
        this.rw = 3
        rw = "READ-WRITE"
    }
    // response text
    resp := fmt.Sprintf("* FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)\r\n")
    resp += fmt.Sprintf("* %d EXISTS\r\n", sum)
    resp += fmt.Sprintf("* %d RECENT\r\n", recentCnt)
    resp += fmt.Sprintf("* OK [UIDVALIDITY %d] UIDs valid\r\n", validCode)
    resp += fmt.Sprintf("* OK [UIDNEXT %d] Predicted next UID\r\n", nextUID)
    // options 当当前邮箱中有被标记为未读的邮件时，则给出该行，该数字为第sn条消息，并非UID
    resp += fmt.Sprintf("* OK [UNSEEN %d] is the first unseen message\r\n", unseen)

    // A list of message flags that the client can change permanently.  If this is missing, the client should assume that all flags can be changed permanently.
    resp += fmt.Sprintf("* OK [PERMANENTFLAGS (\\Deleted \\Seen \\*)] Limited\r\n")
    this.Send(resp + fmt.Sprintf("%s OK [%s] %s completed.\r\n", script.Tag, rw, script.Command))
}

func (this *ImapContext) CLOSE(script *Mas) {
    this.mailBox = ""
    this.Send(fmt.Sprintf("%s OK CLOSE completed.\r\n", script.Tag))
}

/**
 * 创建目录
 * 与select命令的当前选中目录无关
 * @example: CREATE foo/bar/tar（bar不存在）则创建/foo/bar和/foo/bar/tar
 */
func (this *ImapContext) CREATE(script *Mas) {
    err := this.dal.Create(script.Parames)
    if nil != err {
        this.Send(fmt.Sprintf("%s NO %s.\r\n", script.Tag, err.Error()))
        return
    }
    this.Send(fmt.Sprintf("%s OK CREATE completed.\r\n", script.Tag))
}

func (this *ImapContext) DELETE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) RENAME(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) SUBSCRIBE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) UNSUBSCRIBE(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

// ref: {
//   "/":     不可以,
//   "通配符": 不可以,
//   "":      查找mailbox,
//   "name":  查找
// }

// mailbox: {
//   "/":     不可以,
//   "":      返回 "/" "",
//   "通配符": ref只能是邮箱名或"",
//   "name":  ref只能是""
// }

// ["foo",  ""] => [/, ""]
// ["",     ""] => [/, ""]
// ["",  "foo"] => /foo
// ["",    "*"] => [/foo, /bar, /tar, ...]
// ["foo", "*"] => [/foo, /foo/bar, /foo/bar/tar, /foo/...]
func (this *ImapContext) LIST(script *Mas) {
    // [ref: {"", "name"}, mailbox: {"", "通配符", "name"}]
    parames := lexical.Parse(script.Parames)
    if 2 != len(parames) {
        this.Send(fmt.Sprintf("%s NO LIST FAILURE: arguments invalid.\r\n", script.Tag))
        return
    }
    // 获取字符串值
    paths := []string{"", ""}
    for i := 0; i < 2; i++ {
        val := parames[i].Value
        if '"' == val[0] {
            val = val[1 : len(val) - 1]
        }
        paths[i] = val
    }
    path := paths[1]

    resp := ""
    for {
        // 缺失要查询的路径
        if "" == path {
            resp += "* LIST (\\Noselect \\HasChildren) \"/\" \"\"\r\n"
            break
        }
        if "" != paths[0] && "*" != path {
            break
        }
        // 整理为绝对路径
        path = strings.Join(paths, "/")
        // 查询数据
        mailboxes, err := this.dal.List(path)
        if nil != err {
            // TODO log
            break
        }
        // 整理目录列表
        mailboxList := *mailboxes
        for i := 0; i < len(mailboxList); i++ {
            mailbox := mailboxList[i]
            resp += fmt.Sprintf("* LIST (%s) \"/\" \"%s\"\r\n", strings.Join(mailbox.Attributes, " "), mailbox.Name)
        }
        break
    }
    this.Send(fmt.Sprintf("%s%s OK LIST completed.\r\n", resp, script.Tag))
}

func (this *ImapContext) LSUB(script *Mas) {
    this.Send(fmt.Sprintf("%s OK LSUB completed.\r\n", script.Tag))
}

func (this *ImapContext) STATUS(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) APPEND(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}

func (this *ImapContext) EXPUNGE(script *Mas) {
    // this.Send(fmt.Sprintf("* %d EXPUNGE\r\n", serial))
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}





func (this *ImapContext) ID(script *Mas) {
    this.Send(fmt.Sprintf("%s OK %s completed.\r\n", script.Tag, script.Command))
}
