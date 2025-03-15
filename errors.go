package ghb

import (
	"fmt"
	"net/http"
)

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func badRequestf(w http.ResponseWriter, format string, a ...any) {
	err := fmt.Errorf(format, a...)
	badRequest(w, err)
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func internalServerErrorf(w http.ResponseWriter, format string, a ...any) {
	internalServerError(w, fmt.Errorf("internal server error: "+format, a...))
}
