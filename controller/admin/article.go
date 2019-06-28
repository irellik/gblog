package admin

import (
	"gblog/model"
	sl "gblog/service/local"
	"gblog/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// 创建文章
func ArticleStore(c *gin.Context) {
	session := sessions.Default(c)
	user_info := session.Get("user_info")
	u_info := user_info.(map[string]string)
	var article_form model.ArticleForm
	if err := c.ShouldBind(&article_form); err != nil {
		utils.Failed(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	// 时间格式化
	if article_form.PublishedAt == "" {
		article_form.PublishedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	// 插入
	uid, _ := strconv.Atoi(u_info["uid"])
	articleId, err := model.InsertPost(article_form, int64(uid))
	if err != nil {
		utils.Failed(c, http.StatusInternalServerError, "保存失败", nil)
		return
	}
	response := map[string]int64{
		"article_id": articleId,
	}
	utils.Success(c, response)
}

// 文章列表
func ArticleList(c *gin.Context) {
	config := sl.GetConfig()
	page := sl.GetPage(c)
	offset := (page - 1) * config.Site.PageSize
	pageSize := 20
	status := c.Query("status")
	// 是否需要搜索
	keyword := c.Query("s")
	var posts []model.Post
	var total int
	if keyword == "" {
		posts, total = model.GetPosts(offset, pageSize, status)
	} else {
		posts, total = model.SearchPosts(keyword, offset, pageSize)
	}
	responsePosts := make([]map[string]interface{}, 0)
	for _, post := range posts {
		responsePosts = append(responsePosts, map[string]interface{}{
			"id":         post.Id,
			"title":      post.Title,
			"created_at": post.CreatedAt.Format("2006-01-02 15:04:05"),
			"abstract":   post.Abstract,
		})
	}
	data := map[string]interface{}{
		"posts": responsePosts,
		"total": total,
	}
	utils.Success(c, data)
}

// 删除文章
func ArticleDelete(c *gin.Context) {
	idList := c.QueryArray("id_list[]")
	res := false
	if len(idList) > 0 {
		res = model.DeletePost(idList)
	}
	if res {
		utils.Success(c, nil)
	} else {
		utils.Failed(c, http.StatusInternalServerError, "failed", nil)
	}
}

/**
详情
*/
func ArticleDetail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.Failed(c, http.StatusNotFound, "Not Found", nil)
		return
	}
	post, err := model.GetPost(idInt, false)
	if err != nil {
		utils.Failed(c, http.StatusNotFound, "Not Found", nil)
		return
	}
	utils.Success(c, map[string]interface{}{
		"post": map[string]interface{}{
			"id":           post.Id,
			"title":        post.Title,
			"content":      post.Content,
			"status":       post.Status,
			"category":     post.CatId,
			"published_at": post.PublishedAt.Unix(),
			"abstract":     post.Abstract,
			"tags":         post.Tags,
		},
	})
}

func ArticleUpdate(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.Failed(c, http.StatusNotFound, "Not Found", nil)
		return
	}
	_, err = model.GetPost(idInt, false)
	if err != nil {
		utils.Failed(c, http.StatusNotFound, "Not Found", nil)
		return
	}

	var article_form model.ArticleForm
	if err := c.ShouldBind(&article_form); err != nil {
		utils.Failed(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	// 时间格式化
	if article_form.PublishedAt == "" {
		article_form.PublishedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	// 更新
	err = model.UpdatePost(idInt, article_form)
	if err != nil {
		utils.Failed(c, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}
	utils.Success(c, nil)
}
