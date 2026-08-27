package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var errInvalidIdParameter = errors.New("invalid id parameter")

func readIDParam(r *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		return 0, errInvalidIdParameter
	}

	return id, nil
}

func writeJSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// add extra line for nice looking responses via curl
	json = append(json, '\n')

	w.Write(json)

	return nil
}
