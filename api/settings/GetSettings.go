package settings

import (
	"context"
	"featurez/api"
	"featurez/messages"
	"featurez/models"
	"featurez/services"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var GetSettingsHandler = &api.Handler{
	F:       GetSettings,
	Method:  http.MethodGet,
	Request: &messages.GetSettingsRequest{},
}

func GetSettings(ctx context.Context, message io.ReadCloser) ([]byte, error) {
	var settings *models.Settings

	services.PostgresDB.Client.First(&settings)

	respObject := &messages.GetSettingsResponse{}

	if settings.ID != 0 {
		respObject.RedisAddress = settings.RedisAddress
	}

	resp, err := jsoniter.Marshal(respObject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}
