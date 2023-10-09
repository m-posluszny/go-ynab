package auth

import (
	"fmt"

	"github.com/m-posluszny/go-ynab/src/db"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
type RegisterForm struct {
	LoginForm
	RePassword string `form:"repassword" binding:"required"`
}

type Credentials struct {
	Uid          string
	Username     string
	PasswordHash []byte `db:"password_hash"`
}

func (form LoginForm) HashedPassword() []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hashedPassword
}

func (form LoginForm) dbView() Credentials {
	return Credentials{Username: form.Username, PasswordHash: form.HashedPassword()}
}

func GetUserFromUid(dbx *db.DBRead, uid string) (*Credentials, error) {
	creds := Credentials{}
	err := dbx.Get(&creds,
		`SELECT uid, username, password_hash FROM credentials WHERE uid=$1;`,
		uid)
	return &creds, err
}

func GetUserFromName(dbx *db.DBRead, username string) (*Credentials, error) {
	creds := Credentials{}
	err := dbx.Get(&creds,
		`SELECT uid, username, password_hash FROM credentials WHERE username=$1;`,
		username)
	return &creds, err
}
func CreateUser(dbx *db.DBWrite, form RegisterForm) (*Credentials, error) {
	newUser := form.dbView()
	response, err := dbx.NamedExec(
		`INSERT INTO credentials (username, uid, password_hash) VALUES (:username, gen_random_uuid(), :password_hash);`,
		newUser)
	fmt.Println(response, err)
	if err != nil {
		return nil, err
	}
	return GetUserFromName(dbx, newUser.Username)
}

func MatchPassword(dbx *db.DBRead, form LoginForm) bool {
	creds, err := GetUserFromName(dbx, form.Username)
	if err != nil {
		panic(err)
	}
	return bcrypt.CompareHashAndPassword(creds.PasswordHash, []byte(form.Password)) == nil
}