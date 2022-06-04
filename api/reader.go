package api

import (
	"bytes"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func ReadMessage[T any](message io.ReadCloser) (*T, error) {
	var reqMsg *T
	var buf bytes.Buffer

	buf.ReadFrom(message)

	err := jsoniter.Unmarshal(buf.Bytes(), &reqMsg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return reqMsg, nil
}
