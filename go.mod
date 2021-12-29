module dev.cspdls.com/pkg/dhhttp

go 1.15

require (
	dev.cspdls.com/pkg/log v0.0.0-00010101000000-000000000000
	dev.cspdls.com/pkg/storage v0.0.0-00010101000000-000000000000
	github.com/eventials/go-tus v0.0.0-20200718001131-45c7ec8f5d59
	github.com/googleapis/gax-go v1.0.3 // indirect
	github.com/hysios/digest v0.0.0-20201030060155-1de5ed13f2fa
	github.com/hysios/utils v0.0.0-20210311052513-07f534619a64
	github.com/kr/pretty v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/segmentio/ksuid v1.0.3
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.14.1
	google.golang.org/api v0.44.0 // indirect
	gotest.tools v2.2.0+incompatible
)

replace dev.cspdls.com/pkg/log => ../log

replace dev.cspdls.com/pkg/storage => ../storage
