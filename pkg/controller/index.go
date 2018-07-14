package controller

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// Index is the root controller.
// It handles:
// - /
// - /static/** => _static/**
type Index struct {
	Log *logger.Logger
}

// Register adds routes for the controller.
func (i Index) Register(app *web.App) {
	app.ServeStatic("/static", "_static")
	app.GET("/", i.index)
}

// index handles `/`
func (i Index) index(ctx *web.Ctx) web.Result {
	return ctx.Static("_static/index.html")
}
