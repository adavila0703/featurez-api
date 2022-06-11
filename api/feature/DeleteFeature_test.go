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

func TestDeleteFeature(t *testing.T) {
	assert := assert.New(t)

	db, mock := redisMock.NewClientMock()
	redisMock := &services.RedisService{
		Client: db,
	}

	expectedFeatures := []string{"feature1", "feature2"}
	expectedNotFound := []string{"feature3", "feature4"}
	expectedMessage := "keys deleted"

	mock.ExpectDel(expectedFeatures[0]).SetVal(0)
	mock.ExpectDel(expectedFeatures[1]).SetVal(0)
	mock.ExpectDel(expectedNotFound[0]).SetVal(1)
	mock.ExpectDel(expectedNotFound[1]).SetVal(1)

	request, err := jsoniter.Marshal(&messages.DeleteFeatureRequest{
		Name: append(expectedFeatures, expectedNotFound...),
	})
	assert.NoError(err)

	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	respJson, err := DeleteFeature(context.Background(), mockMessage, redisMock)
	assert.NoError(err)

	var resp *messages.DeleteFeatureResponse
	err = jsoniter.Unmarshal(respJson, &resp)
	assert.NoError(err)

	assert.Equal(expectedFeatures, resp.Deleted)
	assert.Equal(expectedNotFound, resp.NotFound)
	assert.Equal(expectedMessage, resp.Message)
}

func TestValidateDeleteFeature(t *testing.T) {
	assert := assert.New(t)

	expectedFeature := []string{"feature"}
	request, err := jsoniter.Marshal(&messages.DeleteFeatureRequest{
		Name: expectedFeature,
	})
	assert.NoError(err)

	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	reqMsg, err := validateDeleteFeature(mockMessage)

	assert.NoError(err)
	assert.Equal(expectedFeature, reqMsg.Name)
}

func TestValidateDeleteFeature_ErrInvalidParameter(t *testing.T) {
	assert := assert.New(t)
	request, err := jsoniter.Marshal(&messages.DeleteFeatureRequest{
		Name: []string{},
	})
	assert.NoError(err)
	mockMessage := io.NopCloser(strings.NewReader(string(request)))
	reqMsg, err := validateDeleteFeature(mockMessage)

	assert.EqualError(err, api.ErrInvalidParameter.Error())
	assert.Nil(reqMsg)
}
