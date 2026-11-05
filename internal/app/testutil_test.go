package app

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"

	"greenlight.dekutyavin.net/internal/config"
	"greenlight.dekutyavin.net/internal/data"
)

type testMovieModel struct {
	insertFn func(*data.Movie) error
	getFn    func(int) (*data.Movie, error)
	updateFn func(*data.Movie) error
	deleteFn func(int) error
	getAllFn func(string, []string, data.Filters) ([]*data.Movie, data.Metadata, error)
}

func (m *testMovieModel) Insert(mov *data.Movie) error { return m.insertFn(mov) }

func (m *testMovieModel) Get(id int) (*data.Movie, error) { return m.getFn(id) }

func (m *testMovieModel) Update(mov *data.Movie) error { return m.updateFn(mov) }

func (m *testMovieModel) Delete(id int) error { return m.deleteFn(id) }

func (m *testMovieModel) GetAll(title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
	return m.getAllFn(title, genres, filters)
}

func newTestMovieModel(t *testing.T) *testMovieModel {
	t.Helper()
	return &testMovieModel{
		insertFn: func(m *data.Movie) error { return nil },
		getFn:    func(i int) (*data.Movie, error) { return nil, data.ErrRecordNotFound },
		updateFn: func(m *data.Movie) error { return nil },
		deleteFn: func(i int) error { return nil },
		getAllFn: func(s1 string, s2 []string, f data.Filters) ([]*data.Movie, data.Metadata, error) {
			return []*data.Movie{}, data.Metadata{}, nil
		},
	}
}

func newTestApplication(t *testing.T) *Application {
	t.Helper()

	return &Application{
		Config: config.Config{
			Env: "testing",
		},
		Logger: slog.New(slog.DiscardHandler),
		Models: Models{Movies: newTestMovieModel(t)},
	}
}

func newTestApplicationWithLogs(t *testing.T) (*Application, *bytes.Buffer) {
	t.Helper()
	buf := &bytes.Buffer{}
	app := &Application{
		Config: config.Config{Env: "testing"},
		Logger: slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Models: Models{Movies: newTestMovieModel(t)},
	}

	return app, buf
}

func assertLogError(t *testing.T, buf *bytes.Buffer, errMsg string, method string, uri string) {
	t.Helper()

	var entry struct {
		Level  string `json:"level"`
		Msg    string `json:"msg"`
		Method string `json:"method"`
		URI    string `json:"uri"`
	}

	err := json.Unmarshal(buf.Bytes(), &entry)
	if err != nil {
		t.Fatalf("log output is not valid json: %v\n%s", err, buf.String())
	}

	if entry.Level != slog.LevelError.String() {
		t.Errorf("want level: %q; got: %q", slog.LevelError.String(), entry.Level)
	}

	if entry.Msg != errMsg {
		t.Errorf("want msg: %q; got: %q", errMsg, entry.Msg)
	}

	if entry.Method != method {
		t.Errorf("want method: %q; got: %q", method, entry.Method)
	}

	if entry.URI != uri {
		t.Errorf("want URI: %q; got: %q", uri, entry.URI)
	}
}

func assertErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, wantStatusCode int, message any) {
	t.Helper()

	wantBody, err := json.MarshalIndent(envelope{"error": message}, "", "\t")
	if err != nil {
		t.Fatalf("could not marshal wantBody: %v", err)
	}
	wantBody = append(wantBody, '\n')

	if rr.Result().StatusCode != wantStatusCode {
		t.Errorf("want status: %d; got: %d", wantStatusCode, rr.Result().StatusCode)
	}

	gotBody := rr.Body.String()
	if gotBody != string(wantBody) {
		t.Errorf("want body: %q; got: %q", string(wantBody), gotBody)
	}
}
