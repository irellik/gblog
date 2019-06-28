package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Pagination struct {
	request *http.Request
	total   int
	size    int
	current int
}

// 获取日期
func TimeToDateStr(time time.Time) string {
	return time.Format("2006-01-02")
}

func TimeFormat(time time.Time, layout string) string {
	return time.Format(layout)
}

func TrimHtmlTag(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	// 去除<!--More-->
	return strings.TrimSpace(src)
}

// 生成分页
func (page *Pagination) Paginate() []map[string]string {
	total := page.total
	size := page.size
	current := page.current
	totalPage := total/size + 1
	var pagination = []map[string]string{}
	if current != 1 {
		pagination = page.addPagination(pagination, current-1, "«")
	}
	if totalPage <= 11 {
		for i := 1; i <= totalPage; i++ {
			pagination = page.addPagination(pagination, i, "")
		}
	} else {
		if current <= 6 {
			for i := 1; i <= 8; i++ {
				pagination = page.addPagination(pagination, i, "")
			}
			pagination = page.addPagination(pagination, "...", "")
			pagination = page.addPagination(pagination, totalPage-1, "")
			pagination = page.addPagination(pagination, totalPage, "")
		} else if current > 6 && current < totalPage-4 {
			pagination = page.addPagination(pagination, 1, "")
			pagination = page.addPagination(pagination, 2, "")
			pagination = page.addPagination(pagination, "...", "")
			for i := current - 3; i < current+3; i++ {
				pagination = page.addPagination(pagination, i, "")
			}
			pagination = page.addPagination(pagination, "...", "")
			pagination = page.addPagination(pagination, totalPage-1, "")
			pagination = page.addPagination(pagination, totalPage, "")
		} else {
			pagination = page.addPagination(pagination, 1, "")
			pagination = page.addPagination(pagination, 2, "")
			pagination = page.addPagination(pagination, "...", "")
			for i := totalPage - 8; i <= totalPage; i++ {
				pagination = page.addPagination(pagination, i, "")
			}
		}
	}
	if current != totalPage {
		pagination = page.addPagination(pagination, current+1, "»")
	}
	return pagination
}

//加法
func MathPlus(a int, b int) int {
	return a + b
}

// 减法
func MathReduce(a int, b int) int {
	return a - b
}

// Int 转 String
func IntToString(num int) string {
	return strconv.Itoa(num)
}

func MakePagination(request *http.Request, total int, size int) *Pagination {
	queryParams := request.URL.Query()
	pageParam := queryParams.Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	current, _ := strconv.Atoi(pageParam)
	return &Pagination{
		request: request,
		total:   total,
		size:    size,
		current: current,
	}
}

type Page interface {
}

// 生成页码信息
func (p *Pagination) addPagination(pagination []map[string]string, page Page, text string) []map[string]string {
	current := p.current
	isCurrent := "no"
	var urlString string
	if reflect.TypeOf(page).String() == "int" {
		if page == current {
			isCurrent = "yes"
		}
		u, _ := url.Parse(p.request.URL.String())
		q := u.Query()
		q.Set("page", fmt.Sprintf("%v", page))
		u.RawQuery = q.Encode()
		urlString = u.String()
	} else {
		urlString = ""
	}
	if text != "" {
		page = text
	}
	pagination = append(pagination, map[string]string{"page": fmt.Sprintf("%v", page), "url": urlString, "isCurrent": isCurrent})
	return pagination
}
