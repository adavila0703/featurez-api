package api

import (
	"bytes"
	"context"
	"featurez/config"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Handler struct {
	F       func(ctx context.Context, message io.ReadCloser) ([]byte, error)
	Method  string
	Request interface{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if h.Method != r.Method {
		http.Error(w, "incorrect method", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	resp, err := h.F(ctx, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)

	logger := &Logger{
		StatusCode:  http.StatusOK,
		Method:      r.Method,
		RequestType: fmt.Sprint(reflect.TypeOf(h.Request)),
	}

	if config.Cfg.DebugMode {
		var buf bytes.Buffer
		buf.ReadFrom(r.Body)
		logger.RequestData = buf.String()
		logger.ResponseData = string(resp)
	}

	logger.LogInfo()
}
