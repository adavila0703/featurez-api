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

var CreateFeatureHandler = &api.Handler{
	F:       CreateFeature,
	Method:  http.MethodPost,
	Request: &messages.CreateFeatureRequest{},
}

func CreateFeature(ctx context.Context, message io.ReadCloser, redis *services.RedisService) ([]byte, error) {
	reqMsg, err := validateCreateFeature(message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result, err := redis.Exists(ctx, reqMsg.Name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if result == 1 {
		return nil, api.ErrFeatureAlreadyExists
	}

	redis.SetKey(ctx, reqMsg.Name, 0)

	respObject := &messages.CreateFeatureResponse{
		Message:     "Feature flag has been set!",
		FeatureFlag: reqMsg.Name,
	}

	resp, err := jsoniter.Marshal(respObject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}

func validateCreateFeature(message io.ReadCloser) (*messages.CreateFeatureRequest, error) {
	reqMsg, err := api.ReadMessage[messages.CreateFeatureRequest](message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(reqMsg.Name) == 0 {
		return nil, errors.WithStack(api.ErrInvalidParameter)
	}

	return reqMsg, nil
}
