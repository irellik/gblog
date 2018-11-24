package admin

import (
	"fmt"
	"gblog/model"
	"github.com/gin-gonic/gin"
	"net/http"
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
	if err := c.ShouldBind(&lf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var remember bool
	if string(c.PostForm("remember")) == "1" {
		remember = true
	} else {
		remember = false
	}
	user, err := model.Auth(c.PostForm("username"), c.PostForm("password"), remember, c.ClientIP())
	fmt.Println(user, err)
}
