package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverPanic(t *testing.T) {
	t.Parallel()

	app, buf := newTestApplicationWithLogs(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/panic", nil)

	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	app.recoverPanic(panicking).ServeHTTP(rr, req)

	assertErrorResponse(t, rr, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
	assertLogError(t, buf, "boom", "GET", "/panic")

	if got := rr.Result().Header.Get("Connection"); got != "close" {
		t.Errorf("want Connection: close; got := %q", got)
	}
}
