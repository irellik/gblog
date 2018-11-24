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
	defer stmt.Close()
	res, err := stmt.Exec(clientIp, rememberToken, username)
	if err != nil {
		return user, false
	}
	_, err = res.RowsAffected()
	if err != nil {
		return user, false
	}
	return user, true
}

func SetAdmin() (string, error) {
	password := sl.RandStr(8)
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return password, err
	}
	sql := "replace into `users` (`id`,`nickname`,`username`,`email`,`password`,`remember_token`,`created_at`,`last_login_at`,`last_login_ip`) values (?,?,?,?,?,?,?,?,?)"
	db := sl.MysqlClient
	stmt, errPre := db.Prepare(sql)
	if errPre != nil {
		return password, err
	}
	defer stmt.Close()
	res, errExc := stmt.Exec(1, "管理员", "admin", "", password_hash, "", time.Now(), "", 0)
	if errExc != nil {
		return password, err
	}
	_, errAff := res.RowsAffected()
	if errAff != nil {
		return password, err
	}
	return password, nil
}
