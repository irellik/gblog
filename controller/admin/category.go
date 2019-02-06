package admin

import (
	"gblog/model"
	"gblog/helpers"
	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context){
	categories := model.GetCategories()
	data := make(map[string]interface{})
	data["categories"] = categories
	helpers.Success(c, data)
}
