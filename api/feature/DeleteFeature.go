package feature

import (
	"context"
	"featurez/api"
	"featurez/messages"
	"featurez/services"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var DeleteFeatureHandler = &api.Handler{
	F:       DeleteFeature,
	Method:  http.MethodDelete,
	Request: &messages.DeleteFeatureRequest{},
}

func DeleteFeature(ctx context.Context, message io.ReadCloser) ([]byte, error) {
	reqMsg, err := validateDeleteFeature(message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, featureName := range reqMsg.Name {
		_, err := services.Redis.Delete(ctx, featureName)

		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	respObject := &messages.DeleteFeatureResponse{
		Message: "keys deleted",
	}

	resp, err := jsoniter.Marshal(respObject)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func validateDeleteFeature(message io.ReadCloser) (*messages.DeleteFeatureRequest, error) {
	reqMsg, err := api.ReadMessage[messages.DeleteFeatureRequest](message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(reqMsg.Name) < 1 {
		return nil, errors.WithStack(api.ErrInvalidParameter)
	}

	return reqMsg, nil
}
