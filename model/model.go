package model

import (
	"database/sql"
	"fmt"
	sl "gblog/service/local"
	"log"
	"strconv"
	"strings"
)


type Tag struct {
	Id        int    `json:"id" form:"id"`
	Name      string `json:"name" form:"name"`
	PostCount int    `json:"post_count"`
	WeightCss int    `json:"weight_css"`
}

// 站点配置
type Setting struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

type Friends struct {
	Id   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Url  string `json:"url" form:"url"`
}

var postTable string = "posts"
var categoryTable string = "categories"
var settingTable string = "settings"
var tagTable string = "tags"
var friendsTable string = "friends"
var postTagTable string = "post_tag"


func Archive(tag string, t string) (map[string][]Post, int) {
	db := sl.MysqlClient
	var rows *sql.Rows
	var err error
	if tag != "" {
		var rowsSql string
		if t == "category" {
			rowsSql = fmt.Sprintf("select p.id,p.title,p.published_at from %s as p left join %s as c on c.id = p.cat_id where c.en_name = ?", postTable, categoryTable)
		} else {
			rowsSql = fmt.Sprintf("select p.id,p.title,p.published_at from %s as p left join %s as pt on pt.post_id = p.id left join %s as t on t.id = pt.tag_id where t.name = ?", postTable, postTagTable, tagTable)
		}
		rows, err = db.Query(rowsSql, tag)
	} else {
		rowsSql := fmt.Sprintf("select id,title,published_at from %s where status = 1 order by published_at desc", postTable)
		rows, err = db.Query(rowsSql)
	}

	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	posts := make(map[string][]Post, 0)
	total := 0
	for rows.Next() {
		total += 1
		var post Post
		rows.Scan(&post.Id, &post.Title, &post.PublishedAt)
		year := post.PublishedAt.Format("2006")
		posts[year] = append(posts[year], post)
	}
	return posts, total
}

// 获取所有设置
func GetSettings() map[string]string {
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select `key`,`value` from %s", settingTable)
	rows, err := db.Query(rowsSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	settingsMap := make(map[string]string)
	for rows.Next() {
		setting := Setting{}
		rows.Scan(&setting.Key, &setting.Value)
		settingsMap[setting.Key] = setting.Value
	}
	return settingsMap
}

// 获取所有tag
func GetTags() []Tag {
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select t.id,t.name,count(p.tag_id) as post_count from %s as t right join %s as p on t.id = p.tag_id group by p.tag_id, t.id,t.name order by post_count desc;", tagTable, postTagTable)
	rows, err := db.Query(rowsSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	tags := make([]Tag, 0)
	for rows.Next() {
		tag := Tag{}
		rows.Scan(&tag.Id, &tag.Name, &tag.PostCount)
		tags = append(tags, tag)
	}
	maxWeight := tags[0].PostCount
	for index, tag := range tags {
		proportion := maxWeight / tag.PostCount
		switch {
		case proportion <= 1:
			tag.WeightCss = 15
		case proportion <= 2 && proportion > 1:
			tag.WeightCss = 10
		case proportion <= 3 && proportion > 2:
			tag.WeightCss = 5
		default:
			tag.WeightCss = 1
		}
		tags[index] = tag
	}
	return tags
}

// 获取所有友情链接
func GetFriends() []Friends {
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select id,name,url from %s", friendsTable)
	rows, err := db.Query(rowsSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	friends := make([]Friends, 0)
	for rows.Next() {
		friend := Friends{}
		rows.Scan(&friend.Id, &friend.Name, &friend.Url)
		friends = append(friends, friend)
	}
	return friends
}


// 获取所有文章的id
func GetAllPostId() []string {
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select id from `%s` where status = 1", postTable)
	rows, err := db.Query(rowsSql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	idList := make([]string, 0)
	for rows.Next() {
		var id string
		rows.Scan(&id)
		idList = append(idList, id)
	}
	return idList
}

// 更新评论数量
func UpdateCommentCount(id int, count int) int64 {
	condition := make([]map[string]string, 0)
	condition = append(condition, map[string]string{
		"column": "id",
		"value":  strconv.Itoa(id),
	})
	target := map[string]string{
		"comments_count": strconv.Itoa(count),
	}
	rowsCount := updateRows(postTable, condition, target)
	return rowsCount
}

// 更新
func updateRows(table string, conditions []map[string]string, targets map[string]string) int64 {
	if len(targets) == 0 || len(conditions) == 0 {
		return 0
	}
	whereList := make([]string, 0)
	for _, condition := range conditions {
		var express string
		if _, ok := condition["express"]; ok {
			express = condition["express"]
		} else {
			express = "="
		}
		whereChild := fmt.Sprintf("%s %s %s", condition["column"], express, condition["value"])
		whereList = append(whereList, whereChild)
	}
	whereStr := strings.Join(whereList, " and ")
	var updateList []string
	for key, value := range targets {
		updateList = append(updateList, fmt.Sprintf("%s = %s", key, value))
	}
	updateStr := strings.Join(updateList, ",")
	rowsSql := fmt.Sprintf("update %s set %s where %s", table, updateStr, whereStr)
	db := sl.MysqlClient
	stmt, err := db.Prepare(rowsSql)
	if err != nil {
		return 0
	}
	res, err := stmt.Exec()
	if err != nil {
		return 0
	}
	defer stmt.Close()
	num, err := res.RowsAffected()
	if err != nil {
		return 0
	}
	return num
}
