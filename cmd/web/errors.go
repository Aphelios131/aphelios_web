package main

import "net/http"

func (app *application) logError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "url", url)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// errorResponse() 方法是一个辅助方法，用于向客户端发送带有给定状态码的 JSON 格式错误消息。请注意，我们对 message 参数使用了 any 类型，而不是字符串类型，因为这样可以让我们更灵活地控制响应中包含的值
// func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
// 	type envelope struct {
// 		error_message any
// 	}

// 	env := envelope{
// 		error_message: message,
// 	}

// 	err := app.
// }
