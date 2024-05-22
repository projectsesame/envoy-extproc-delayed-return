package main

import (
	"log"
	"strconv"
	"time"

	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	ep "github.com/wrossmorrow/envoy-extproc-sdk-go"
)

type delayedReturnRequestProcessor struct {
	opts *ep.ProcessingOptions
}

func (s *delayedReturnRequestProcessor) GetName() string {
	return "delayed-return"
}

func (s *delayedReturnRequestProcessor) GetOptions() *ep.ProcessingOptions {
	return s.opts
}

const kDelayedTime = "delayed-time"

func (s *delayedReturnRequestProcessor) ProcessRequestHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {
	cancel := func(code int32, msg string) error {
		return ctx.CancelRequest(code, map[string]ep.HeaderValue{}, typev3.StatusCode_name[code])
	}

	var delayedTime int64
	var err error

	if headers.RawHeaders[kDelayedTime] == nil {
		delayedTime = 100
		log.Printf("[%s], delayed-time header is not set, use default value 100", ctx.RequestID)

	} else {
		delayedTime, err = strconv.ParseInt(string(headers.RawHeaders[kDelayedTime]), 10, 64)
		if err != nil {
			return cancel(400, "delayed-time header is invalid")
		}
	}

	if delayedTime < 0 {
		return cancel(400, "delayed-time header is invalid")
	}

	ctx.SetValue("delayed-time", delayedTime)

	logDelayedTime(ctx, "ProcessRequestHeaders")

	return ctx.ContinueRequest()
}

func (s *delayedReturnRequestProcessor) ProcessRequestBody(ctx *ep.RequestContext, body []byte) error {

	logDelayedTime(ctx, "ProcessRequestBody")

	return ctx.ContinueRequest()
}

func logDelayedTime(ctx *ep.RequestContext, method string) {

	var delayedTime int64

	v, _ := ctx.GetValue(`delayed-time`)

	delayedTime = v.(int64)

	now := time.Now()
	time.Sleep(time.Duration(delayedTime) * time.Millisecond)

	log.Printf("[%s],process %s, delayed time: %d, elapsed time: %d\n", ctx.RequestID, method, delayedTime, time.Since(now).Milliseconds())

}

func (s *delayedReturnRequestProcessor) ProcessRequestTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {

	logDelayedTime(ctx, "ProcessRequestTrailers")

	return ctx.ContinueRequest()
}

func (s *delayedReturnRequestProcessor) ProcessResponseHeaders(ctx *ep.RequestContext, headers ep.AllHeaders) error {

	logDelayedTime(ctx, "ProcessResponseHeaders")
	return ctx.ContinueRequest()
}

func (s *delayedReturnRequestProcessor) ProcessResponseBody(ctx *ep.RequestContext, body []byte) error {

	logDelayedTime(ctx, "ProcessResponseBody")
	return ctx.ContinueRequest()
}

func (s *delayedReturnRequestProcessor) ProcessResponseTrailers(ctx *ep.RequestContext, trailers ep.AllHeaders) error {

	logDelayedTime(ctx, "ProcessResponseTrailers")
	return ctx.ContinueRequest()
}

func (s *delayedReturnRequestProcessor) Init(opts *ep.ProcessingOptions, nonFlagArgs []string) error {
	s.opts = opts
	return nil
}

func (s *delayedReturnRequestProcessor) Finish() {}
