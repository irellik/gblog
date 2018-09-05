package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Throw404(c *gin.Context) {
	c.String(http.StatusNotFound, "%s", "404 page not found")
	c.Abort()
}
