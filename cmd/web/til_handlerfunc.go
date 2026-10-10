package main

import (
	"aphelios_website/internal/models"
	"aphelios_website/validator"
	"fmt"
	"log"
	"net/http"
)

type TilForm struct {
	Til models.Til
	validator.Validator
}

func (app *application) TodayILearn(w http.ResponseWriter, r *http.Request) {
	data := app.NewtemplateData(r)
	data.Title = "Today I Learn"
	data.RootPath = "/today-i-learned"
	tils, err := app.tils.GetAll()
	if err != nil {
		log.Println(err)
		return
	}

	data.TILs = tils
	data.SideBars = false
	app.templateExcute(w, "tils.html", 200, data)
}

func (app *application) TodayILearnPath(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	til, err := app.tils.GetByPath(path)
	if err != nil {
		fmt.Printf("app.til.GetBy(%v): %v", path, err)
		http.NotFound(w, r)
		return
	}

	data := app.NewtemplateData(r)
	data.Form = TilForm{Til: til}
	app.templateExcute(w, "til_detail.html", 200, data)

}

// 显示新增til界面
func (app *application) New_til(w http.ResponseWriter, r *http.Request) {
	data := app.NewtemplateData(r)
	data.Form = TilForm{Til: models.Til{}}
	app.templateExcute(w, "new_til.html", 200, data)
}

func (app *application) Update_til(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	til, err := app.tils.GetByPath(path)
	if err != nil {
		fmt.Printf("GetByPath err:%v", err)
		return
	}
	data := app.NewtemplateData(r)
	data.Form = TilForm{Til: til}
	app.templateExcute(w, "new_til.html", 200, data)
}

func (app *application) TilsPost(w http.ResponseWriter, r *http.Request) {
	// var TilsForm struct {
	// 	Id int `form:"til-id"`
	// 	Path string `form:"til-path"`
	// 	Title string `form:"til-title"`
	// 	Category string `form:"til-category"`
	// 	Summary string `form:"til-summary"`
	// 	Text string `form:"til-text"`
	// }

	var TilsForm TilForm

	err := app.decodePostForm(r, &TilsForm.Til)
	if err != nil {
		fmt.Printf("decodePostForm err:%v", err)
		return
	}

	TilsForm.CheckField(validator.NotBlack(TilsForm.Til.Path), "path", "This field cannot be blank")
	TilsForm.CheckField(validator.MaxChars(TilsForm.Til.Path, 20), "path", "This field cannot be more than 20 characters long")

	TilsForm.CheckField(validator.NotBlack(TilsForm.Til.Title), "title", "This field cannot be blank")
	TilsForm.CheckField(validator.MaxChars(TilsForm.Til.Title, 20), "Title", "This field cannot be more than 20 characters long")

	TilsForm.CheckField(validator.NotBlack(TilsForm.Til.Category), "category", "This field cannot be blank")
	TilsForm.CheckField(validator.MaxChars(TilsForm.Til.Category, 10), "category", "This field cannot be more than 10 characters long")

	TilsForm.CheckField(validator.NotBlack(TilsForm.Til.Summary), "summary", "This field cannot be blank")
	TilsForm.CheckField(validator.MaxChars(TilsForm.Til.Summary, 1000), "summary", "This field cannot be more than 1000 characters long")

	TilsForm.CheckField(validator.NotBlack(TilsForm.Til.Text), "text", "This field cannot be blank")
	TilsForm.CheckField(validator.MaxChars(TilsForm.Til.Text, 1000), "text", "This field cannot be more than 1000 characters long")

	if !TilsForm.Valid() {
		data := app.NewtemplateData(r)
		data.Form = TilsForm
		app.templateExcute(w, "new_til.html", http.StatusUnprocessableEntity, data)
		return
	}

	if TilsForm.Til.Id == 0 {
		err = app.tils.Insert(&TilsForm.Til)
		if err != nil {
			log.Printf("app.til.Insert(): %v\n", err)
			// TODO: send some notification to the UI
		}
	} else {
		err = app.tils.Update(&TilsForm.Til)
		if err != nil {
			log.Printf("app.til.UpdateExisting(): %v\n", err)
			// TODO: send some notification to the UI
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/today-i-learned/%v", TilsForm.Til.Path), http.StatusSeeOther)

}

func (app *application) Delete_til(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	err := app.tils.Delete(path)
	if err != nil {
		fmt.Printf("delete err:%v", err)
		return
	}

	http.Redirect(w, r, "/til", http.StatusSeeOther)

}
