package feature

import (
	"context"
	"featurez/api"
	"featurez/messages"
	"featurez/services"
	"fmt"
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

	expectedKey := "key"
	expectedMessage := "Feature flag has been set!"

	mock.ExpectExists(expectedKey).SetVal(0)
	mock.ExpectSet(expectedKey, 0, 0)

	mockMessage := io.NopCloser(strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, expectedKey)))
	respJson, err := CreateFeature(context.Background(), mockMessage, redisMock)
	assert.NoError(err)

	var resp *messages.CreateFeatureResponse
	err = jsoniter.Unmarshal(respJson, &resp)
	assert.NoError(err)

	assert.Equal(expectedKey, resp.FeatureFlag)
	assert.Equal(expectedMessage, resp.Message)
}

func TestCreateFeature_KeyExists(t *testing.T) {
	assert := assert.New(t)

	db, mock := redisMock.NewClientMock()
	redisMock := &services.RedisService{
		Client: db,
	}

	expectedKey := "key"

	mock.ExpectExists(expectedKey).SetVal(1)

	mockMessage := io.NopCloser(strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, expectedKey)))
	respJson, err := CreateFeature(context.Background(), mockMessage, redisMock)

	assert.EqualError(err, api.ErrFeatureAlreadyExists.Error())
	assert.Nil(respJson)
}

func TestValidateCreateFeature(t *testing.T) {
	assert := assert.New(t)
	mockMessage := io.NopCloser(strings.NewReader(fmt.Sprintf(`{"name":"test"}`)))
	reqMsg, err := validateCreateFeature(mockMessage)

	assert.NoError(err)
	assert.Equal("test", reqMsg.Name)
}

func TestValidateCreateFeature_ErrInvalidParameter(t *testing.T) {
	assert := assert.New(t)
	mockMessage := io.NopCloser(strings.NewReader(`{"name":""}`))
	reqMsg, err := validateCreateFeature(mockMessage)

	assert.EqualError(err, api.ErrInvalidParameter.Error())
	assert.Nil(reqMsg)
}
