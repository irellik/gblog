package model

import (
	"fmt"
	sl "github.com/irellik/gblog/service/local"
)

// 栏目
type Category struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"name"`
	EnName      string `json:"en_name" form:"en_name"`
	PostCount   int    `json:"post_count"`
}

// 获取所有栏目
func GetCategories() []Category {
	category := Category{}
	categories := make([]Category, 0)
	db := sl.MysqlClient
	rowsSql := fmt.Sprintf("select `c`.`id`, `c`.`name`, `c`.`description`, `c`.`en_name`, `p`.`count` as post_count from `%s` as `c` left join (select cat_id,count(*) as count from `%s` where status = ? group by cat_id) as `p` on `c`.`id` = `p`.`cat_id` where `c`.status = 1 order by `c`.`id` asc;", categoryTable, postTable)
	rows, err := db.Query(rowsSql, 1)
	if err != nil {
		return categories
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&category.Id, &category.Name, &category.Description, &category.EnName, &category.PostCount)
		categories = append(categories, category)
	}
	return categories
}
