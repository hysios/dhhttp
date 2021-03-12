package dhhttp

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/hysios/digest"
	"github.com/pkg/errors"
)

type Client struct {
	Addr     string
	Username string
	Password string

	trans *digest.Transport
	test  bool
}

func New(addr string, user, pass string) (*Client, error) {
	return &Client{Addr: addr, Username: user, Password: pass, trans: digest.NewTransport(user, pass)}, nil
}

func (client *Client) api(apidir string) url.URL {
	u, _ := parseUrl(client.Addr)
	u.Path = path.Join(u.Path, apidir)
	return *u
}

func parseUrl(addr string) (*url.URL, error) {
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return url.Parse(addr)
	}
	return url.Parse("http://" + addr)
}

func (client *Client) newRequest(method string, uri url.URL, body io.Reader) (*http.Response, error) {
	// b, _ := bodyLen(body)
	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "hkhttp: new request failed")
	}

	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	// req.Header.Add("Content-Length", strconv.Itoa(n))

	res, err := client.trans.RoundTrip(req)
	if err != nil {
		return nil, errors.Wrap(err, "dhhttp: digest roundTrip request failed")
	}
	return res, nil
}

func errString(rd io.Reader) string {
	b, _ := ioutil.ReadAll(rd)
	return string(b)
}

func (client *Client) decodeResp(res *http.Response, val interface{}) error {
	dec := NewDecoder(res.Body)

	if res.StatusCode == 200 {
		if err := dec.Decode(val); err != nil {
			return errors.Wrap(err, "hkhttp: decode result failed")
		}
		return nil
	} else if res.StatusCode == 401 {
		return errors.New("invalid authority: " + errString(res.Body))
	} else {
		return errors.New("server error: " + errString(res.Body))
	}
}

func (client *Client) decodeBytes(buf []byte, val interface{}) error {
	dec := NewDecoder(bytes.NewBuffer(buf))

	if err := dec.Decode(val); err != nil {
		return errors.Wrap(err, "hkhttp: decode bytes failed")
	}
	return nil
}

func (client *Client) AudioInput() (int, error) {
	var (
		reply  AudioInputReply
		uri    = client.api("cgi-bin/devAudioInput.cgi?action=getCollect")
		method = "GET"
	)

	res, err := client.newRequest(method, uri, nil)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if err = client.decodeResp(res, &reply); err != nil {
		return 0, err
	}

	return reply.Result, err
}

type EventCallback func(events Events, buf []byte)

func (client *Client) extractCallback(args []interface{}) ([]string, EventCallback) {
	var other []string

	if len(args) < 2 {
		return nil, nil
	}

	_cb := args[len(args)-1]

	cb, _ := _cb.(EventCallback)

	for i := 0; i < len(args)-1; i++ {
		other = append(other, args[i].(string))
	}

	return other, cb
}

func (client *Client) SubscribeFunc(channel int, eventCode []string, cb EventCallback) error {
	var (
		// uri       = client.api("cgi-bin/snapManager.cgi?action=attachFileProc&Flags[0]=Event")
		uri       = client.api("/cgi-bin/snapManager.cgi?action=attachFileProc&Flags[0]=Event&Events=[TrafficParking,VideoMotion]&channel=1&heartbeat=30")
		method    = "GET"
		heartbeat = 30
	)

	if len(eventCode) > 0 {
		uri.Query().Set("Events", strings.Join(eventCode, ","))
	}

	uri.Query().Set("channel", strconv.Itoa(channel))
	uri.Query().Set("heartbeat", strconv.Itoa(heartbeat))
	log.Printf("uri %s", &uri)
	res, err := client.newRequest(method, uri, nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if err = client.decodeStream(res, 2, cb); err != nil {
		return err
	}

	return err
}

type Event struct {
	Events Events
	Image  []byte
}

func (client *Client) Subscribe(channel int, eventCode ...string) (chan Event, error) {
	var ch = make(chan Event)

	go client.SubscribeFunc(channel, eventCode, func(events Events, buf []byte) {
		ch <- Event{Events: events, Image: buf}
	})

	return ch, nil
}

func (client *Client) decodeStream(resp *http.Response, count int, _cb interface{}) error {
	var (
		cb         = reflect.ValueOf(_cb)
		decodeType reflect.Type
	)

	t := cb.Type()
	if t.Kind() == reflect.Func && t.NumIn() >= count {
		decodeType = t.In(0)
	} else {
		return errors.New("callback function args less count")
	}
	log.Printf("decode Type %s", decodeType.Name())

	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(resp.Body, params["boundary"])
		i := 0

		var (
			args  []reflect.Value
			clear = func() {
				args = make([]reflect.Value, 0)
			}
		)

		for {

			p, err := mr.NextPart()
			if err == io.EOF {
				return nil
			}

			if client.test && p == nil {
				return nil
			}
			log.Printf("context-type %s", p.Header.Get("Content-Type"))
			switch p.Header.Get("Content-Type") {
			case "text/plain":
				buf, err := ioutil.ReadAll(p)
				if err != nil && !client.test {
					// skip
					return err
				}

				if bytes.HasPrefix(buf, []byte("Heartbeat")) {
					log.Println("heartbeat")
					continue
				}

				val := reflect.New(decodeType)
				log.Printf("frame\n%s", buf)
				if err := client.decodeBytes(buf, val.Interface()); err != nil {
					log.Printf("decode Bytes error %s", err)
					continue
					// return err
				}
				log.Printf("val %v", val.Interface())
				i++
				args = append(args, val.Elem())
			case "image/jpeg":
				buf, err := ioutil.ReadAll(p)
				if err != nil {
					log.Printf("read image  error %s", err)
					continue
					// return err
				}
				log.Printf("read %d length image", len(buf))
				i++
				args = append(args, reflect.ValueOf(buf))
				// events.Image = buf
			default:
				log.Printf("unknown Content-Type")
			}

			if i%count == 0 && len(args) >= count {
				log.Println("call")
				cb.Call(args)
				clear()
			}

			// fmt.Printf("Part %q: %q\n", p.Header.Get("Foo"), slurp)
		}
	}
	// for s.Scan() {
	// 	val := reflect.New(decodeType)
	// 	log.Printf("frame\n%s", s.Bytes())
	// 	if err := client.decodeBytes(s.Bytes(), val.Interface()); err != nil {
	// 		return err
	// 	}
	// 	cb.Call([]reflect.Value{reflect.ValueOf(val.Interface())})
	// }

	return nil
}
