package imapd

import (
	"database/sql"
)

type DataAccessLayer struct {
    db      *sql.DB
    StmtMap map[string]*sql.Stmt
}

func NewDAL(db *sql.DB) *DataAccessLayer {
	this := &DataAccessLayer{}
	this.StmtMap = make(map[string]*sql.Stmt)

	this.prepare("select", "select box_id where name=? and user_id=?")
	return this
}

func (this *DataAccessLayer) prepare(index string, query string) {
    stmt ,err := this.db.Prepare(query)
    if nil != err {
        panic(err)
    }
    this.StmtMap[index] = stmt
}

func (this *DataAccessLayer) Query(index string, query string) {
}

func (this *DataAccessLayer) Exec(index string, query string) {
}
