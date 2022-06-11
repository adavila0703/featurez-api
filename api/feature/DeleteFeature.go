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

func DeleteFeature(ctx context.Context, message io.ReadCloser, redis *services.RedisService) ([]byte, error) {
	reqMsg, err := validateDeleteFeature(message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var deleted []string
	var notFound []string

	for _, featureName := range reqMsg.Name {
		found, err := redis.Delete(ctx, featureName)

		if found == 0 {
			deleted = append(deleted, featureName)
		} else if found == 1 {
			notFound = append(notFound, featureName)
		}

		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	respObject := &messages.DeleteFeatureResponse{
		Deleted:  deleted,
		NotFound: notFound,
		Message:  "keys deleted",
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
