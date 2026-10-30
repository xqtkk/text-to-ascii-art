package web

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func testHandler() http.Handler {
	app := NewApplication(log.New(io.Discard, "", 0), log.New(io.Discard, "", 0))
	return app.Routes()
}

func TestHomeAndMissingRouteStatuses(t *testing.T) {
	handler := testHandler()
	for _, test := range []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/", http.StatusOK},
		{http.MethodGet, "/missing", http.StatusNotFound},
		{http.MethodGet, "/ascii-art", http.StatusMethodNotAllowed},
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(test.method, test.path, nil))
		if recorder.Code != test.status {
			t.Errorf("%s %s status = %d, want %d", test.method, test.path, recorder.Code, test.status)
		}
	}
}

func TestConvertAndExport(t *testing.T) {
	handler := testHandler()
	values := url.Values{"input": {"A"}, "font": {"standard"}}
	request := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), "YOUR RESULT") {
		t.Fatalf("POST /ascii-art status=%d, response missing result", recorder.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/export", strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("POST /export status=%d, want 200", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("export Content-Type = %q", got)
	}
	if got := recorder.Header().Get("Content-Disposition"); got != `attachment; filename="ascii-art.txt"` {
		t.Errorf("export Content-Disposition = %q", got)
	}
	if got := recorder.Header().Get("Content-Length"); got != strconv.Itoa(recorder.Body.Len()) {
		t.Errorf("export Content-Length = %q, body is %d bytes", got, recorder.Body.Len())
	}
}

func TestBadSubmissionReturns400(t *testing.T) {
	handler := testHandler()
	for _, values := range []url.Values{
		{"input": {""}, "font": {"standard"}},
		{"input": {"hello"}, "font": {"../../README.md"}},
		{"input": {"snowman ☃"}, "font": {"standard"}},
	} {
		request := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(values.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("bad submission status = %d, want 400", recorder.Code)
		}
	}
}
