module github.com/hysios/dhhttp

go 1.15

require (
	cskyzn.com/pkg/server v0.0.0-00010101000000-000000000000 // indirect
	github.com/eventials/go-tus v0.0.0-20200718001131-45c7ec8f5d59
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/googleapis/gax-go v1.0.3 // indirect
	github.com/hysios/digest v0.0.0-20201030060155-1de5ed13f2fa
	github.com/hysios/log v0.0.0-20210420091742-d54e2f0555dd
	github.com/hysios/utils v0.0.10
	github.com/kr/pretty v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/segmentio/ksuid v1.0.3
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.14.1
	google.golang.org/api v0.44.0 // indirect
	gotest.tools v2.2.0+incompatible
)

replace github.com/hysios/log => ../log

replace cskyzn.com/pkg/server => ../../keyue/server

replace github.com/hysios/edgekv => ../edgekv
