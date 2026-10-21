package models

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Blogs []*Blog

type Blog struct {
	Id       int    `form:"blog-id"`
	Path     string `form:"blog-path"`
	Title    string `form:"blog-title"`
	Category string `form:"blog-category"`
	Summary  string `form:"blog-summary"`
	Text     string `form:"blog-text"`
	HTML     struct {
		Text    template.HTML
		Summary template.HTML
	}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlogModel struct {
	DB *sql.DB
}

func (Bm *BlogModel) GetAll() (Blogs, error) {
	stmt := `select id, path, title, category, summary, text, created_at, updated_at from  blog order by created_at desc;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := Bm.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blogs []*Blog
	for rows.Next() {
		var blog Blog
		err = rows.Scan(
			&blog.Id,
			&blog.Path,
			&blog.Title,
			&blog.Category,
			&blog.Summary,
			&blog.Text,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		htmlBytes := blackfriday.Run([]byte(blog.Summary))
		blog.HTML.Summary = template.HTML(htmlBytes)
		blogs = append(blogs, &blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (bm *BlogModel) GetByPath(path string) (Blog, error) {
	stmt := `select id, title, category, summary, text, created_at, updated_at from blog 
	where path = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := bm.DB.QueryRowContext(ctx, stmt, path)

	var blog Blog
	err := row.Scan(
		&blog.Id,
		&blog.Title,
		&blog.Category,
		&blog.Summary,
		&blog.Text,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil {
		fmt.Printf("data scanf err:%v", err)
		return blog, err
	}

	htmlBytes := blackfriday.Run([]byte(blog.Text))
	blog.HTML.Text = template.HTML(htmlBytes)
	blog.Path = path

	return blog, nil

}

func (bm *BlogModel) Insert(blog Blog) error {
	stmt := `insert into blog (path, title, category, summary, text) values
	(?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := bm.DB.ExecContext(ctx, stmt, blog.Path, blog.Title, blog.Category, blog.Summary, blog.Text)
	if err != nil {
		fmt.Printf("insert exec err:%v", err)
		return err
	}

	return nil
}

func (bm *BlogModel) Update(blog Blog) error {
	stmt := `update 
	blog SET path = ?,
		title = ?,
		category = ?,
		summary = ?,
		text = ?,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := bm.DB.ExecContext(ctx, stmt, blog.Path, blog.Title, blog.Category, blog.Summary, blog.Text, blog.Id)
	if err != nil {
		return err
	}

	return nil
}
