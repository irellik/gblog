package model

import (
	"errors"
	sl "gblog/service/local"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id            int       `json:"id"`
	Nickname      string    `json:"nickname"`
	Username      string    `json:"username"`
	Secret        string    `json:"secret"`
	BindSecret    int       `json:"bind_secret"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	RememberToken string    `json:"remember_token"`
	CreatedAt     time.Time `json:"created_at"`
	LastLoginAt   time.Time `json:"last_login_at"`
	LastLoginIp   int       `json:"last_login_ip"`
}

var UserNotBindSecretError = errors.New("user not bind error")
var UserSecretNotMatchError = errors.New("user secret not match error")

func Auth(username string, password string, authenticatorCode string, clientIp int64) (User, error) {
	db := sl.MysqlClient
	rowSql := "select `id`, `secret`,`bind_secret`, `nickname`,`username`,`email`,`password`,`remember_token`,`created_at`,`last_login_at`,`last_login_ip` from users where `username` = ? limit 1"
	var user User
	var err error
	err = db.QueryRow(rowSql, username).Scan(&user.Id, &user.Secret, &user.BindSecret, &user.Nickname, &user.Username, &user.Email, &user.Password, &user.RememberToken, &user.CreatedAt, &user.LastLoginAt, &user.LastLoginIp)
	if err != nil {
		return user, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}
	// 先检查是否匹配
	otpConf := &sl.OTPConfig{
		Secret: user.Secret,
	}
	match, err := otpConf.Authenticate(authenticatorCode)
	if err != nil {
		return user, err
	}
	if !match {
		if user.BindSecret == 0 {
			return user, UserNotBindSecretError
		} else {
			return user, UserSecretNotMatchError
		}
	}

	// 更新绑定状态
	sql := "update users set `last_login_ip` = ? , `bind_secret` = 1 where `username` = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(clientIp, username)
	if err != nil {
		return user, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return user, err
	}
	return user, err
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
		return password, errPre
	}
	defer stmt.Close()
	res, errExc := stmt.Exec(1, "管理员", "admin", "", string(password_hash), "", time.Now(), time.Now(), 0)
	if errExc != nil {
		return password, errExc
	}
	_, errAff := res.RowsAffected()
	if errAff != nil {
		return password, errAff
	}
	return password, nil
}
