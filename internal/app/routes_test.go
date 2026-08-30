package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFound(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/not-existing-path/", nil)

	app.Routes().ServeHTTP(rr, req)
	assertErrorResponse(t, rr, http.StatusNotFound, "the requested resource could not be found")
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/healthcheck", nil)

	message := fmt.Sprintf("the %s method is not supported for this resourse", http.MethodPost)
	app.Routes().ServeHTTP(rr, req)
	assertErrorResponse(t, rr, http.StatusMethodNotAllowed, message)

}
