package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/helpers"
	"github.com/irellik/gblog/model"
	"github.com/irellik/gblog/service"
)

// 首页
func Index(c *gin.Context) {
	config := service.GetConfig()
	// 获取页码
	page := service.GetPage(c)
	offset := (page - 1) * config.Site.PageSize
	postList, total := model.GetPosts("", offset, config.Site.PageSize, false)

	pagination := helpers.MakePagination(c.Request, total, config.Site.PageSize).Paginate()
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
