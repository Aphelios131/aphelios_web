package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/russross/blackfriday/v2"
)

type TILs []*Til

type Til struct {
	Id       int    `form:"til-id"`
	Path     string `form:"til-path"`
	Title    string `form:"til-title"`
	Category string `form:"til-category"`
	Summary  string `form:"til-summary"`
	Text     string `form:"til-text"`
	HTML     struct {
		Text    template.HTML
		Summary template.HTML
	}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TilModel struct {
	DB *sql.DB
}

func (tilm *TilModel) GetAll() (TILs, error) {
	stmt := `select id, path, title, category, summary, created_at, updated_at from til order by created_at desc`

	rows, err := tilm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tils TILs

	for rows.Next() {
		var til Til
		err = rows.Scan(
			&til.Id,
			&til.Path,
			&til.Title,
			&til.Category,
			&til.Summary,
			&til.CreatedAt,
			&til.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		htmlBytes := blackfriday.Run([]byte(til.Summary))
		til.HTML.Summary = template.HTML(htmlBytes)
		tils = append(tils, &til)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tils, err
}

//路径查询
func (tilm *TilModel) GetByPath(path string) (Til, error) {
	stmt := `select id, title, category, summary, text, created_at, updated_at from til where path = $1`
	row := tilm.DB.QueryRow(stmt, path)

	var Til Til
	err := row.Scan(&Til.Id, &Til.Title, &Til.Category, &Til.Summary, &Til.Text, &Til.CreatedAt, &Til.UpdatedAt)
	if err != nil {
		fmt.Printf("data scanf err:%v", err)
		return Til, err
	}

	htmlBytes := blackfriday.Run([]byte(Til.Text))
	Til.HTML.Text = template.HTML(htmlBytes)
	Til.Path = path
	return Til, nil

}

func (tilm *TilModel) Insert(t *Til) error {
	stmt := `insert into til (path, title, category, summary, text)
	values ($1,$2,$3,$4,$5)`
	_, err := tilm.DB.Exec(stmt, t.Path, t.Title, t.Category, t.Summary, t.Text)
	if err != nil {
		return err
	}

	return nil
}

func (tilm *TilModel) Update(t *Til) error {
	stmt := `update til SET path = $1,
		title = $2,
		category = $3,
		summary = $4,
		text = $5,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $6;`
	_, err := tilm.DB.Exec(stmt, t.Path, t.Title, t.Category, t.Summary, t.Text, t.Id)
	if err != nil {
		return err
	}

	return nil
}

func (tilm *TilModel) Delete(path string) error {
	stmt := `delete from til where path = $1`
	_, err := tilm.DB.Exec(stmt, path)
	if err != nil {
		return err
	}

	return nil
}
