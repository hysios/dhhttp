package dhhttp

type (
	PlateColor   string
	PlateType    string
	VehicleColor string
	VehicleType  string
)

const (
	ColorOther                 PlateColor = "Other"                 // 其他颜色 "Other"
	ColorBlue                  PlateColor = "Blue"                  // 蓝色 "Blue"
	ColorYellow                PlateColor = "Yellow"                // 黄色 "Yellow"
	ColorWhite                 PlateColor = "White"                 // 白色 "White"
	ColorBlack                 PlateColor = "Black"                 // 黑色 "Black"
	ColorYellowBottomBlackText PlateColor = "YellowbottomBlackText" // 黄底黑字 "YellowbottomBlackText"
	ColorBlueBottomWhiteText   PlateColor = "BluebottomWhiteText"   // 蓝底白字 "BluebottomWhiteText"
	ColorBlackBottomWhiteText  PlateColor = "BlackBottomWhiteText"  // 黑底白字 "BlackBottomWhiteText"
	ColorShadowGreen           PlateColor = "ShadowGreen"           // 渐变绿 "ShadowGreen"
	ColorYellowGreen           PlateColor = "YellowGreen"           // 黄绿双拼 "YellowGreen"
)

const (
	PlateTypeUnknown                PlateType = "Unknown"                // "Unknown" 未知
	PlateTypeNormal                 PlateType = "Normal"                 // "Normal" 蓝牌黑牌
	PlateTypeYellow                 PlateType = "Yellow"                 // "Yellow" 黄牌
	PlateTypeDoubleyellow           PlateType = "DoubleYellow"           // "DoubleYellow" 双层黄尾牌
	PlateTypePolice                 PlateType = "Police"                 // "Police" 警牌
	PlateTypeArmed                  PlateType = "Armed"                  // "Armed" 武警牌
	PlateTypeMilitary               PlateType = "Military"               // "Military" 部队号牌
	PlateTypeDoublemilitary         PlateType = "DoubleMilitary"         // "DoubleMilitary" 部队双层
	PlateTypeSar                    PlateType = "SAR"                    // "SAR" 港澳特区号牌
	PlateTypeTrainning              PlateType = "Trainning"              // "Trainning" 教练车号牌
	PlateTypePersonal               PlateType = "Personal"               // "Personal" 个性号牌
	PlateTypeAgri                   PlateType = "Agri"                   // "Agri" 农用牌
	PlateTypeEmbassy                PlateType = "Embassy"                // "Embassy" 使馆号牌
	PlateTypeMoto                   PlateType = "Moto"                   // "Moto" 摩托车号牌
	PlateTypeTractor                PlateType = "Tractor"                // "Tractor" 拖拉机号牌
	PlateTypeOfficialcar            PlateType = "OfficialCar"            // "OfficialCar" 公务车
	PlateTypePersonalcar            PlateType = "PersonalCar"            // "PersonalCar" 私家车
	PlateTypeWarcar                 PlateType = "WarCar"                 // "WarCar"  军用
	PlateTypeOther                  PlateType = "Other"                  // "Other" 其他号牌
	PlateTypeCivilaviation          PlateType = "Civilaviation"          // "Civilaviation" 民航号牌
	PlateTypeBlack                  PlateType = "Black"                  // "Black" 黑牌
	PlateTypePurenewenergymicrocar  PlateType = "PureNewEnergyMicroCar"  // "PureNewEnergyMicroCar" 纯电动新能源小车
	PlateTypeMixednewenergymicrocar PlateType = "MixedNewEnergyMicroCar" // "MixedNewEnergyMicroCar" 混合新能源小车
	PlateTypePurenewenergylargecar  PlateType = "PureNewEnergyLargeCar"  // "PureNewEnergyLargeCar" 纯电动新能源大车
	PlateTypeMixednewenergylargecar PlateType = "MixedNewEnergyLargeCar" // "MixedNewEnergyLargeCar" 混合新能源大车
)

