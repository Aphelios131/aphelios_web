package models

import (
	"database/sql"
	"html/template"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Id           int
	Name         string
	Version      int
	Content      string
	CreatedAt    time.Time
	HTMLContent  template.HTML
}

type PageModel struct {
	DB *sql.DB
}

// func (pg *PageModel) Insert(page *Page) error {
// 	var GetPage Page
// 	stmt := `select COALESCE(MAX(version), 0) from aphelios_website.pages where name = $1`
// 	row := pg.DB.QueryRow(stmt, page.Name)
// 	err := row.Scan(&GetPage.Version)
// 	if err != nil {
// 		return fmt.Errorf("row.Scan(&page.Version): %v", err)
// 	}

// 	page.Version += 1
// 	stmt := `insert into aphelios_website.pages(name, version, content) values($1, $2, $3)`
// 	return nil
// }

func (pg *PageModel) Get(name string) (*Page, error) {
	var GetPage Page
	stmt := `select id, name, version, content, created_at from pages where name = $1 order by version desc limit 1`
	row := pg.DB.QueryRow(stmt, name)
	err := row.Scan(
		&GetPage.Id,
		&GetPage.Name,
		&GetPage.Version,
		&GetPage.Content,
		&GetPage.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	//读取 Page 里的 Markdown 原文，用 blackfriday 把 md 翻译成 HTML，
	// 包装成 template.HTML 类型，传给 Go 模板，让页面正常渲染富文本，
	// 而不是原样输出 html 源码。
	htmlBytes := blackfriday.Run([]byte(GetPage.Content))
	GetPage.HTMLContent = template.HTML(htmlBytes)

	return &GetPage, nil
}
