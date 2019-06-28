package admin

import (
	"gblog/model"
	"gblog/utils"
	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context) {
	categories := model.GetCategories()
	data := make(map[string]interface{})
	data["categories"] = categories
	utils.Success(c, data)
}
