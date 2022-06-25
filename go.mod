module github.com/hysios/dhhttp

go 1.15

require (
	cskyzn.com/pkg/server v0.0.0-00010101000000-000000000000
	github.com/eventials/go-tus v0.0.0-20200718001131-45c7ec8f5d59
	github.com/hysios/digest v0.0.0-20201030060155-1de5ed13f2fa
	github.com/hysios/log v0.0.1
	github.com/hysios/utils v0.0.11
	github.com/kr/pretty v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/segmentio/ksuid v1.0.4
	github.com/stretchr/testify v1.7.0
)

replace github.com/hysios/log => ../log

replace cskyzn.com/pkg/server => ../../keyue/server

replace github.com/hysios/edgekv => ../edgekv

replace cskyzn.com/pkg/bimgserver => ../../keyue/bimgserver

replace cskyzn.com/pkg/flowproxy => ../../keyue/flowproxy

replace cskyzn.com/pkg/stream => ../../keyue/stream

replace cskyzn.com/pkg/utils => ../../keyue/utils
