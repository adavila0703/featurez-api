package feature

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFeatureList(t *testing.T) {
	assert := assert.New(t)

	message := io.NopCloser(strings.NewReader("message"))
	resp, err := GetFeatureList(context.Background(), message)

	assert.NoError(err)
	fmt.Println(resp)
}
