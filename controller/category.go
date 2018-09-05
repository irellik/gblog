package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"model"
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
	c.HTML(http.StatusOK, "archive/index.html", gin.H{
		"settings":   model.GetSettings(),
		"categories": model.GetCategories(),
		"totalTag":   len(tags),
		"friends":    friends,
		"posts":      posts,
		"sortedKey":  sortedKey,
		"message":    fmt.Sprintf("分类 %s 下的文章", category),
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
	c.HTML(http.StatusOK, "category/all.html", gin.H{
		"settings":   model.GetSettings(),
		"categories": model.GetCategories(),
		"totalTag":   len(tags),
		"friends":    friends,
		"posts":      posts,
		"sortedKey":  sortedKey,
		"totalPost":  total,
	})
}
