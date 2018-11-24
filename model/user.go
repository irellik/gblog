package model

import (
	"fmt"
	sl "github.com/irellik/gblog/service/local"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id            int       `json:"id"`
	Nickname      string    `json:"nickname"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	RememberToken string    `json:"remember_token"`
	CreatedAt     time.Time `json:"created_at"`
	LastLoginAt   time.Time `json:"last_login_at"`
	LastLoginIp   int       `json:"last_login_ip"`
}

func Auth(username string, password string, remember bool, clientIp string) (User, bool) {
	db := sl.MysqlClient
	rowSql := "select * from users where `username` = ? limit 1"
	var user User
	var err error
	var rememberToken []byte
	err = db.QueryRow(rowSql, username).Scan(&user)
	if err != nil {
		return user, false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, false
	} else {
		// 密码验证成功
		if remember {
			sourceToken := []byte(fmt.Sprintf("%s-%s", sl.RandStr(16), string(time.Now().AddDate(0, 0, 30).Unix())))
			rememberToken, err = sl.Encrypt(sourceToken)
			if err != nil {
				return user, false
			}
		} else {
			rememberToken = []byte("")
		}
	}
	sql := "update users set `last_login_ip` = ?,`remember_token` = ? where `username` = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return user, false
	}
	res, err := stmt.Exec(clientIp, rememberToken, username)
	if err != nil {
		return user, false
	}
	defer stmt.Close()
	_, err = res.RowsAffected()
	if err != nil {
		return user, false
	}
	return user, true
}
