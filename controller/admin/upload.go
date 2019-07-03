package admin

import (
	"fmt"
	"gblog/service/local"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func UploadImage(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File

	date := time.Now().Format("20060102")
	config := local.GetConfig()
	data := make([]string, 0)
	for filename, file := range files {
		path := fmt.Sprintf("%s/%s", config.Site.UploadPath, date)
		dst := fmt.Sprintf("%s/%s", path, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}
		// Upload the file to specific dst.
		err := c.SaveUploadedFile(file[0], dst)
		if err == nil {
			data = append(data, fmt.Sprintf("https://%s/uploads/%s/%s", config.Site.Domain, date, filename))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"errno": 0,
		"data":  data,
	})
}