const (
	VehicleTypeUnknow                   VehicleType = "Unknow"                   // "Unknow" 未知类型
	VehicleTypeMotor                    VehicleType = "Motor"                    // "Motor" 机动车
	VehicleTypeNonMotor                 VehicleType = "Non-Motor"                // "Non-Motor" 非机动车
	VehicleTypeBus                      VehicleType = "Bus"                      // "Bus" 公交车
	VehicleTypeBicycle                  VehicleType = "Bicycle"                  // "Bicycle" 自行车
	VehicleTypeMotorcycle               VehicleType = "Motorcycle"               // "Motorcycle" 摩托车
	VehicleTypeUnlicensedmotor          VehicleType = "UnlicensedMotor"          // "UnlicensedMotor" 无牌机动车
	VehicleTypeLargecar                 VehicleType = "LargeCar"                 // "LargeCar"  大型汽车
	VehicleTypeMicrocar                 VehicleType = "MicroCar"                 // "MicroCar" 小型汽车
	VehicleTypeEmbassycar               VehicleType = "EmbassyCar"               // "EmbassyCar" 使馆汽车
	VehicleTypeMarginalcar              VehicleType = "MarginalCar"              // "MarginalCar" 领馆汽车
	VehicleTypeAreaoutcar               VehicleType = "AreaoutCar"               // "AreaoutCar" 境外汽车
	VehicleTypeForeigncar               VehicleType = "ForeignCar"               // "ForeignCar" 外籍汽车
	VehicleTypeDualtriwheelmotorcycle   VehicleType = "DualTriWheelMotorcycle"   // "DualTriWheelMotorcycle" 两、三轮摩托车
	VehicleTypeLightmotorcycle          VehicleType = "LightMotorcycle"          // "LightMotorcycle" 轻便摩托车
	VehicleTypeEmbassymotorcycle        VehicleType = "EmbassyMotorcycle"        // "EmbassyMotorcycle" 使馆摩托车
	VehicleTypeMarginalmotorcycle       VehicleType = "MarginalMotorcycle"       // "MarginalMotorcycle" 领馆摩托车
	VehicleTypeAreaoutmotorcycle        VehicleType = "AreaoutMotorcycle"        // "AreaoutMotorcycle" 境外摩托车
	VehicleTypeForeignmotorcycle        VehicleType = "ForeignMotorcycle"        // "ForeignMotorcycle" 外籍摩托车
	VehicleTypeFarmtransmitcar          VehicleType = "FarmTransmitCar"          // "FarmTransmitCar" 农用运输车
	VehicleTypeTractor                  VehicleType = "Tractor"                  // "Tractor" 拖拉机
	VehicleTypeTrailer                  VehicleType = "Trailer"                  // "Trailer"  挂车
	VehicleTypeCoachcar                 VehicleType = "CoachCar"                 // "CoachCar" 教练汽车
	VehicleTypeCoachmotorcycle          VehicleType = "CoachMotorcycle"          // "CoachMotorcycle" 教练摩托车
	VehicleTypeTrialcar                 VehicleType = "TrialCar"                 // "TrialCar" 试验汽车
	VehicleTypeTrialmotorcycle          VehicleType = "TrialMotorcycle"          // "TrialMotorcycle" 试验摩托车
	VehicleTypeTemporaryentrycar        VehicleType = "TemporaryEntryCar"        // "TemporaryEntryCar" 临时入境汽车
	VehicleTypeTemporaryentrymotorcycle VehicleType = "TemporaryEntryMotorcycle" // "TemporaryEntryMotorcycle" 临时入境摩托车
	VehicleTypeTemporarysteercar        VehicleType = "TemporarySteerCar"        // "TemporarySteerCar" 临时行驶车
	VehicleTypePassengercar             VehicleType = "PassengerCar"             // "PassengerCar" 客车
	VehicleTypeLargetruck               VehicleType = "LargeTruck"               // "LargeTruck" 大货车
	VehicleTypeMidtruck                 VehicleType = "MidTruck"                 // "MidTruck" 中货车
	VehicleTypeSalooncar                VehicleType = "SaloonCar"                // "SaloonCar" 轿车
	VehicleTypeMicrobus                 VehicleType = "Microbus"                 // "Microbus" 面包车
	VehicleTypeMicrotruck               VehicleType = "MicroTruck"               // "MicroTruck" 小货车
	VehicleTypeTricycle                 VehicleType = "Tricycle"                 // "Tricycle" 三轮车
	VehicleTypePasserby                 VehicleType = "Passerby"                 // "Passerby" 行人
)

const (
	VehicleColorOther  VehicleColor = "Other"  //其他颜色 "OTHER"
	VehicleColorWhite  VehicleColor = "White"  //白色	"White"
	VehicleColorBlack  VehicleColor = "Black"  //黑色	"Black"
	VehicleColorRed    VehicleColor = "Red"    //红色	"Red"
	VehicleColorYellow VehicleColor = "Yellow" //黄色	"Yellow"
	VehicleColorGray   VehicleColor = "Gray"   //灰色	"Gray"
	VehicleColorBlue   VehicleColor = "Blue"   //蓝色	"Blue"
	VehicleColorGreen  VehicleColor = "Green"  //绿色	"Green"
	VehicleColorPink   VehicleColor = "Pink"   //粉红色 "Pink"
	VehicleColorPurple VehicleColor = "Purple" //紫色	"Purple"
	VehicleColorBrown  VehicleColor = "Brown"  //棕色	"Brown"
)
