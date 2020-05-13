package admin

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/model"
	"github.com/irellik/gblog/service/local"
	"github.com/irellik/gblog/utils"
	"github.com/skip2/go-qrcode"
	"net/http"
	"strconv"
	"time"
)

type loginForm struct {
	Username          string `json:"username" xml:"username" form:"username" binding:"required,min=1"`
	Password          string `json:"password" xml:"password" form:"password" binding:"required,min=8"`
	AuthenticatorCode string `json:"authenticator_code" xml:"authenticator_code" form:"authenticator_code" binding:"required,min=6,max=6"`
	Token             string `json:"token" xml:"token" form:"token" binding:"required,min=1"`
	Remember          string `json:"remember" xml:"remember" form:"remember" binding:"-"`
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func LoginView(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", gin.H{})
}

func Login(c *gin.Context) {
	var lf loginForm
	config := local.GetConfig()
	if err := c.ShouldBind(&lf); err != nil {
		utils.Failed(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	// 验证码检查
	httpClient := utils.MakeHttpClient()
	response, err := httpClient.HttpPost("https://www.google.com/recaptcha/api/siteverify", map[string]string{
		"secret":   config.Recaptcha.Secret,
		"response": lf.Token,
		"remoteip": c.ClientIP(),
	}, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		utils.Failed(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var rcMap map[string]interface{}
	err = json.Unmarshal(response, &rcMap)
	if err != nil {
		utils.Failed(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	success := rcMap["success"].(bool)
	if !success {
		utils.Failed(c, http.StatusForbidden, "403 Forbidden", nil)
		return
	}
	score := rcMap["score"].(float64)
	if score < config.Recaptcha.Score {
		utils.Failed(c, http.StatusForbidden, "403 Forbidden", nil)
		return
	}
	var remember bool
	if string(c.PostForm("remember")) == "true" {
		remember = true
	} else {
		remember = false
	}
	userDB, err := model.Auth(lf.Username, lf.Password, lf.AuthenticatorCode, utils.InetAton(c.ClientIP()))
	if err != nil {
		msg := err.Error()
		// 还未绑定
		if err == model.UserNotBindSecretError {
			otpConf := &local.OTPConfig{
				Secret: userDB.Secret,
			}
			siteName := model.GetSettings()["name"]
			optUri := otpConf.ProvisionURIWithIssuer(userDB.Username, siteName)
			png, err := qrcode.Encode(optUri, qrcode.Medium, 256)
			if err != nil {
				utils.Failed(c, http.StatusInternalServerError, msg, nil)
				return
			} else {
				utils.Failed(c, 461, msg, map[string]string{
					"qrImg": base64.StdEncoding.EncodeToString(png),
				})
				return
			}
		} else if err == model.UserSecretNotMatchError {
			utils.Failed(c, http.StatusConflict, msg, nil)
			return
		}
		utils.Failed(c, http.StatusConflict, msg, nil)
		return
	}
	// 认证成功之后设置session和对应cookie
	session := sessions.Default(c)
	user := make(map[string]string, 0)
	user["uid"] = strconv.Itoa(userDB.Id)
	user["username"] = userDB.Username
	user["nickname"] = userDB.Nickname
	user["email"] = userDB.Email
	user["expired"] = strconv.Itoa(int(time.Now().AddDate(0, 0, 30).Unix()))
	gob.Register(map[string]string{})
	session.Set("user_info", user)
	session.Options(sessions.Options{
		MaxAge:   config.Session.MaxAge,
		HttpOnly: false,
		Secure:   false,
		Domain:   config.Site.Domain,
		Path:     "/",
	})
	session.Save()
	if remember {
		jsonUser, err := json.Marshal(user)
		if err != nil {
			utils.Failed(c, http.StatusInternalServerError, "Internal Server Error", nil)
			return
		}
		cookieByte, err := local.Encrypt([]byte(jsonUser))
		if err != nil {
			utils.Failed(c, http.StatusInternalServerError, "Internal Server Error", nil)
			return
		}
		c.SetCookie("g_u", base64.StdEncoding.EncodeToString(cookieByte), config.Site.RememberDays*24*60*60, "/", config.Site.Domain, false, false)
	}
	data := make(map[string]map[string]string, 0)
	data["user"] = user
	utils.Success(c, data)
}
