package main

import (
	"go.uber.org/zap"
	"html/template"
	"io"
	"monowork/monowork"
	"net/http"

	"github.com/labstack/echo/v4"
)

const sampleRate = 44100
const seconds = 1

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Static("/assets", "www/dist/assets")

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("www/dist/*.html")),
	}

	e.Renderer = renderer

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
	}).Name = "index"

	reader := monowork.NewReader()
	logger, _ := zap.NewProduction()
	ws := monowork.NewWebSocket(reader, &logger)

	// e.GET("/stream", func(c echo.Context) error {
	// 	c.Response().Header().Set("Connection", "Keep-Alive")
	// 	c.Response().Header().Set("Transfer-Encoding", "chunked")
	// 	for true {
	// 		binary.Write(c.Response().Writer, binary.BigEndian, &buffer)
	// 		flusher.Flush() // Trigger "chunked" encoding
	// 		return
	// 	}
	// }).Name = "index"
	e.Logger.Fatal(e.Start(":8000"))
}
