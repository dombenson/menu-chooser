package users

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Name() string
	Email() string
	VerifyPassword(string) error
}

type user struct {
	id    int
	name  string
	email string
	pass  string
}

const GetUserSQL = "SELECT `userid`,`name`,`email`,`pass` FROM `users` WHERE `userid` = ?"
const GetUserByEmailSQL = "SELECT `userid`,`name`,`email`,`pass` FROM `users` WHERE `email` = ?"

func MakeUser(rows *sql.Rows) (newUser User, err error) {
	retUser := &user{}
	newUser = retUser
	err = rows.Scan(&retUser.id, &retUser.name, &retUser.email, &retUser.pass)
	return
}

func (this *user) Name() string {
	return this.name
}
func (this *user) Email() string {
	return this.email
}
func (this *user) VerifyPassword(pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(this.pass), []byte(pass))
}
