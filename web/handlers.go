package web

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"text-to-ascii-art/program"
)

const maxRequestBytes = 16 << 10
const maxInputRunes = 500

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.renderError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.renderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	app.render(w, http.StatusOK, Data{Font: "standard"})
}

func (app *Application) Output(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.renderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	data, status, message := app.convertRequest(w, r)
	if status != http.StatusOK {
		app.renderError(w, status, message)
		return
	}
	app.render(w, http.StatusOK, data)
}

func (app *Application) Export(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.renderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	data, status, message := app.convertRequest(w, r)
	if status != http.StatusOK {
		app.renderError(w, status, message)
		return
	}
	content := []byte(data.Output)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func (app *Application) convertRequest(w http.ResponseWriter, r *http.Request) (Data, int, string) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	if err := r.ParseForm(); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return Data{}, http.StatusRequestEntityTooLarge, "The submitted form is too large."
		}
		return Data{}, http.StatusBadRequest, "The form could not be read. Please try again."
	}
	input := strings.ReplaceAll(r.FormValue("input"), "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	if input == "" {
		return Data{}, http.StatusBadRequest, "Enter some text to convert."
	}
	if !utf8.ValidString(input) || utf8.RuneCountInString(input) > maxInputRunes {
		return Data{}, http.StatusBadRequest, "Text must be valid UTF-8 and no longer than 500 characters."
	}
	font := r.FormValue("font")
	if font == "" {
		font = "standard"
	}
	output, err := program.TextToASCIIArt(input, font)
	if err != nil {
		return Data{}, http.StatusBadRequest, err.Error()
	}
	return Data{Input: input, Output: output, Font: font}, http.StatusOK, ""
}

func (app *Application) render(w http.ResponseWriter, status int, data Data) {
	var body bytes.Buffer
	if err := app.templates.ExecuteTemplate(&body, "home.html", data); err != nil {
		app.errorLog.Printf("render home template: %v", err)
		app.renderError(w, http.StatusInternalServerError, "Something went wrong while rendering the page.")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = io.Copy(w, &body)
}

func (app *Application) renderError(w http.ResponseWriter, status int, message string) {
	var body bytes.Buffer
	data := ApplicationError{Message: message, Code: status}
	if err := app.templates.ExecuteTemplate(&body, "error.html", data); err != nil {
		app.errorLog.Printf("render error template: %v", err)
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = io.Copy(w, &body)
}
