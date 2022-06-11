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

func TestCreateFeature(t *testing.T) {
	assert := assert.New(t)

	db, mock := redisMock.NewClientMock()
	redisMock := &services.RedisService{
		Client: db,
	}

	expectedFeature := "key"
	expectedMessage := "Feature flag has been set!"

	mock.ExpectExists(expectedFeature).SetVal(0)
	mock.ExpectSet(expectedFeature, 0, 0)

	request, err := jsoniter.Marshal(&messages.CreateFeatureRequest{
		Name: expectedFeature,
	})
	assert.NoError(err)

	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	respJson, err := CreateFeature(context.Background(), mockMessage, redisMock)
	assert.NoError(err)

	var resp *messages.CreateFeatureResponse
	err = jsoniter.Unmarshal(respJson, &resp)
	assert.NoError(err)

	assert.Equal(expectedFeature, resp.FeatureFlag)
	assert.Equal(expectedMessage, resp.Message)
}

func TestCreateFeature_KeyExists(t *testing.T) {
	assert := assert.New(t)

	db, mock := redisMock.NewClientMock()
	redisMock := &services.RedisService{
		Client: db,
	}

	expectedFeature := "key"

	mock.ExpectExists(expectedFeature).SetVal(1)

	request, err := jsoniter.Marshal(&messages.CreateFeatureRequest{
		Name: expectedFeature,
	})
	assert.NoError(err)

	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	respJson, err := CreateFeature(context.Background(), mockMessage, redisMock)

	assert.EqualError(err, api.ErrFeatureAlreadyExists.Error())
	assert.Nil(respJson)
}

func TestValidateCreateFeature(t *testing.T) {
	assert := assert.New(t)

	expectedFeature := "key"

	request, err := jsoniter.Marshal(&messages.CreateFeatureRequest{
		Name: expectedFeature,
	})
	assert.NoError(err)

	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	reqMsg, err := validateCreateFeature(mockMessage)

	assert.NoError(err)
	assert.Equal(expectedFeature, reqMsg.Name)
}

func TestValidateCreateFeature_ErrInvalidParameter(t *testing.T) {
	assert := assert.New(t)
	request, err := jsoniter.Marshal(&messages.CreateFeatureRequest{
		Name: "",
	})
	assert.NoError(err)
	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	reqMsg, err := validateCreateFeature(mockMessage)

	assert.EqualError(err, api.ErrInvalidParameter.Error())
	assert.Nil(reqMsg)
}
