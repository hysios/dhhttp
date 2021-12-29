package dhhttp

type AudioInputReply struct {
	Result int
}

type Events struct {
	Events []EventInfo
}

type EventBaseInfo struct {
	Action string
	Code   string
	Index  int
}

type EventInfo struct {
	ID                 []int
	CountInGroup       int
	RegionName         []string
	EventBaseInfo      EventBaseInfo
	SmartMotionEnable  bool
	FrameSequence      int
	GroupID            string
	IndexInGroup       int
	Lane               int
	LowerSpeedLimit    float64
	MachineName        string
	Mark               int
	PTS                float64
	PreAlarm           int
	RedLightMargin     int
	RedLightUTCMS      int
	Sequence           int
	SnapFlag           string
	SnapTimes          int
	Source             float64
	Speed              int
	SpeedLimit         []int
	SpeedingPercentage int
	StartParking       int
	TrafficCar         TrafficCar
	UTC                int
	UTCMS              int
	UpperSpeedLimit    int
	VehicleLength      float64
}

type Object struct {
	ObjectID            string
	ObjectType          string
	BoundingBox         []int
	Category            string
	Center              []int
	Confidence          int
	Country             string
	FrameSequence       int
	MainColor           []int
	MainSeat            Seat
	SlaveSeat           Seat
	OriginalBoundingBox []int
	RelativeID          int
	SerialUUID          string
	Source              float64
	Speed               float64
	SpeedTypeInternal   int
	Text                string
	TrackType           int
}

type Seat struct {
	DriverCalling string
	DriverSmoking string
	SafeBelt      string
	SunShade      string
}

type TrafficCar struct {
	CountInGroup     int
	DeviceAddress    string
	Direction        int
	DrivingDirection []string
	Event            string
	GroupID          string
	IndexInGroup     int
	Lane             int
	LowerSpeedLimit  int
	MachineAddress   string
	MachineGroup     string
	MachineName      string
	OverSpeedMargin  int
	PhysicalLane     int
	PlateColor       string
	PlateNumber      string
	PlateType        string
	RoadwayNo        string
	SnapshotMode     int
	Speed            int
	UnderSpeedMargin int
	UpperSpeedLimit  int
	VehicleColor     string
	VehicleLength    float64
	VehicleSign      string
	VehicleSize      string
	ViolationCode    string
	ViolationDesc    string
}

type Vehicle struct {
	Action            string
	BoundingBox       []int
	BrandYear         int
	CarLogoIndex      int
	CarSeriesIndex    int
	Category          string
	Center            []int
	Confidence        int
	DriverFace        [][]int
	Extra             Extra
	FrameSequence     int
	MainColor         []int
	MainSeat          Seat
	ObjectID          int
	ObjectType        string
	Posture           int
	RelativeID        int
	SerialUUID        string
	SlaveSeat         Seat
	Source            float64
	Speed             int
	SpeedTypeInternal int
	SubBrand          int
	Text              string
}

type Extra struct {
	DirectionRun int
}

type CommInfo struct {
	Country string
}

type TimeReply struct {
	Result string
}
