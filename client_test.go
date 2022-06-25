package dhhttp

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestClient_AudioInput(t *testing.T) {
	client, err := New("192.168.1.108", "admin", "admin123")
	if err != nil {
		t.Fatalf("new client %s", err)
	}

	cap, err := client.AudioInput()
	if err != nil {
		t.Fatalf("get audio input error %s", err)
	}
	t.Logf("audioinput % #v", pretty.Formatter(cap))
	assert.NoError(t, err)
	// assert.Equal(t, cap, "")
}

func TestClient_decodeStream(t *testing.T) {
	var (
		host = os.Getenv("DHHTTP_HOST")
	)
	var (
		buf    bytes.Buffer
		resp   = &http.Response{Body: ioutil.NopCloser(bufio.NewReader(&buf)), Header: make(http.Header)}
		mw     = multipart.NewWriter(&buf)
		client = Client{Addr: host, test: true}
	)

	resp.Header.Set("Content-Type", "multipart/mixed; boundary="+mw.Boundary())
	p, err := mw.CreatePart(textproto.MIMEHeader{"Content-Type": []string{"text/plain"}})
	if err != nil {
		t.Fatalf("create part error %s", err)
	}

	io.WriteString(p, `Events[0].CommInfo.Country=CHN
Events[0].CountInGroup=1
Events[0].EventBaseInfo.Action=Start
Events[0].EventBaseInfo.Code=TrafficParking
Events[0].EventBaseInfo.Index=0
Events[0].FrameSequence=32855766
Events[0].GroupID=1647090670
Events[0].IndexInGroup=1
Events[0].Lane=1
Events[0].LowerSpeedLimit=0
Events[0].MachineName=V1FC049A
Events[0].Mark=6
Events[0].Object.Action=Appear
Events[0].Object.BoundingBox[0]=5954
Events[0].Object.BoundingBox[1]=2234
Events[0].Object.BoundingBox[2]=6346
Events[0].Object.BoundingBox[3]=2404
Events[0].Object.Category=Normal
Events[0].Object.Center[0]=6150
Events[0].Object.Center[1]=2319
Events[0].Object.Confidence=0
Events[0].Object.Country=CHN
Events[0].Object.FrameSequence=0
Events[0].Object.MainColor[0]=0
Events[0].Object.MainColor[1]=0
Events[0].Object.MainColor[2]=255
Events[0].Object.MainColor[3]=0
Events[0].Object.MainSeat.DriverCalling=unknow
Events[0].Object.MainSeat.DriverSmoking=unknow
Events[0].Object.MainSeat.SafeBelt=unknow
Events[0].Object.ObjectID=1647090670
Events[0].Object.ObjectType=Plate
Events[0].Object.OriginalBoundingBox[0]=2791
Events[0].Object.OriginalBoundingBox[1]=646
Events[0].Object.OriginalBoundingBox[2]=2975
Events[0].Object.OriginalBoundingBox[3]=695
Events[0].Object.RelativeID=1
Events[0].Object.SerialUUID=
Events[0].Object.SlaveSeat.DriverCalling=unknow
Events[0].Object.SlaveSeat.DriverSmoking=unknow
Events[0].Object.SlaveSeat.SafeBelt=unknow
Events[0].Object.Source=0.000000
Events[0].Object.Speed=0
Events[0].Object.SpeedTypeInternal=0
Events[0].Object.Text=鄂A06654
Events[0].Object.TrackType=0
Events[0].PTS=67737170160.000000
Events[0].PreAlarm=0
Events[0].RedLightMargin=0
Events[0].RedLightUTCMS=0
Events[0].Sequence=0
Events[0].SnapFlag=ABC123
Events[0].SnapTimes=3
Events[0].Source=-1.000000
Events[0].Speed=0
Events[0].SpeedLimit[0]=70
Events[0].SpeedLimit[1]=70
Events[0].SpeedingPercentage=0
Events[0].StartParking=1615554696
Events[0].TrafficCar.CountInGroup=1
Events[0].TrafficCar.DeviceAddress=
Events[0].TrafficCar.Direction=0
Events[0].TrafficCar.DrivingDirection[0]=Approach
Events[0].TrafficCar.Event=TrafficParking
Events[0].TrafficCar.GroupID=1647090670
Events[0].TrafficCar.IndexInGroup=1
Events[0].TrafficCar.Lane=1
Events[0].TrafficCar.LowerSpeedLimit=0
Events[0].TrafficCar.MachineAddress=
Events[0].TrafficCar.MachineGroup=
Events[0].TrafficCar.MachineName=V1FC049A
Events[0].TrafficCar.OverSpeedMargin=0
Events[0].TrafficCar.PhysicalLane=1
Events[0].TrafficCar.PlateColor=Blue
Events[0].TrafficCar.PlateNumber=鄂A06654
Events[0].TrafficCar.PlateType=Normal
Events[0].TrafficCar.RoadwayNo=
Events[0].TrafficCar.SnapshotMode=4
Events[0].TrafficCar.Speed=0
Events[0].TrafficCar.UnderSpeedMargin=0
Events[0].TrafficCar.UpperSpeedLimit=70
Events[0].TrafficCar.VehicleColor=Black
Events[0].TrafficCar.VehicleLength=0.000000
Events[0].TrafficCar.VehicleSign=Unknown
Events[0].TrafficCar.VehicleSize=Light-duty
Events[0].TrafficCar.ViolationCode=1044
Events[0].TrafficCar.ViolationDesc=A类违法停车
Events[0].UTC=1615554670
Events[0].UTCMS=619
Events[0].UpperSpeedLimit=70
Events[0].Vehicle.Action=Appear
Events[0].Vehicle.BoundingBox[0]=2640
Events[0].Vehicle.BoundingBox[1]=1872
Events[0].Vehicle.BoundingBox[2]=5864
Events[0].Vehicle.BoundingBox[3]=6432
Events[0].Vehicle.BrandYear=0
Events[0].Vehicle.CarLogoIndex=0
Events[0].Vehicle.CarSeriesIndex=0
Events[0].Vehicle.Category=Unknown
Events[0].Vehicle.Center[0]=4248
Events[0].Vehicle.Center[1]=4152
Events[0].Vehicle.Confidence=0
Events[0].Vehicle.DriverFace[0][0]=0
Events[0].Vehicle.DriverFace[0][1]=0
Events[0].Vehicle.DriverFace[0][2]=0
Events[0].Vehicle.DriverFace[0][3]=0
Events[0].Vehicle.DriverFace[1][0]=0
Events[0].Vehicle.DriverFace[1][1]=0
Events[0].Vehicle.DriverFace[1][2]=0
Events[0].Vehicle.DriverFace[1][3]=0
Events[0].Vehicle.Extra.directionRun=0
Events[0].Vehicle.FrameSequence=0
Events[0].Vehicle.MainColor[0]=0
Events[0].Vehicle.MainColor[1]=0
Events[0].Vehicle.MainColor[2]=0
Events[0].Vehicle.MainColor[3]=0
Events[0].Vehicle.MainSeat.DriverCalling=unknow
Events[0].Vehicle.MainSeat.DriverSmoking=unknow
Events[0].Vehicle.MainSeat.SafeBelt=unknow
Events[0].Vehicle.MainSeat.SunShade=unknow
Events[0].Vehicle.ObjectID=1647090670
Events[0].Vehicle.ObjectType=Vehicle
Events[0].Vehicle.Posture=0
Events[0].Vehicle.RelativeID=1
Events[0].Vehicle.SerialUUID=
Events[0].Vehicle.SlaveSeat.DriverCalling=unknow
Events[0].Vehicle.SlaveSeat.DriverSmoking=unknow
Events[0].Vehicle.SlaveSeat.SafeBelt=unknow
Events[0].Vehicle.SlaveSeat.SunShade=unknow
Events[0].Vehicle.Source=0.000000
Events[0].Vehicle.Speed=0
Events[0].Vehicle.SpeedTypeInternal=0
Events[0].Vehicle.SubBrand=0
Events[0].Vehicle.Text=Unknown
Events[0].VehicleLength=0.000000

`)
	client.decodeStream(resp, 1, func(event Events) {
		log.Printf("event % #v", pretty.Formatter(event))
	})
}

