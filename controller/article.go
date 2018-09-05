package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"helpers"
	"html/template"
	"model"
	"net/http"
	"service"
	"strconv"
)

// 详情页
func Article(c *gin.Context) {
	// 文章ID合法性检测
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Throw404(c)
		return
	}
	// 获取文章内容
	post := model.GetPost(postId)
	if post.Id == 0 {
		Throw404(c)
		return
	}
	post.ContentHtml = template.HTML(post.Content)
	config := service.GetConfig()
	// 获取页码
	page := service.GetPage(c)
	offset := (page - 1) * config.Site.PageSize
	_, total := model.GetPosts("", offset, config.Site.PageSize, false)
	c.HTML(http.StatusOK, "index/article.html", gin.H{
		"settings":   model.GetSettings(),
		"categories": model.GetCategories(),
		"totalTag":   len(model.GetTags()),
		"totalPost":  total,
		"friends":    model.GetFriends(),
		"post":       post,
	})
}

// 搜索
func Search(c *gin.Context) {
	config := service.GetConfig()
	// 获取页码
	page := service.GetPage(c)
	offset := (page - 1) * config.Site.PageSize
	postList, total := model.GetPosts(c.Param("keyword"), offset, config.Site.PageSize, true)
	totalPage := total/config.Site.PageSize + 1
	pagination := helpers.MakePagination(c.Request, total, config.Site.PageSize).Paginate()
	tags := model.GetTags()
	friends := model.GetFriends()
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"postList":    postList,
		"settings":    model.GetSettings(),
		"categories":  model.GetCategories(),
		"pagination":  pagination,
		"currentPage": page,
		"totalPage":   totalPage,
		"totalPost":   total,
		"totalTag":    len(tags),
		"friends":     friends,
		"message":     fmt.Sprintf("包含关键字 %s 的文章", c.Param("keyword")),
	})
}
