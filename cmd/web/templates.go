package main

import (
	"aphelios_website/internal/models"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
)

type templateData struct {
	Page            models.Page
	HTML            template.HTML
	Flash           string
	IsAuthenticated bool //判断用户是否登陆
	Title           string
	RootPath        string
	TILs            models.TILs
	SideBars        bool
	Form            any
	CSRFToken       string
	Blogs           models.Blogs
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

//本质上是一个字符串键映射，它充当查找表，将名称映射到函数。
var functions = template.FuncMap{
	"formatDate": formatDate,
}

// 如果当前请求来自已认证用户，则返回 true，否则返回 false
// 返回 false。
func (app *application) isAuthenticated(r *http.Request) bool {
	//return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) NewtemplateData(r *http.Request) templateData {
	return templateData{
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

//将所有模版加载
func newTemplates() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			page,
		}

		tpl, err := template.New(name).Funcs(functions).ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = tpl
	}

	return cache, nil

}

func (app *application) templateExcute(w http.ResponseWriter, page string, status int, data templateData) {
	tpl, ok := app.templateCache[page]
	if !ok {
		fmt.Printf("the template %v dones not exist", page)
		return
	}

	buf := new(bytes.Buffer)

	err := tpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		fmt.Printf("excutetemplate err:%s", err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
