package dhhttp

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"dev.cspdls.com/pkg/log"
)

func (client *Client) CurrentTime() (time.Time, error) {
	var (
		uri    = client.api("cgi-bin/global.cgi?action=getCurrentTime.cgi")
		method = "GET"
		reply  TimeReply
	)

	log.Infof("url %s", uri.String())
	res, err := client.newRequest(method, uri, nil)
	if err != nil {
		return time.Time{}, err
	}

	if err = client.decodeResp(res, &reply); err != nil {
		return time.Time{}, err
	}

	return time.Parse("2006-01-02 15:04:05", reply.Result)
}

func (client *Client) SetCurrentTime(t time.Time) error {
	var (
		uri    = client.api("cgi-bin/global.cgi")
		method = "GET"
	)
	uri.RawQuery = url.PathEscape(fmt.Sprintf("action=setCurrentTime&time=%s", t.Format("2006-01-02 15:04:05")))
	// uri.RawQuery = q.Encode() + "&time=" + strings.Replace(t.Format("2006-01-02 15:04:05"), " ", "%20", 1)
	log.Infof("url %s", uri.String())
	resp, err := client.newRequest(method, uri, nil)
	if err != nil {
		return err
	}

	b, _ := ioutil.ReadAll(resp.Body)
	log.Infof("set time result %s", string(b))

	return nil
}
