package dhhttp

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_DownloadMedia(t *testing.T) {
	client, err := New("192.168.1.108", "admin", "admin123")
	if err != nil {
		t.Fatalf("new client %s", err)
	}

	startTime := time.Now().Add(-30 * time.Minute)
	endTime := startTime.Add(5 * time.Second)
	media, err := client.DownloadMedia(1, startTime, endTime, 0)
	assert.NoError(t, err)
	assert.NotNil(t, media)
	b, err := ioutil.ReadAll(media)
	assert.NoError(t, err)
	assert.True(t, bytes.HasPrefix(b, []byte{'D', 'H', 'A', 'V'}))
	t.Logf("b %s", hex.Dump(b[:32]))
}
