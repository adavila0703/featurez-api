package settings

import (
	"context"
	"featurez/messages"
	"featurez/models"
	"featurez/services"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetSettings(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	assert.NoError(err)

	postgresMock, err := gorm.Open("postgres", db)
	assert.NoError(err)

	services.PostgresDB.Client = postgresMock
	defer db.Close()

	var settings *models.Settings

	expectedAddr := "localhost:1000"

	mock.ExpectQuery(regexp.QuoteMeta(`select * from settings`)).
		WithArgs(settings).
		WillReturnRows(sqlmock.NewRows([]string{"redis_address"}).AddRow(expectedAddr))

	request, err := jsoniter.Marshal(&messages.GetSettingsRequest{})
	assert.NoError(err)
	message := io.NopCloser(strings.NewReader(string(request)))

	respJSON, err := GetSettings(context.Background(), message, nil)
	assert.NoError(err)

	var resp *messages.GetSettingsResponse
	err = jsoniter.Unmarshal(respJSON, &resp)
	assert.NoError(err)

	assert.Equal(expectedAddr, resp.RedisAddress)
}
