package main

import (
	"flag"
	"fmt"
	"gblog/controller"
	"gblog/controller/admin"
	admin2 "gblog/controller/admin"
	"gblog/middleware"
	"gblog/model"
	sl "gblog/service/local"
	"gblog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"path/filepath"
)

var command string

func init() {
	flag.StringVar(&command, "command", "", "some command")
}

func main() {
	// 加载全局配置
	globalConfig := sl.LoadConfig()
	// 连接MySQL
	sl.MysqlInit()

	flag.Parse()
	if command == "set:admin" {
		password, err := model.SetAdmin()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("管理员已经设置，账号是admin,密码是%s\n", password)
		}
		os.Exit(0)
	}

	//var LayoutView = filepath.Join(getCurrentPath(), "./views/layout.html")
	// Disable Console Color
	// gin.DisableConsoleColor()

	// 更新评论
	//go st.UpdateCommentCount()
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	//// session
	store := cookie.NewStore([]byte(globalConfig.AppKey))
	router.Use(sessions.Sessions("g_session", store))

	// 静态文件
	router.Static("/static", "./static")
	router.Static("/admin", "./static/html/admin")
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

	router_admin := router.Group("/api")

	router_admin.Use(middleware.AuthMiddleware)

	{
		router_admin.GET("/", admin.Index)
		router_admin.GET("/category", admin.CategoryList)
		router_admin.POST("/article", admin.ArticleStore)
		router_admin.DELETE("/article", admin.ArticleDelete)
		router_admin.GET("/article", admin.ArticleList)
		router_admin.GET("/article/:id", admin.ArticleDetail)
		router_admin.PUT("/article/:id", admin.ArticleUpdate)
	}
	router.GET("/login", admin2.LoginView)
	router.POST("/api/login", admin2.Login)
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(globalConfig.Site.Address)
}

func loadView(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"toDate":      utils.TimeToDateStr,
		"timeFormat":  utils.TimeFormat,
		"mathPlus":    utils.MathPlus,
		"mathReduce":  utils.MathReduce,
		"intToString": utils.IntToString,
	}
	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob(filepath.Join(sl.GetCurrentPath(), "./views/**/*"))
}
