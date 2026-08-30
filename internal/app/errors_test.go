package app

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogError(t *testing.T) {
	t.Parallel()

	app, buf := newTestApplicationWithLogs(t)

	errMsg := "boom"
	method := http.MethodGet
	uri := "/v1/movies/42"
	req := httptest.NewRequest(method, uri, nil)

	app.logError(req, errors.New(errMsg))
	assertLogError(t, buf, errMsg, method, uri)
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/any-path-will-do", nil)

	wantStatusCode := http.StatusBadRequest
	message := map[string]any{
		"key": "value",
	}

	app.errorResponse(rr, req, wantStatusCode, message)
	assertErrorResponse(t, rr, wantStatusCode, message)
}

func TestServerErrorResponse(t *testing.T) {
	t.Parallel()

	method := http.MethodGet
	uri := "/v1/any-path-will-do"

	app, buf := newTestApplicationWithLogs(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(method, uri, nil)

	wantStatusCode := http.StatusInternalServerError
	errMsg := "some error message"
	message := "the server encountered a problem and could not process your request"

	app.serverErrorResponse(rr, req, errors.New(errMsg))

	assertLogError(t, buf, errMsg, method, uri)
	assertErrorResponse(t, rr, wantStatusCode, message)
}

func TestNotFoundResponse(t *testing.T) {
	t.Parallel()

	method := http.MethodGet
	uri := "/v1/any-path-will-do"

	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(method, uri, nil)

	wantStatusCode := http.StatusNotFound
	message := "the requested resource could not be found"

	app.notFoundResponse(rr, req)
	assertErrorResponse(t, rr, wantStatusCode, message)
}

func TestMethodNotAllowedResponse(t *testing.T) {
	t.Parallel()

	method := http.MethodPost
	uri := "/v1/any-path-will-do"

	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(method, uri, nil)

	wantStatusCode := http.StatusMethodNotAllowed
	message := fmt.Sprintf("the %s method is not supported for this resourse", method)

	app.methodNotAllowedResponse(rr, req)
	assertErrorResponse(t, rr, wantStatusCode, message)
}
