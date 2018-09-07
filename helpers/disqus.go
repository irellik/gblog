package helpers

import (
	"encoding/json"
	"github.com/irellik/gblog/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CommentCount struct {
	Url   string `json:"id"`
	Count int    `json:"comments"`
}

func UpdateCommentCount() {
	idList := model.GetAllPostId()
	chunkSize := 5
	jobList := make([][]string, 0)
	for i := 0; i < len(idList); i += chunkSize {
		end := i + chunkSize
		if end > len(idList) {
			end = len(idList)
		}
		jobList = append(jobList, idList[i:end])
	}
	var paramStr string
	for _, arr := range jobList {
		paramStr += strings.Join(arr, ".html&1=")
		commentCountUrl := "http://iwww.disqus.com/count-data.js?1=" + paramStr + ".html"
		response := HttpGet(commentCountUrl)
		countRegexp := regexp.MustCompile(`\"counts\":(\[\{\"id\":\".*\])\}\);`)
		params := countRegexp.FindStringSubmatch(response)
		countList := params[1]
		var cc []CommentCount
		json.Unmarshal([]byte(countList), &cc)
		for _, post := range cc {
			strSlice := strings.Split(post.Url, ".")
			id, _ := strconv.Atoi(strSlice[0])
			model.UpdateCommentCount(id, post.Count)
		}
		time.Sleep(time.Second * 5)
	}
}
