package controller

import (
	"fmt"
	"gblog/model"
	sl "gblog/service/local"
	"gblog/utils"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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
	post, err := model.GetPost(postId, true)
	if err != nil {
		Throw404(c)
		return
	}
	post.ContentHtml = template.HTML(post.Content)
	config := sl.GetConfig()
	// 获取页码
	page := sl.GetPage(c)
	offset := (page - 1) * config.Site.PageSize
	_, total := model.GetPosts(offset, config.Site.PageSize, "1")
	settings := model.GetSettings()
	tagList := make([]string, 0)
	for _, keywordStruct := range post.Tags {
		tagList = append(tagList, keywordStruct.Name)
	}
	c.HTML(http.StatusOK, "index/article.html", gin.H{
		"settings":    settings,
		"categories":  model.GetCategories(),
		"totalTag":    len(model.GetTags()),
		"totalPost":   total,
		"friends":     model.GetFriends(),
		"post":        post,
		"title":       post.Title + " - " + settings["name"],
		"keywords":    strings.Join(tagList, ","),
		"description": post.Abstract,
	})
}

// 搜索
func Search(c *gin.Context) {
	config := sl.GetConfig()
	// 获取页码
	page := sl.GetPage(c)
	offset := (page - 1) * config.Site.PageSize
	postList, total := model.SearchPosts(c.Param("keyword"), offset, config.Site.PageSize)
	totalPage := total/config.Site.PageSize + 1
	pagination := utils.MakePagination(c.Request, total, config.Site.PageSize).Paginate()
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