func TestClient_GetPtzStatus(t *testing.T) {
	var (
		host = os.Getenv("DHHTTP_HOST")
		user = os.Getenv("DHHTTP_USER")
		pwd  = os.Getenv("DHHTTP_PASSWORD")
	)
	client, err := New(host, user, pwd)
	if err != nil {
		t.Fatalf("new client %s", err)
	}

	status, err := client.GetPtzStatus(1)
	if err != nil {
		t.Fatalf("get ptz status error %s", err)
	}
	t.Logf("status % #v", pretty.Formatter(status))
	assert.NoError(t, err)
	// assert.Equal(t, cap, "")
}

func TestClient_GetPresets(t *testing.T) {
	var (
		host = os.Getenv("DHHTTP_HOST")
		user = os.Getenv("DHHTTP_USER")
		pwd  = os.Getenv("DHHTTP_PASSWORD")
	)
	client, err := New(host, user, pwd)
	if err != nil {
		t.Fatalf("new client %s", err)
	}

	presets, err := client.GetPresets(1)
	if err != nil {
		t.Fatalf("get ptz presets error %s", err)
	}
	t.Logf("presets % #v", pretty.Formatter(presets))
	assert.NoError(t, err)
	// assert.Equal(t, cap, "")
}

func TestClient_GetMachineName(t *testing.T) {
	var (
		host = os.Getenv("DHHTTP_HOST")
		user = os.Getenv("DHHTTP_USER")
		pwd  = os.Getenv("DHHTTP_PASSWORD")
	)
	client, err := New(host, user, pwd)
	if err != nil {
		t.Fatalf("new client %s", err)
	}

	name, err := client.GetMachineName()
	if err != nil {
		t.Fatalf("get machine name error %s", err)
	}
	t.Logf("name %s", name)
	assert.NoError(t, err)
	// assert.Equal(t, cap, "")
}
