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

var UpdateSettingsHandler = &api.Handler{
	F:       UpdateSettings,
	Method:  http.MethodPost,
	Request: &messages.UpdateSettingsRequest{},
}

func UpdateSettings(ctx context.Context, message io.ReadCloser, redis *services.RedisService) ([]byte, error) {
	reqMsg, err := validateUpdateSettings(message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var usrSettings *models.Settings

	services.PostgresDB.Client.First(&usrSettings)

	respObject := &messages.UpdateSettingsResponse{}

	if usrSettings.ID != 0 {
		usrSettings.RedisAddress = reqMsg.RedisAddress
		services.PostgresDB.Client.Save(&usrSettings)
		services.Redis = services.NewRedisService(reqMsg.RedisAddress)
		respObject.Message = "address has been saved"
	}

	resp, err := jsoniter.Marshal(respObject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}

func validateUpdateSettings(message io.ReadCloser) (*messages.UpdateSettingsRequest, error) {
	reqMsg, err := api.ReadMessage[messages.UpdateSettingsRequest](message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if reqMsg.RedisAddress == "" {
		return nil, errors.WithStack(api.ErrInvalidParameter)
	}

	return reqMsg, nil
}
