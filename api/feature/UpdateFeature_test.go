package feature

import (
	"context"
	"featurez/api"
	"featurez/messages"
	"featurez/services"
	"io"
	"strings"
	"testing"

	redisMock "github.com/go-redis/redismock/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFeature(t *testing.T) {
	assert := assert.New(t)

	db, mock := redisMock.NewClientMock()

	expectedFeature := "key"
	expectedValue := "value"
	expectedMessage := "Updated feature"

	mock.ExpectSet(expectedFeature, expectedValue, 0)

	request, err := jsoniter.Marshal(&messages.UpdateFeatureRequest{
		Name:  expectedFeature,
		Value: expectedValue,
	})
	assert.NoError(err)

	message := io.NopCloser(strings.NewReader(string(request)))
	redisMock := &services.RedisService{
		Client: db,
	}

	respJson, err := UpdateFeature(context.Background(), message, redisMock)
	assert.NoError(err)

	var resp *messages.UpdateFeatureResponse
	err = jsoniter.Unmarshal(respJson, &resp)
	assert.NoError(err)

	assert.Equal(expectedMessage, resp.Message)
}

func TestValidateUpdateFeature(t *testing.T) {
	assert := assert.New(t)

	expectedFeature := "key"
	expectedValue := "value"

	request, err := jsoniter.Marshal(&messages.UpdateFeatureRequest{
		Name:  expectedFeature,
		Value: expectedValue,
	})
	assert.NoError(err)

	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	reqMsg, err := validateUpdateFeature(mockMessage)

	assert.NoError(err)
	assert.Equal(expectedFeature, reqMsg.Name)
	assert.Equal(expectedValue, reqMsg.Value)
}

func TestValidateUpdateFeature_ErrInvalidParameter(t *testing.T) {
	assert := assert.New(t)
	request, err := jsoniter.Marshal(&messages.UpdateFeatureRequest{
		Name: "",
	})
	assert.NoError(err)
	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	reqMsg, err := validateUpdateFeature(mockMessage)

	assert.EqualError(err, api.ErrInvalidParameter.Error())
	assert.Nil(reqMsg)
}
