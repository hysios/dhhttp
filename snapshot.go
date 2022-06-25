package dhhttp

import (
	"fmt"
	"io"

	"github.com/hysios/log"
)

func (client *Client) Snapshot(channel int) ([]byte, error) {
	var (
		uri    = client.api("/cgi-bin/snapshot.cgi")
		method = "GET"
		query  = uri.Query()
	)

	query.Set("channel", fmt.Sprint(channel))
	uri.RawQuery = query.Encode()

	log.Infof("uri %s", uri.String())
	res, err := client.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}
