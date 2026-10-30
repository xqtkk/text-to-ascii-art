package web

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	assets "text-to-ascii-art"
)

type Data struct {
	Input  string
	Output string
	Font   string
}

type ApplicationError struct {
	Message string
	Code    int
}

type Application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	templates *template.Template
}

func NewApplication(errorLog, infoLog *log.Logger) *Application {
	return &Application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		templates: template.Must(template.ParseFS(assets.Files, "templates/*.html")),
	}
}

func NewServer(addr *string, errorLog *log.Logger, mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:              *addr,
		ErrorLog:          errorLog,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func Web(addr *string) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	app := NewApplication(errorLog, infoLog)
	server := NewServer(addr, errorLog, app.Routes())

	stop, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	serverErrors := make(chan error, 1)
	go func() {
		infoLog.Printf("Server listening at http://localhost%s", *addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			errorLog.Printf("server: %v", err)
		}
	case <-stop.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			errorLog.Printf("server shutdown: %v", err)
		}
		infoLog.Println("Server stopped")
	}
}
