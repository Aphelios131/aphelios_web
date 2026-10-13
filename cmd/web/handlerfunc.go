package main

import (
	"aphelios_website/validator"
	"net/http"
)

// Get Home
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	//`ParseFiles` 支持多文件 多个文件时，**第一个文件名字作为模板名**，其他作为关联模板。
	page, err := app.pages.Get("Home")
	if err != nil {
		app.logError(r, err)
		return
	}
	data := app.NewtemplateData(r)
	data.Page = *page
	data.HTML = page.HTMLContent
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data.Flash = flash
	app.templateExcute(w, "home.html", 200, data)

	// tmp, err := template.ParseFiles("./ui/html/base.html", "./ui/html/pages/home.html", "./ui/html/partials/nav.html")
	// if err != nil {
	// 	fmt.Printf("html ParseFiles fail:%v", err)
	// 	return
	// }

	// err = tmp.ExecuteTemplate(w, "base", page)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// tpl, ok := app.templateCache["home.html"]
	// if !ok {
	// 	fmt.Printf("the template %v dones not exist", page)
	// 	return
	// }

	// err = tpl.ExecuteTemplate(w, "base", page)
	// if err != nil {
	// 	fmt.Printf("excutetemplate err:%s", err)
	// 	return
	// }
	//w.Write([]byte("This is my website"))
}

type UserSignupForm struct {
		Email    string `form:"email"`
		Password string `form:"password"`
		validator.Validator `form:"-"`
}
// Get SignUp
func (app *application) SignUp(w http.ResponseWriter, r *http.Request) {
	data := app.NewtemplateData(r)
	data.Form = UserSignupForm{}
	app.templateExcute(w, "signup.html", 200, data)
}

func (app *application) SignUpPost(w http.ResponseWriter, r *http.Request) {

	var userSignupForm UserSignupForm 
	err := app.decodePostForm(r, &userSignupForm)
	if err != nil {
		app.logError(r, err)
		return
	}

	//验证邮箱格式
	userSignupForm.CheckField(validator.Matches(userSignupForm.Email, validator.EmailRX), "Email", "Email格式不符合")

	if !userSignupForm.Valid() {
		data := app.NewtemplateData(r)
		data.Form = userSignupForm
		app.templateExcute(w, "signup.html", http.StatusUnprocessableEntity, data)
	}

	err = app.users.Insert(userSignupForm.Email, userSignupForm.Password)
	if err != nil {
		app.logError(r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// Get Login
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	data := app.NewtemplateData(r)
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data.Flash = flash
	app.templateExcute(w, "login.html", 200, data)
}

// Post login
func (app *application) LoginPost(w http.ResponseWriter, r *http.Request) {
	var UserLoginForm struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}
	err := app.decodePostForm(r, &UserLoginForm)
	if err != nil {
		app.logError(r, err)
		return
	}

	id, err := app.users.Get(UserLoginForm.Email, UserLoginForm.Password)
	if err != nil {
		app.logError(r, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.logError(r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.logError(r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
