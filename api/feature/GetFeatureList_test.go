package feature

import (
	"context"
	"featurez/messages"
	"featurez/services"
	"io"
	"strings"
	"testing"

	redisMock "github.com/go-redis/redismock/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestGetFeatureList(t *testing.T) {
	assert := assert.New(t)

	db, mock := redisMock.NewClientMock()

	mockKey := "key"
	mockValue := "value"

	mock.ExpectKeys("*").SetVal([]string{mockKey})
	mock.ExpectGet(mockKey).SetVal(mockValue)

	request, err := jsoniter.Marshal(&messages.GetFeatureListRequest{})
	assert.NoError(err)

	message := io.NopCloser(strings.NewReader(string(request)))
	redisMock := &services.RedisService{
		Client: db,
	}

	respJson, err := GetFeatureList(context.Background(), message, redisMock)
	assert.NoError(err)

	var resp *messages.GetFeatureListResponse
	err = jsoniter.Unmarshal(respJson, &resp)
	assert.NoError(err)

	assert.Equal(mockKey, resp.FeatureList[0].Name)
	assert.Equal(mockValue, resp.FeatureList[0].Value)
}
