package model

import (
	"html/template"
	"time"
	"fmt"
	"database/sql"
	"gblog/helpers"
	sl "gblog/service/local"
	"strings"
)

type ArticleForm struct {
	Title string `json:"title" form:"title" binding:"required"`
	CategoryId string `json:"category_id" form:"category" binding:"required,gt=0"`
	Content string `json:"content" form:"content"`
	PublishedAt string `json:"published_at" form:"published_at"`
	Summary string `json:"summary" form:"summary"`
	Tags []string `json:"tags" form:"tags[]"`
	Status string `json:"status" form:"status"`
}

//文章
type Post struct {
	Id           int    `json:"id" form:"id"`
	Title        string `json:"title" form:"title"`
	Content      string `json:"content" form:"content"`
	ContentHtml  template.HTML
	Status       int      `json:"status" form:"status"`
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

// 获取文章列表
func GetPosts(offset int, limit int) ([]Post, int) {
	posts := make([]Post, 0)
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select p.id,p.title,p.content,p.status,p.published_at,c.id as cat_id,c.name,c.en_name,p.comments_count,p.abstract,p.content,p.created_at from %s as p left join %s as c on p.cat_id = c.id where p.status = 1 order by p.id desc limit %d,%d", postTable, categoryTable, offset, limit)
	countSql := fmt.Sprintf("select count(*) as total from %s where status = 1", postTable)
	var rows *sql.Rows
	var err error
	rows, err = db.Query(rowsSql)
	defer rows.Close()
	for rows.Next() {
		var post Post
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.Status, &post.PublishedAt,&post.CatId, &post.CName, &post.CEnName, &post.CommentCount, &post.Abstract, &post.Content, &post.CreatedAt)
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
	err = db.QueryRow(countSql).Scan(&total)
	if err != nil {
		return posts,0
	}

	return posts, total
}

// 搜索文章
func SearchPosts(keyword string, offset int, limit int)([]Post, int){
	posts := make([]Post, 0)
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select p.id,p.title,p.content,p.published_at,c.name,c.en_name,p.comments_count,p.abstract,p.content,p.created_at from %s as p left join %s as c on p.cat_id = c.id where (p.content like ? or p.title like ? ) and p.status = 1 order by p.id desc limit %d,%d", postTable, categoryTable, offset, limit)
	countSql := fmt.Sprintf("select count(*) as total from %s as p left join %s as c on p.cat_id = c.id where (p.content like ? or p.title like ? ) and p.status = 1", postTable, categoryTable)
	keyword = "%" + keyword + "%"
	rows, err := db.Query(rowsSql, keyword, keyword)
	defer rows.Close()
	if err != nil {
		return posts,0
	}
	for rows.Next(){
		var post Post
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.PublishedAt, &post.CName, &post.CEnName, &post.CommentCount, &post.Abstract, &post.Content, &post.CreatedAt)
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
	err = db.QueryRow(countSql,keyword,keyword).Scan(total)
	if err != nil {
		return posts,0
	}
	return posts, total
}



// 插入文章
func InsertPost(article_form ArticleForm, uid int64) (int64, error){
	// tag
	db := sl.MysqlClient
	type tagSt struct {
		Id	int8 `json:"id"`
		Name string `json:"name"`
	}
	// 插入tag
	tagIdList := make([]int64, 0)
	if len(article_form.Tags) > 0 {

		for _,tag := range article_form.Tags{
			var tagId int64
			tagSql := fmt.Sprintf("select id from `%s` where `name` = ?", tagTable)
			errTagQuery := db.QueryRow(tagSql, tag).Scan(&tagId)
			switch {
				case errTagQuery == sql.ErrNoRows:
					// 不存在则插入
					tagInsertSql := fmt.Sprintf("insert into `%s` (`name`) values (?)", tagTable)
					res, err := db.Exec(tagInsertSql, tag)
					if err != nil{
						return 0,err
					}
					tagId,err = res.LastInsertId()
					if err != nil{
						return 0,err
					}
				case errTagQuery != nil:
					return 0,errTagQuery
			}
			tagIdList = append(tagIdList, tagId)
		}

	}
	// 插入文章和关联tag
	insertArticleSql := fmt.Sprintf("insert into `%s` (`title`, `content`,`status`,`author_id`,`cat_id`,`published_at`,`abstract`) values (?,?,?,?,?,?,?)", postTable)
	res,err := db.Exec(insertArticleSql, article_form.Title, article_form.Content,article_form.Status,uid,article_form.CategoryId,article_form.PublishedAt,article_form.Summary)
	if err != nil {
		return 0,err
	}
	articleId,errArticleLastId := res.LastInsertId()
	if errArticleLastId != nil {
		return 0,err
	}
	// 关联文章和tag
	tagRelation := make([]string, 0)
	for _,tagId := range tagIdList{
		tagRelation = append(tagRelation, fmt.Sprintf("(%d,%d)", articleId, tagId))
	}
	tagRelationSql := fmt.Sprintf("replace into `%s` (`post_id`,`tag_id`) values %s", postTagTable, strings.Join(tagRelation, ","))
	db.Exec(tagRelationSql)
	return articleId,nil
	//sql := "INSERT INTO `%s` (`title`,`content`,`category_id`,`published_at`,`summary`,``)"
}

// 删除文章
func ArticleDelete(idList []string)bool{
	deleteSql := fmt.Sprintf("update `%s` set status = 3 where `id` in (?)", postTable)
	db := sl.MysqlClient
	_,err := db.Exec(deleteSql, strings.Join(idList, ","))
	if err != nil {
		return false
	}
	return true
}

// 获取单篇文章
func GetPost(id int, filterStatus bool) (Post,error) {
	post := Post{}
	db := sl.MysqlClient
	var statusList []int
	if filterStatus {
		statusList = []int{
			1,
		}
	}else{
		statusList = []int{
			1,2,3,
		}
	}
	rowsSql := fmt.Sprintf("select `t`.`name` as `tag_name`,`p`.`id`,`p`.`title`,`p`.`content`,`p`.`status`,`p`.`published_at`,`p`.`comments_count`,`c`.`id` as `cat_id`, `c`.`name`,`c`.`en_name`,`p`.`Abstract` from `%s` as `p` left join `%s` as `c` on `p`.`cat_id` = `c`.`id` left join `%s` as `pt` on `p`.`id` = `pt`.`post_id` left join `%s` as `t` on `t`.`id` = `pt`.`tag_id` where `p`.`id` = ? and `p`.`status` in (?)", postTable, categoryTable, postTagTable, tagTable)
	rows, err := db.Query(rowsSql, id, statusList)
	if err != nil {
		return post,err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		rows.Scan(&post.TagName, &post.Id, &post.Title, &post.Content, &post.Status, &post.PublishedAt, &post.CommentCount, &post.CatId, &post.CName, &post.CEnName, &post.Abstract)
		if post.Abstract == "" {
			if len([]rune(post.Content)) > 250 {
				post.Abstract = helpers.TrimHtmlTag(string([]rune(post.Content)[:250]))
			} else {
				post.Abstract = helpers.TrimHtmlTag(string([]rune(post.Content)))
			}

		}
		tag.Name = post.TagName
		post.Tags = append(post.Tags, tag)
	}

	return post,err
}