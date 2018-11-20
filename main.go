package main

import (
	adminController "gblog/controller/admin"
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/controller"
	"github.com/irellik/gblog/helpers"
	sl "github.com/irellik/gblog/service/local"
	st "github.com/irellik/gblog/service/third"
	"html/template"
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

	admin := router.Group("/admin")
	admin.Use(authMiddleware)
	{
		admin.GET("/", adminController.Index)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(globalConfig.Site.Address)
	// router.Run(":3000") for a hard coded port
}

func authMiddleware(c *gin.Context) {

	// Pass on to the next-in-chain

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
