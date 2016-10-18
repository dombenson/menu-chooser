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
	id  int
	name  string
	email string
	pass  string
}

func MakeUser(rows *sql.Rows) (user User, err error) {
	retUser := &user{}
	user = retUser
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
	return bcrypt.CompareHashAndPassword(this.pass, []byte(pass))
}
