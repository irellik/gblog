package middleware

import (
	"encoding/json"
	"strconv"
	"time"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	sl "gblog/service/local"
	"gblog/helpers"
	"net/http"
	"encoding/gob"
)


func AuthMiddleware(c *gin.Context) {
	gob.Register(map[string]string{})
	globalConfig := sl.GetConfig()
	session := sessions.Default(c)
	user_info := session.Get("user_info")
	if user_info == nil {
		user := make(map[string]string, 0)
		// get cookie
		cookieBase64Encode,err := c.Cookie("g_u")
		if err != nil {
			helpers.Failed(c, http.StatusUnauthorized,err.Error())
			return
		}
		// base64 decode
		cookieBase64Decode,err := base64.StdEncoding.DecodeString(cookieBase64Encode)
		if err != nil {
			helpers.Failed(c, http.StatusUnauthorized,err.Error())
			return
		}
		// decrypt cookie
		cookieByte,err := sl.Decrypt(cookieBase64Decode)
		if err != nil {
			helpers.Failed(c, http.StatusUnauthorized,err.Error())
			return
		}
		// json to map
		err = json.Unmarshal(cookieByte, &user)
		if err != nil{
			helpers.Failed(c, http.StatusUnauthorized,err.Error())
			return
		}
		// if cookie is valid,set session
		expired,_ := strconv.Atoi(user["expired"])
		if int(time.Now().Unix()) < expired {
			session.Set("user_info", user)
			session.Options(sessions.Options{
				MaxAge:globalConfig.Session.MaxAge,
			})
			session.Save()
		}else{
			helpers.Failed(c, http.StatusUnauthorized,"unauthorized")
			return
		}
	}
	c.Next()
}
