package main

import (
	"bytes"
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

// CleanMarkdownSource 统一清洗markdown原始文本，文件上传、表单编辑都调用它
func CleanMarkdownSource(raw string) string {
	b := []byte(raw)
	// 顺序：先 \r\n，再单独 \r，去掉BOM
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	b = bytes.ReplaceAll(b, []byte("\r"), []byte("\n"))
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	return string(b)
}
