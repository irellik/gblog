package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/model"
	"net/http"
	"sort"
)

// 频道首页
func CategoryIndex(c *gin.Context) {
	category := c.Param("category")
	posts, _ := model.Archive(category, "category")
	tags := model.GetTags()
	friends := model.GetFriends()
	sortedKey := make([]string, 0)
	for key := range posts {
		sortedKey = append(sortedKey, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(sortedKey)))
	settings := model.GetSettings()
	c.HTML(http.StatusOK, "archive/index.html", gin.H{
		"settings":    settings,
		"categories":  model.GetCategories(),
		"totalTag":    len(tags),
		"friends":     friends,
		"posts":       posts,
		"sortedKey":   sortedKey,
		"message":     fmt.Sprintf("分类 %s 下的文章", category),
		"description": fmt.Sprintf("%s - %s", category, settings["name"]),
		"title":       fmt.Sprintf("分类 %s 下的文章 - %s", category, settings["name"]),
		"keywords":    fmt.Sprintf(category),
	})
}

// 分类页面
func CategoryAll(c *gin.Context) {
	posts, total := model.Archive("", "")
	tags := model.GetTags()
	friends := model.GetFriends()
	sortedKey := make([]string, 0)
	for key := range posts {
		sortedKey = append(sortedKey, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(sortedKey)))
	settings := model.GetSettings()
	c.HTML(http.StatusOK, "category/all.html", gin.H{
		"settings":    settings,
		"categories":  model.GetCategories(),
		"totalTag":    len(tags),
		"friends":     friends,
		"posts":       posts,
		"sortedKey":   sortedKey,
		"totalPost":   total,
		"description": fmt.Sprintf("分类 - %s", settings["name"]),
		"title":       fmt.Sprintf("分类 - %s", settings["name"]),
		"keywords":    fmt.Sprintf("分类"),
	})
}
