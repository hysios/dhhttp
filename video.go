package dhhttp

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/hysios/log"
)

type Frame struct {
	b           bytes.Buffer
	ContentType string
	Size        int
	// Mime
}

func (frame *Frame) Read(p []byte) (n int, err error) {
	return frame.b.Read(p)
}

func (client *Client) MjpegStream(channel int, subtype int) (chan Frame, error) {
	var (
		uri    = client.api("/cgi-bin/mjpg/video.cgi")
		method = "GET"
		query  = uri.Query()
		ch     = make(chan Frame)
	)

	query.Set("channel", fmt.Sprint(channel))
	query.Set("subtype", fmt.Sprint(subtype))
	uri.RawQuery = query.Encode()

	log.Infof("uri %s", uri.String())
	res, err := client.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	mediaType, params, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if _err := recover(); _err != nil {
				return
			}
		}()

		log.Infof("mediaType %s with params %v", mediaType, params)
		if strings.HasPrefix(mediaType, "multipart/") {
			mr := multipart.NewReader(res.Body, params["boundary"])
			for {

				p, err := mr.NextPart()
				if err == io.EOF {
					return
				}

				if err != nil {
					return
				}
				log.Infof("header %v", p.Header)
				if p.Header.Get("Content-Type") != "image/jpeg" {
					continue
				}
				var frame Frame
				b, err := io.ReadAll(p)
				if err != nil {
					continue
				}
				frame.b.Write(b)
				frame.ContentType = p.Header.Get("Content-Type")
				frame.Size, _ = strconv.Atoi(p.Header.Get("Content-Length"))

				ch <- frame
			}
		}
	}()

	return ch, nil
}
