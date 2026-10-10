package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// 中间件 保护路由 确认访问路由需要登陆 没有登陆返回登陆界面
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 如果用户未通过身份验证，则将其重定向到登录页面并
		// 从中间件链返回，以便后续处理程序不再运行
		// 链式结构将被执行。
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		// 需要身份验证的信息不会存储在用户的浏览器缓存中（或其他中间缓存）。
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// 中间件 判断用户登陆,将isAuthenticatedContextKey改为true
func (app *application) authenticates(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.users.GetId(id)
		if err != nil {
			fmt.Printf("Getid err:%v", err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func preventCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
