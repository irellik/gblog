package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/model"
	"net/http"
	"sort"
)

func Tag(c *gin.Context) {
	tags := model.GetTags()
	friends := model.GetFriends()
	c.HTML(http.StatusOK, "tag/index.html", gin.H{
		"settings":   model.GetSettings(),
		"categories": model.GetCategories(),
		"tags":       tags,
		"friends":    friends,
	})
}

func TagIndex(c *gin.Context) {
	tag := c.Param("tag")
	posts, _ := model.Archive(tag, "")
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
		"message":    fmt.Sprintf("标签 %s 下的文章", tag),
	})
}
