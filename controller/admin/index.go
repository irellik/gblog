package admin

import (
	"gblog/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
	"gblog/service/local"
	"strconv"
	"encoding/gob"
	"gblog/helpers"
	"encoding/json"
	"encoding/base64"
	"time"
)

type loginForm struct {
	Username string `json:"username" xml:"username" form:"username" binding:"required,min=1"`
	Password string `json:"password" xml:"password" form:"password" binding:"required,min=8"`
	Remember string `json:"remember" xml:"remember" form:"remember" binding:"-"`
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
		helpers.Failed(c, http.StatusBadRequest, err.Error())
		return
	}
	var remember bool
	if string(c.PostForm("remember")) == "true" {
		remember = true
	} else {
		remember = false
	}
	userDB, err := model.Auth(c.PostForm("username"), c.PostForm("password"), helpers.InetAton(c.ClientIP()))
	if err != nil {
		helpers.Failed(c, http.StatusConflict, "incorrect username or password")
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
		MaxAge:config.Session.MaxAge,
		HttpOnly:false,
		Secure:false,
		Domain:config.Site.Domain,
		Path:"/",
	})
	session.Save()
	if remember {
		jsonUser,err := json.Marshal(user)
		if err != nil {
			helpers.Failed(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cookieByte,err := local.Encrypt([]byte(jsonUser))
		if err != nil {
			helpers.Failed(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.SetCookie("g_u", base64.StdEncoding.EncodeToString(cookieByte), config.Site.RememberDays * 24 * 60 * 60, "/", config.Site.Domain, false, false)
	}
	data := make(map[string]map[string]string, 0)
	data["user"] = user
	helpers.Success(c, data)
}
