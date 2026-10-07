package main

import (
	"aphelios_website/internal/models"
	"aphelios_website/validator"
	"fmt"
	"net/http"
)

type BlogForm struct {
	Blog models.Blog
	validator.Validator
}

func (app *application) Blog(w http.ResponseWriter, r *http.Request) {
	data := app.NewtemplateData(r)
	blogs, err := app.blogs.GetAll()
	if err != nil {
		fmt.Printf("Get all err:%v", err)
		return
	}
	data.Blogs = blogs
	app.templateExcute(w, "blog.html", 200, data)

}

func (app *application) blogPath(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")

	blog, err := app.blogs.GetByPath(path)
	if err != nil {
		fmt.Printf("GetByPath err:%v", err)
		return
	}

	data := app.NewtemplateData(r)
	data.Form = BlogForm{Blog: blog}
	app.templateExcute(w, "blog_detail.html", 200, data)

}

func (app *application) New_blog(w http.ResponseWriter, r *http.Request) {
	data := app.NewtemplateData(r)
	data.Form = BlogForm{Blog: models.Blog{}}
	app.templateExcute(w, "new_blog.html", 200, data)
}

func (app *application) Update_blog(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	Blog, err := app.blogs.GetByPath(path)
	if err != nil {
		fmt.Printf("GetByPath:%v", err)
		return
	}

	data := app.NewtemplateData(r)
	data.Form = BlogForm{Blog: Blog}
	app.templateExcute(w, "new_blog.html", 200, data)
}

func (app *application) BlogPost(w http.ResponseWriter, r *http.Request) {
	var BlogForm BlogForm

	err := app.decodePostForm(r, &BlogForm.Blog)
	if err != nil {
		fmt.Printf("decodePostForm err:%v", err)
		return
	}

	BlogForm.CheckField(validator.NotBlack(BlogForm.Blog.Path), "path", "This field cannot be blank")
	BlogForm.CheckField(validator.MaxChars(BlogForm.Blog.Path, 20), "path", "This field cannot be more than 20 characters long")

	BlogForm.CheckField(validator.NotBlack(BlogForm.Blog.Title), "title", "This field cannot be blank")
	BlogForm.CheckField(validator.MaxChars(BlogForm.Blog.Title, 20), "Title", "This field cannot be more than 20 characters long")

	BlogForm.CheckField(validator.NotBlack(BlogForm.Blog.Category), "category", "This field cannot be blank")
	BlogForm.CheckField(validator.MaxChars(BlogForm.Blog.Category, 10), "category", "This field cannot be more than 10 characters long")

	BlogForm.CheckField(validator.NotBlack(BlogForm.Blog.Summary), "summary", "This field cannot be blank")
	BlogForm.CheckField(validator.MaxChars(BlogForm.Blog.Summary, 50), "summary", "This field cannot be more than 50 characters long")

	BlogForm.CheckField(validator.NotBlack(BlogForm.Blog.Text), "text", "This field cannot be blank")
	BlogForm.CheckField(validator.MaxChars(BlogForm.Blog.Text, 500), "text", "This field cannot be more than 50 characters long")

	if !BlogForm.Valid() {
		data := app.NewtemplateData(r)
		data.Form = BlogForm
		app.templateExcute(w, "new_blog.html", http.StatusUnprocessableEntity, data)
		return
	}

	if BlogForm.Blog.Id == 0 {
		err = app.blogs.Insert(BlogForm.Blog)
		if err != nil {
			fmt.Printf("app.til.Insert(): %v\n", err)
			// TODO: send some notification to the UI
		}
	} else {
		err = app.blogs.Update(BlogForm.Blog)
		if err != nil {
			fmt.Printf("app.til.UpdateExisting(): %v\n", err)
			// TODO: send some notification to the UI
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/blog/%v", BlogForm.Blog.Path), http.StatusSeeOther)
}
