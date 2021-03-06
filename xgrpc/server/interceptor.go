/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2021/1/2 12:50
 */
package serverinterceptors

import (
	"bytes"
	"context"
	"fmt"
	"github.com/coder2m/component/xcode"
	"github.com/coder2m/g-saber/xlog"
	"github.com/coder2m/component/xmonitor"
	"github.com/coder2m/component/xtrace"
	"github.com/coder2m/g-saber/xcast"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"runtime/debug"
	"strings"
	"time"
)

func extractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	var buf bytes.Buffer
	stack := debug.Stack()
	buf.Write(stack)
	xlog.Error(fmt.Sprintf("%+v", r), xlog.FieldValue(buf.String()))
	return xcode.SystemCodeAdd(xcast.ToUint32(codes.Internal), "server internal error")
}

func CrashUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		return handler(ctx, req)
	}
}

func CrashStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) (err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})
		return handler(srv, stream)
	}
}

func PrometheusUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		code := xcode.ExtractCodes(err)
		xmonitor.ServerHandleHistogram.WithLabelValues(xmonitor.TypeGRPCUnary, info.FullMethod, extractAID(ctx)).Observe(time.Since(startTime).Seconds())
		xmonitor.ServerHandleCounter.WithLabelValues(xmonitor.TypeGRPCUnary, info.FullMethod, extractAID(ctx), xcast.ToString(code.GetCode())).Inc()
		if code != xcode.OK {
			xmonitor.ServerErrorCounter.WithLabelValues(xmonitor.TypeGRPCUnary, info.FullMethod, extractAID(ctx), xcast.ToString(code.GetCode())).Inc()
		}
		return resp, err
	}
}

func PrometheusStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		err := handler(srv, ss)
		code := xcode.ExtractCodes(err)
		xmonitor.ServerHandleHistogram.WithLabelValues(xmonitor.TypeGRPCStream, info.FullMethod, extractAID(ss.Context())).Observe(time.Since(startTime).Seconds())
		xmonitor.ServerHandleCounter.WithLabelValues(xmonitor.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()), xcast.ToString(code.GetCode())).Inc()
		if code != xcode.OK {
			xmonitor.ServerErrorCounter.WithLabelValues(xmonitor.TypeGRPCUnary, info.FullMethod, extractAID(ss.Context()), xcast.ToString(code.GetCode())).Inc()
		}
		return err
	}
}

func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		span, ctx := xtrace.StartSpanFromContext(
			ctx,
			info.FullMethod,
			xtrace.FromIncomingContext(ctx),
			xtrace.TagComponent("gRPC"),
			xtrace.TagSpanKind("server.unary"),
		)
		defer span.Finish()
		resp, err := handler(ctx, req)
		if err != nil {
			spbStatus := xcode.ExtractCodes(err)
			span.SetTag("code", spbStatus.GetCode())
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		return resp, err
	}
}

func XTimeoutUnaryServerInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		if deadline, ok := ctx.Deadline(); ok {
			leftTime := time.Until(deadline)
			if leftTime < timeout {
				timeout = leftTime
			}
		}
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
		defer cancel()
		return handler(ctx, req)
	}
}
