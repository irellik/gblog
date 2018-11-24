package main

import (
	admin2 "gblog/controller/admin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/controller"
	"github.com/irellik/gblog/controller/admin"
	"github.com/irellik/gblog/helpers"
	sl "github.com/irellik/gblog/service/local"
	st "github.com/irellik/gblog/service/third"
	"html/template"
	"net/http"
	"path/filepath"
)

func main() {
	// 加载全局配置
	sl.LoadConfig()
	globalConfig := sl.GetConfig()
	// 连接MySQL
	sl.MysqlInit()
	//var LayoutView = filepath.Join(getCurrentPath(), "./views/layout.html")
	// Disable Console Color
	// gin.DisableConsoleColor()

	// 更新评论
	go st.UpdateCommentCount()
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	// 静态文件
	router.Static("/static", "./static")
	// 定义模板位置
	loadView(router)
	router.GET("/", controller.Index)
	router.GET("/category/:category", controller.CategoryIndex)
	router.GET("/category", controller.CategoryAll)
	router.GET("/archive", controller.Archive)
	router.GET("/tag", controller.Tag)
	router.GET("/tag/:tag", controller.TagIndex)
	router.GET("/post/:id", controller.Article)
	router.GET("/search/:keyword", controller.Search)

	router_admin := router.Group("/admin")
	//// session
	store := cookie.NewStore([]byte("secret"))
	router_admin.Use(sessions.Sessions("my_session", store))
	router_admin.Use(authMiddleware)

	{
		router_admin.GET("/", admin.Index)
	}
	router.GET("/login", admin2.LoginView)
	router.POST("/login", admin2.Login)
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(globalConfig.Site.Address)
}

func authMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user_info := session.Get("user_info")
	if user_info == nil {
		c.Redirect(http.StatusFound, "/login")
	}
	c.Next()
}

func loadView(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"toDate":      helpers.TimeToDateStr,
		"timeFormat":  helpers.TimeFormat,
		"mathPlus":    helpers.MathPlus,
		"mathReduce":  helpers.MathReduce,
		"intToString": helpers.IntToString,
	}
	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob(filepath.Join(sl.GetCurrentPath(), "./views/**/*"))
}
