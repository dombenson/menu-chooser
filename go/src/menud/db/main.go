package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"menud/config"
	"menud/users"
)

type Connection interface {
	GetUser(int) (users.User, error)
	GetUserByEmailPassword(email, password string) (users.User, error)
}

type connection struct {
	baseConn    *sql.DB
	getUserStmt *sql.Stmt
	getUserEmailStmt *sql.Stmt
}

func GetConnection() Connection {
	obj := &connection{}
	obj.connect()
	return obj
}

func (this *connection) connect() {
	var err error
	this.baseConn, err = sql.Open("mysql", config.DBConnString())
	if err != nil {
		panic("Unable to connect to database")
	}
	this.getUserStmt, err = this.baseConn.Prepare("SELECT `id`,`name`,`email`,`pass` FROM `users` WHERE `userid` = ?")
	this.getUserEmailStmt, err = this.baseConn.Prepare("SELECT `id`,`name`,`email`,`pass` FROM `users` WHERE `email` = ?")
	if err != nil {
		panic("Unable to prepare user statement")
	}
}

func (this *connection) GetUser(id int) (_ users.User, err error) {
	var rows *sql.Rows
	rows, err = this.getUserStmt.Query(id)
	if err != nil {
		return
	}
	return users.MakeUser(rows)
}

func (this *connection) GetUserByEmailPassword(email, password string) (user users.User, err error) {
	var rows *sql.Rows
	rows, err = this.getUserEmailStmt.Query(email)
	if err != nil {
		return
	}
	user, err = users.MakeUser(rows)
	if err != nil {
		return
	}
	err = user.VerifyPassword(password)
	return
}


