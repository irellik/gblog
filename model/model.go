package model

import (
	"database/sql"
	"fmt"
	"github.com/irellik/gblog/helpers"
	sl "github.com/irellik/gblog/service/local"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"
)

//文章

type Post struct {
	Id           int    `json:"id" form:"id"`
	Title        string `json:"title" form:"title"`
	Content      string `json:"content" form:"content"`
	ContentHtml  template.HTML
	Status       int8      `json:"status" form:"status"`
	AuthorId     int       `json:"author_id" form:"author_id"`
	CatId        int       `json:"cat_id" form:"cat_id"`
	PublishedAt  time.Time `json:"published_at" form:"published_at"`
	CreatedAt    time.Time `json:"created_at" form:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" form:"updated_at"`
	Abstract     string    `json:"abstract" form:"abstract"`
	CommentCount uint64    `json:"comments_count" form:"comments_count"`
	TagImg       string    `json:"tag_img" form:"tag_img"`
	TagName      string    `json:"tag_name"`
	CName        string    `json:"name" form:"name"`
	CEnName      string    `json:"en_name" form:"en_name"`
	Tags         []Tag     `json:"tag"`
}

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

// 栏目
type Category struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"name"`
	EnName      string `json:"en_name" form:"en_name"`
	PostCount   int    `json:"post_count"`
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

// 获取文章列表
func GetPosts(category string, offset int, limit int, search bool) ([]Post, int) {
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select p.id,p.title,p.content,p.published_at,c.name,c.en_name,p.comments_count,p.abstract,p.content from %s as p left join %s as c on p.cat_id = c.id where p.status = 1 order by p.id desc limit %d,%d", postTable, categoryTable, offset, limit)
	countSql := fmt.Sprintf("select count(*) as total from %s where status = 1", postTable)
	var rows *sql.Rows
	var err error
	if category != "" {
		if search {
			rowsSql = fmt.Sprintf("select p.id,p.title,p.content,p.published_at,c.name,c.en_name,p.comments_count,p.abstract,p.content from %s as p left join %s as c on p.cat_id = c.id where p.content like ? and p.status = 1 order by p.id desc limit %d,%d", postTable, categoryTable, offset, limit)
			countSql = fmt.Sprintf("select count(*) as total from %s as p left join %s as c on p.cat_id = c.id where p.content like ? and p.status = 1", postTable, categoryTable)
			category = "%" + category + "%"
		} else {
			rowsSql = fmt.Sprintf("select p.id,p.title,p.content,p.published_at,c.name,c.en_name,p.comments_count,p.abstract,p.content from %s as p left join %s as c on p.cat_id = c.id where c.en_name = ? and p.status = 1 order by p.id desc limit %d,%d", postTable, categoryTable, offset, limit)
			countSql = fmt.Sprintf("select count(*) as total from %s as p left join %s as c on p.cat_id = c.id where c.en_name = ? and p.status = 1", postTable, categoryTable)
		}
		rows, err = db.Query(rowsSql, category)
	} else {
		rows, err = db.Query(rowsSql)
	}

	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.PublishedAt, &post.CName, &post.CEnName, &post.CommentCount, &post.Abstract, &post.Content)
		if post.Abstract == "" {
			if len([]rune(post.Content)) > 250 {
				post.Abstract = helpers.TrimHtmlTag(string([]rune(post.Content)[:250]))
			} else {
				post.Abstract = helpers.TrimHtmlTag(string([]rune(post.Content)))
			}

		}
		posts = append(posts, post)
	}
	// 获取文章总数
	var total int
	if category != "" {
		err = db.QueryRow(countSql, category).Scan(&total)
	} else {
		err = db.QueryRow(countSql).Scan(&total)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return posts, total
}

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

// 获取所有栏目
func GetCategories() []Category {
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select `c`.`id`, `c`.`name`, `c`.`description`, `c`.`en_name`, `p`.`count` as post_count from `%s` as `c` left join (select cat_id,count(*) as count from `%s` where status = ? group by cat_id) as `p` on `c`.`id` = `p`.`cat_id` where `c`.status = 1 order by `c`.`id` asc;", categoryTable, postTable)
	rows, err := db.Query(rowsSql, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	categories := make([]Category, 0)
	for rows.Next() {
		category := Category{}
		rows.Scan(&category.Id, &category.Name, &category.Description, &category.EnName, &category.PostCount)
		categories = append(categories, category)
	}
	return categories
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

// 获取单篇文章
func GetPost(id int) Post {
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select `t`.`name` as `tag_name`,`p`.`id`,`p`.`title`,`p`.`content`,`p`.`published_at`,`p`.`comments_count`,`c`.`name`,`c`.`en_name` from `%s` as `p` left join `%s` as `c` on `p`.`cat_id` = `c`.`id` left join `%s` as `pt` on `p`.`id` = `pt`.`post_id` left join `%s` as `t` on `t`.`id` = `pt`.`tag_id` where `p`.`id` = ?", postTable, categoryTable, postTagTable, tagTable)
	rows, err := db.Query(rowsSql, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	post := Post{}
	for rows.Next() {
		var tag Tag
		rows.Scan(&post.TagName, &post.Id, &post.Title, &post.Content, &post.PublishedAt, &post.CommentCount, &post.CName, &post.CEnName)
		tag.Name = post.TagName
		post.Tags = append(post.Tags, tag)
	}

	return post
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
	fmt.Printf("update %s set %s where %s", table, updateStr, whereStr)
	db := sl.MysqlClient
	stmt, _ := db.Prepare(rowsSql)
	res, _ := stmt.Exec()
	num, _ := res.RowsAffected()
	return num
}
