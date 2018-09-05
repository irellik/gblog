package main

import (
	"controller"
	"github.com/gin-gonic/gin"
	"helpers"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"service"
	"strings"
)

func main() {
	// 加载全局配置
	service.LoadConfig()
	globalConfig := service.GetConfig()
	// 连接MySQL
	service.MysqlInit()
	//var LayoutView = filepath.Join(getCurrentPath(), "./views/layout.html")
	// Disable Console Color
	// gin.DisableConsoleColor()

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

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(globalConfig.Site.Address)
	// router.Run(":3000") for a hard coded port
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
	engine.LoadHTMLGlob(filepath.Join(getCurrentPath(), "./views/**/*"))
}

func getCurrentPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("%s", err)
	}
	return strings.Replace(path, "\\", "/", -1)
}
