package settings

import (
	"context"
	"featurez/api"
	"featurez/clients"
	"featurez/messages"
	"featurez/models"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var UpdateSettingsHandler = &api.Handler{
	F:      UpdateSettings,
	Method: http.MethodPost,
}

func UpdateSettings(ctx context.Context, message io.ReadCloser) ([]byte, error) {
	reqMsg, err := validateUpdateSettings(message)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var usrSettings *models.Settings

	clients.PostgresDB.Client.First(&usrSettings)

	respObject := &messages.UpdateSettingsResponse{}

	if usrSettings.ID != 0 {
		usrSettings.RedisAddress = reqMsg.RedisAddress
		clients.PostgresDB.Client.Save(&usrSettings)
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
