package main

import (
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
