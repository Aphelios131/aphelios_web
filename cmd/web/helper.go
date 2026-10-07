package main

import (
	"errors"
	"net/http"

	"github.com/go-playground/form/v4"
)

func (app *application) decodePostForm(r *http.Request, dst any) error {

	//解析 URL query + POST body
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		//遍历 `err` 的错误链，看里面是否包含 `*form.InvalidDecoderError`：
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}

		return err
	}

	return nil
}
