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

var GetFeatureListHandler = &api.Handler{
	F:       GetFeatureList,
	Method:  http.MethodGet,
	Request: &messages.GetFeatureListRequest{},
}

func GetFeatureList(ctx context.Context, message io.ReadCloser) ([]byte, error) {
	keys, err := services.Redis.GetAllKeys(ctx)
	if err != nil {
		return nil, err
	}

	var featuresList []*messages.Feature

	for _, key := range keys {
		value, err := services.Redis.GetValues(ctx, key)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		feature := &messages.Feature{
			Name:  key,
			Value: value,
		}

		featuresList = append(featuresList, feature)
	}

	respObject := &messages.GetFeatureListResponse{
		FeatureList: featuresList,
	}

	resp, err := jsoniter.Marshal(respObject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}
