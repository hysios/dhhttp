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
	"dev.cspdls.com/pkg/storage"
	"github.com/eventials/go-tus"
	"github.com/kr/pretty"
	"github.com/segmentio/ksuid"

	_ "time/tzdata"
)

var (
	addr        string
	user        string
	pass        string
	output      string
	upload      bool
	s3Endpont   string
	s3appkey    string
	s3secret    string
	s3bucket    string
	deviceNo    string
	test        bool
	fixed       bool
	rebootCmd   bool
	shutdownCmd bool
)

var (
	store storage.Storage
)

func init() {
	flag.StringVar(&addr, "addr", "192.168.1.108", "大华球机服务器")
	flag.StringVar(&user, "user", "admin", "用户名")
	flag.StringVar(&pass, "pass", "admin123", "密码")
	flag.StringVar(&output, "output", "./images", "输出目录")

	flag.BoolVar(&upload, "upload", false, "上传到服务器")
	flag.StringVar(&s3Endpont, "s3-endpont", "8.134.49.62:9000", "s3 服务地址")
	flag.StringVar(&s3appkey, "s3-appkey", "minioadmin", "s3 Appkey")
	flag.StringVar(&s3secret, "s3-secret", "minioadmin", "s3 Secret")
	flag.StringVar(&s3bucket, "s3-bucket", "forensics", "S3 Bucket")
	flag.StringVar(&deviceNo, "deviceNO", "TEST", "设备编号")
	flag.BoolVar(&test, "test", false, "测试数据上传")
	flag.BoolVar(&fixed, "fixed-tz", false, "修正时区")

	flag.BoolVar(&rebootCmd, "reboot", false, "重启服务")
	flag.BoolVar(&shutdownCmd, "shutdown", false, "关机服务")
}

func main() {
	flag.Parse()

	if test {
		doTest()
		return
	}

	if fixed {
		initTimefix()
	}

	log.Infof("current utc time %s", time.Now().UTC())

	client, err := dhhttp.New(addr, user, pass)
	if err != nil {
		log.Fatalf("connect addr %s error %s", addr, err)
	}

	if rebootCmd {
		log.Infof("reboot %v", client.Reboot())
		os.Exit(-1)
	}

	if shutdownCmd {
		log.Infof("shutdown %v", client.Shutdown())
		os.Exit(-1)
	}

	client.AudioInput()

	os.MkdirAll(output, os.ModeDir|os.ModePerm)

	repeat := func(fn func()) {
		safeRun := func() {
			defer func() {
				r := recover()
				switch err := r.(type) {
				case error:
					log.Errorf("error %s", err)
				case string:
					log.Errorf("error %s", err)
				case nil:
				default:
					log.Errorf("error %s", err)
				}

			}()

			fn()
		}
		for {
			safeRun()
			time.Sleep(10 * time.Second)
		}
	}

	repeat(func() {
		err := client.SubscribeFunc(1, []string{"All"}, func(_events dhhttp.Events, image []byte) {
			if len(_events.Events) == 0 {
				return
			}

			var events dhhttp.EventInfo = _events.Events[0]

			log.Infof("events % #v", pretty.Formatter(events))

			t := time.Now()
			log.Infof("eventType %s", events.TrafficCar.Event)
			switch events.TrafficCar.Event {
			case "TrafficParking":
				// filename := filepath.Join(output, fmt.Sprintf("%d.json", t.UnixNano()))
				buf, err := json.Marshal(events)
				if err != nil {
					log.Infof("marshal events error %s", err)
				}

				saveOrUpload(&KeyGenerator{
					DeviceNo:  deviceNo,
					CreatedAt: t,
					Plate:     false,
					Prefix:    "",
					Ext:       "json",
				}, buf)

				if len(image) > 0 {
					saveOrUpload(&KeyGenerator{
						DeviceNo:  deviceNo,
						CreatedAt: t,
						Plate:     false,
						Prefix:    "combine",
						Ext:       "jpeg",
					}, image)
				}
			}
		})
		if err != nil {
			log.Errorf("subscribe exit %v", err)
		}
	})
}

func initStorage() {
	var err error
	store, err = storage.NewMinioV2(s3appkey, s3secret, s3bucket, storage.MinioEndpoint(s3Endpont))
	if err != nil {
		log.Fatalf("打开 minio s3 失败 %s", err)
	}
}

func initTimefix() {
	// loc, _ := time.LoadLocation("Asia/Shanghai")
	loc := time.FixedZone("UTC-8", -28800)
	time.Local = loc
	time.UTC = loc

}

type KeyGenerator struct {
	DeviceNo  string
	CreatedAt time.Time
	Plate     bool
	Prefix    string
	Ext       string
}

func (key *KeyGenerator) String() string {
	return fmt.Sprintf("%s/%4d/%02d/%02d/%s-%s.%s", key.DeviceNo,
		key.CreatedAt.Year(),
		key.CreatedAt.Month(),
		key.CreatedAt.Day(),
		key.Prefix,
		key.UUID(),
		key.ExtFile(),
	)
}

func (key *KeyGenerator) UUID() string {
	return ksuid.New().String()
}

func (key *KeyGenerator) ExtFile() string {
	if len(key.Ext) > 0 {
		return key.Ext
	}

	return "jpg"
}

func (key *KeyGenerator) PlateOrPrefix() string {
	if key.Plate {
		return "plate"
	}

	return key.Prefix
}

func saveOrUpload(key *KeyGenerator, b []byte) error {
	if upload {
		log.Infof("upload to %s", key.String())
		log.Infof("result %v", key.String())

		// create the tus client.
		client, err := tus.NewClient("http://8.134.49.62:1080/files", nil)
		if err != nil {
			return fmt.Errorf("new client error %s", err)
		}

		// create an upload from a file.
		upload := tus.NewUploadFromBytes(b)
		upload.Metadata["filename"] = key.String()

		// create the uploader.
		uploader, err := client.CreateUpload(upload)
		if err != nil {
			return fmt.Errorf("create upload error %s", err)
		}
		// start the uploading process.
		log.Info(uploader.Upload())
	} else {
		t := key.CreatedAt
		filename := filepath.Join(output, fmt.Sprintf("%d.json", t.UnixNano()))
		log.Info(ioutil.WriteFile(filename, b, os.ModePerm))
	}
	return nil
}

func doTest() {
	initStorage()
	log.Infof("test result %v", store.Put("test", []byte("this is test content")))
}
