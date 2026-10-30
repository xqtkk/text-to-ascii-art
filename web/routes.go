package web

import (
	"io/fs"
	"net/http"

	assets "text-to-ascii-art"
)

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.Home)
	mux.HandleFunc("POST /ascii-art", app.Output)
	// Retain the original route so existing links continue to work.
	mux.HandleFunc("POST /text-to-ascii-art", app.Output)
	mux.HandleFunc("POST /export", app.Export)
	staticFiles, err := fs.Sub(assets.Files, "ui/static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	return mux
}
