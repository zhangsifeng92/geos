package reuseport_test

import (
	"fmt"
	"log"

	"github.com/zhangsifeng92/geos/plugins/http_plugin/fasthttp"
	"github.com/zhangsifeng92/geos/plugins/http_plugin/fasthttp/reuseport"
)

func ExampleListen() {
	ln, err := reuseport.Listen("tcp4", "localhost:12345")
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	if err = fasthttp.Serve(ln, requestHandler); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!")
}
