package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routers() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticates, preventCSRF)
	mux.Handle("GET /", dynamic.ThenFunc(app.Home))

	mux.Handle("GET /admin/signup", dynamic.ThenFunc(app.SignUp))
	mux.Handle("POST /admin/signup", dynamic.ThenFunc(app.SignUpPost))
	mux.Handle("GET /admin/login", dynamic.ThenFunc(app.Login))
	mux.Handle("POST /admin/login", dynamic.ThenFunc(app.LoginPost))

	mux.Handle("GET /til", dynamic.ThenFunc(app.TodayILearn))
	mux.Handle("GET /today-i-learned/{path}", dynamic.ThenFunc(app.TodayILearnPath))

	mux.Handle("GET /blog", dynamic.ThenFunc(app.Blog))
	mux.Handle("GET /blog/{path}", dynamic.ThenFunc(app.blogPath))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("POST /admin/logout", protected.ThenFunc(app.Logout))

	//til
	mux.Handle("GET /admin/new_til", protected.ThenFunc(app.New_til))
	mux.Handle("GET /admin/til/{path}", protected.ThenFunc(app.Update_til))
	mux.Handle("POST /admin/new_til", protected.ThenFunc(app.TilsPost))
	mux.Handle("POST /admin/delete_til/{path}", protected.ThenFunc(app.Delete_til))

	//blog
	mux.Handle("GET /admin/new_blog", protected.ThenFunc(app.New_blog))
	mux.Handle("GET /admin/blog/{path}", protected.ThenFunc(app.Update_blog))
	mux.Handle("POST /admin/new_blog", protected.ThenFunc(app.BlogPost))

	return mux
}
