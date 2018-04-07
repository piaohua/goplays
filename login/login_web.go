package main

import (
	"fmt"
	"time"

	"goplays/pb"

	"github.com/valyala/fasthttp"
)

// web
func webHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case "POST":
	default:
		fmt.Fprintf(ctx, "%s", "failed")
		return
	}
	switch getIP(ctx) {
	case "127.0.0.1":
	default:
		fmt.Fprintf(ctx, "%s", "failed")
		return
	}
	result := ctx.PostBody()
	msg1 := new(pb.WebRequest)
	err1 := msg1.Unmarshal(result)
	if err1 != nil {
		fmt.Fprintf(ctx, "%v", err1)
		return
	}
	timeout := 3 * time.Second
	res2, err2 := nodePid.RequestFuture(msg1, timeout).Result()
	if err2 != nil {
		fmt.Fprintf(ctx, "%v", err2)
		return
	}
	var response2 *pb.WebResponse
	var ok bool
	if response2, ok = res2.(*pb.WebResponse); !ok {
		fmt.Fprintf(ctx, "%s", "failed")
		return
	}
	body, err1 := response2.Marshal()
	if err1 != nil {
		fmt.Fprintf(ctx, "%v", err1)
		return
	}
	fmt.Fprintf(ctx, "%s", body)
}
