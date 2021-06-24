package dhhttp

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hysios/log"
)

var dateLayout = "2006-01-02 15:04:05"

func qtime(t time.Time) string {
	return strings.Replace(t.Format(dateLayout), " ", "%20", 1)
}

func (client *Client) DownloadMedia(channel int, startTime, endTime time.Time, subType int) (io.ReadCloser, error) {
	var (
		uri    = client.api("cgi-bin/loadfile.cgi")
		method = "GET"
	)

	uri.RawQuery = fmt.Sprintf("action=startLoad&channel=%d&startTime=%s&endTime=%s&subtype=%d",
		channel, qtime(startTime), qtime(endTime), subType)

	log.Infof("url %s", uri.String())
	res, err := client.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
