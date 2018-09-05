package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"model"
	"net/http"
	"sort"
)

func Archive(c *gin.Context) {
	posts, total := model.Archive("", "")
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
		"message":    fmt.Sprintf("好! 目前共计 %d 篇日志。 继续努力。", total),
	})
}
