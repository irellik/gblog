package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/irellik/gblog/model"
	"github.com/irellik/gblog/utils"
)

func CategoryList(c *gin.Context) {
	categories := model.GetCategories()
	data := make(map[string]interface{})
	data["categories"] = categories
	utils.Success(c, data)
}
