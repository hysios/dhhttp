package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"dev.cspdls.com/pkg/dhhttp"
	"dev.cspdls.com/pkg/log"
	"github.com/kr/pretty"
)

var (
	addr   string
	user   string
	pass   string
	output string
)

func init() {
	flag.StringVar(&addr, "addr", "192.168.1.108", "大华球机服务器")
	flag.StringVar(&user, "user", "admin", "用户名")
	flag.StringVar(&pass, "pass", "admin123", "密码")
	flag.StringVar(&output, "output", "./images", "输出目录")
}

func main() {
	flag.Parse()

	client, err := dhhttp.New(addr, user, pass)
	if err != nil {
		log.Fatalf("connect addr %s error %s", addr, err)
	}

	client.AudioInput()
	ch, err := client.Subscribe(1, "TrafficParking", "All")
	if err != nil {
		log.Fatalf("subscribe exit %v", err)
	}

	os.MkdirAll(output, os.ModeDir|os.ModePerm)

	for event := range ch {

		if len(event.Events.Events) == 0 {
			continue
		}
		var events dhhttp.EventInfo = event.Events.Events[0]

		log.Infof("events % #v", pretty.Formatter(events))

		t := time.Now()
		log.Infof("eventType %s", events.TrafficCar.Event)
		switch events.TrafficCar.Event {
		case "TrafficParking":
			filename := filepath.Join(output, fmt.Sprintf("%d.json", t.UnixNano()))
			buf, err := json.Marshal(events)
			if err != nil {
				log.Infof("marshal events error %s", err)
			}
			log.Info(ioutil.WriteFile(filename, buf, os.ModePerm))

			if len(event.Image) > 0 {
				imagename := filepath.Join(output, fmt.Sprintf("%d.jpeg", t.UnixNano()))
				log.Info(ioutil.WriteFile(imagename, event.Image, os.ModePerm))
			}
		}
	}
}
