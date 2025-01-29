package err

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Message string `json:"error"`
}

func newErrorResponse(err error) *errorResponse {
	return &errorResponse{Message: err.Error()}
}

func (e *errorResponse) asBytes() ([]byte, error) {
	return json.Marshal(e)
}

var (
	RespErrJsonEncode = errors.New("json encode failure")
	RespErrJsonDecode = errors.New("json decode failure")
	RespErrDbAccess   = errors.New("database access error")
	RespErrDbInsert   = errors.New("database insert error")
	RespErrDbUpdate   = errors.New("database update error")
	RespErrDbDelete   = errors.New("database delete error")
)

func writeErrorResponse(w http.ResponseWriter, err error) {
	bytes, err := newErrorResponse(err).asBytes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		slog.Info("error", slog.String("message", string(bytes)))
		w.Write(bytes)
	}
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	writeErrorResponse(w, err)
}

func ServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	writeErrorResponse(w, err)
}

func NotFound(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	writeErrorResponse(w, err)
}
