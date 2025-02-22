module github.com/zhangsifeng92/geos

go 1.13

require (
	github.com/eapache/channels v1.1.0
	github.com/fatih/color v1.12.0
	github.com/go-stack/stack v1.8.0
	github.com/peterh/liner v1.2.1
	github.com/robertkrimen/otto v0.0.0-20210614181706-373ff5438452
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/tidwall/gjson v1.8.1
	github.com/urfave/cli v1.22.5
	github.com/zhangsifeng92/geos/plugins/http_plugin/fasthttp v0.0.0
	github.com/zhangsifeng92/geos/wasmgo/wagon v0.0.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

replace (
	github.com/zhangsifeng92/geos/plugins/http_plugin/fasthttp => ./plugins/http_plugin/fasthttp
	github.com/zhangsifeng92/geos/wasmgo/wagon => ./wasmgo/wagon
)
