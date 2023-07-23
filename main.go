package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/gordonklaus/portaudio"

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

	portaudio.Initialize()
	defer portaudio.Terminate()
	buffer := make([]float32, sampleRate*seconds)
	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(buffer), func(in []float32) {
		for i := range buffer {
			buffer[i] = in[i]
		}
	})
	if err != nil {
		panic(err)
	}
	stream.Start()
	defer stream.Close()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
	}).Name = "index"

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
