package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/model"
	sl "github.com/irellik/gblog/service/local"
	"github.com/irellik/gblog/utils"
	"net/http"
)

// 首页
func Index(c *gin.Context) {
	config := sl.GetConfig()
	// 获取页码
	page := sl.GetPage(c)
	offset := (page - 1) * config.Site.PageSize
	postList, total := model.GetPosts(offset, config.Site.PageSize, "1")

	pagination := utils.MakePagination(c.Request, total, config.Site.PageSize).Paginate()
	tags := model.GetTags()
	friends := model.GetFriends()
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"postList":    postList,
		"settings":    model.GetSettings(),
		"categories":  model.GetCategories(),
		"pagination":  pagination,
		"currentPage": page,
		"totalPost":   total,
		"totalTag":    len(tags),
		"friends":     friends,
	})
}
