package dhhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_MjpegStream(t *testing.T) {
	var client, err = New("192.168.1.108", "admin", "admin123")

	assert.NoError(t, err)
	assert.NotNil(t, client)

	frameCh, err := client.MjpegStream(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, frameCh)
}
