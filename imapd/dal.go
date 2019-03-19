package imapd

import (
	"database/sql"
)

type MailBox struct {
    Name string
    Attributes []string
}

func InitMailBox(name string, hc, sa, nf, mk, sd bool) *MailBox {
    ret := &MailBox{Name: name}
    var hasChildren string

    hasChildren = "\\HasChildren"
    if !hc {
        hasChildren = "\\HasNoChildren"
    }
    ret.Attributes = append(ret.Attributes, hasChildren)

    if !sa {
        ret.Attributes = append(ret.Attributes, "\\Noselect")
    }

    if nf {
        ret.Attributes = append(ret.Attributes, "\\Noinferiors")
    }

    if mk {
        ret.Attributes = append(ret.Attributes, "\\Marked")
    }

    if sd {
        ret.Attributes = append(ret.Attributes, "\\" + name)
    }
    return ret
}




type DataAccessLayer struct {
    db      *sql.DB
    StmtMap map[string]*sql.Stmt
}

func NewDAL(db *sql.DB) *DataAccessLayer {
	this := &DataAccessLayer{}
	this.StmtMap = make(map[string]*sql.Stmt)

	this.prepare("list",   "SELECT box_name, has_children, selectable, noinferiors, marked, sys_defined FROM mail_boxes WHERE user_id=? AND box_name like ?")
	this.prepare("select", "SELECT box_id FROM mail_boxes WHERE user_id=? AND box_name=?")
	return this
}

func (this *DataAccessLayer) prepare(index string, query string) {
    stmt ,err := this.db.Prepare(query)
    if nil != err {
        panic(err)
    }
    this.StmtMap[index] = stmt
}

func (this *DataAccessLayer) query(index string, query ...interface{}) (*sql.Rows, error) {
    return this.StmtMap[index].Query(query...)
}

func (this *DataAccessLayer) exec(index string, query string) {
}

func (this *DataAccessLayer) Select(index string, query string) {
}

func (this *DataAccessLayer) List(query string) (*([]MailBox), error) {
    mailBoxes := make([]MailBox, 0)

    // rows, err := this.query("list", query)
    // if nil != err {
    //     return nil, err
    // }
    // defer rows.Close()

    // var name string
    // var has_children, selectable, noinferiors, marked, sys_defined bool
	// for rows.Next() {
    //     err = rows.Scan(
    //         &name,
    //         &has_children, &selectable, &noinferiors, &marked, &sys_defined,
    //     )
	// 	if err != nil {
	// 		return nil, err
    //     }
    //     mailbox := InitMailBox(name, has_children, selectable, noinferiors, marked, sys_defined)
	// 	mailBoxes = append(mailBoxes, *mailbox)
	// }
	// // Check for errors from iterating over rows.
	// if err = rows.Err(); err != nil {
	// 	return nil, err
    // }
    names := [] string {"Inbox", "Sent", "Notes", "Junk", "Drafts", "Trash"}
    for i := 0; i < len(names); i++ {
        mailbox := InitMailBox(names[i], false, true, false, false, true)
        mailBoxes = append(mailBoxes, *mailbox)
    }
    return &mailBoxes, nil
}
