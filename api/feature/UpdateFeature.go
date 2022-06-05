package feature

import (
	"context"
	"featurez/api"
	"featurez/clients"
	"featurez/messages"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var UpdateFeatureHandler = &api.Handler{
	F:       UpdateFeature,
	Method:  http.MethodPost,
	Request: &messages.UpdateFeatureRequest{},
}

func UpdateFeature(ctx context.Context, message io.ReadCloser) ([]byte, error) {
	reqMsg, err := validateUpdateFeature(message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clients.Redis.SetKey(ctx, reqMsg.Name, reqMsg.Value)

	respObject := &messages.UpdateFeatureResponse{
		Message: "Updated feature",
	}

	resp, err := jsoniter.Marshal(respObject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}

func validateUpdateFeature(message io.ReadCloser) (*messages.UpdateFeatureRequest, error) {
	reqMsg, err := api.ReadMessage[messages.UpdateFeatureRequest](message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if reqMsg.Name == "" || reqMsg.Value == "" {
		return nil, errors.WithStack(api.ErrInvalidParameter)
	}

	return reqMsg, nil
}
