package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.routers(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server", "addr", srv.Addr)

	//创建一个 shutdownError 通道。我们将使用它来接收返回的任何错误
	shutdownError := make(chan error)

	//捕获 SIGINT 和 SIGTERM 信号
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Info("stopping server", "addr", srv.Addr, "signal", s.String())
		// 调用服务器上的 Shutdown() 函数，传入一个不包含任何值或截止时间的空上下文。如果优雅关闭成功，Shutdown() 将返回 nil；否则返回错误（可能是由于关闭监听器时出现问题）。我们将此返回值传递给 shutdownError 通道。
		shutdownError <- srv.Shutdown(context.Background())
	}()

	// 调用服务器上的 Shutdown() 函数会导致 ListenAndServe() 函数立即返回一个 http.ErrServerClosed 错误。因此，如果我们看到这个错误，实际上是件好事，表明优雅关闭已经开始。所以我们专门检查这个错误，只有当它不是 http.ErrServerClosed 时才返回错误。
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("shutdown complete")
	return nil
}
