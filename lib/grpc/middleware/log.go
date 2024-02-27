package middleware

import (
	"context"
	"path"
	"path/filepath"
	"time"

	"triones-one/lib/zlog"

	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	logger *zap.Logger
	LogDir = "../log"
)

func InitLogger() error {
	access, _, err := zlog.InitLogger(&zlog.Config{
		Level:            "info",
		Format:           "json",
		DisableTimestamp: false,
		File: zlog.FileLogConfig{
			Filename:   filepath.Join(LogDir, "access.log"),
			MaxSize:    100,
			MaxDays:    180,
			MaxBackups: 90,
			Compress:   true,
		},
		DisableStacktrace:   true,
		DisableErrorVerbose: true,
	})
	if err != nil {
		return err
	}
	logger = access
	return nil
}

// UnaryServerLogInterceptor log 拦截
func UnaryServerLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var (
		addr    string
		traceId string
	)
	if peer, ok := peer.FromContext(ctx); ok {
		addr = peer.Addr.String()
	}

	traceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()

	start := time.Now()
	defer func() {
		logger.Info("",
			zap.Error(err),
			zap.String("traceId", traceId),
			zap.String("addr", addr),
			zap.Duration("cost", time.Since(start)),
			zap.String("service", path.Dir(info.FullMethod)[1:]),
			zap.String("method", path.Base(info.FullMethod)),
			zap.String("req", jsonFormat(req)),
			zap.String("resp", jsonFormat(resp)),
		)
	}()

	resp, err = handler(ctx, req)
	return resp, err
}

// StreamServerRecoveryInterceptor recover
func StreamServerLogInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	var (
		addr    string
		traceId string
		ctx     = stream.Context()
	)
	if peer, ok := peer.FromContext(ctx); ok {
		addr = peer.Addr.String()
	}

	traceId = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	start := time.Now()
	defer func() {
		logger.Info("",
			zap.Error(err),
			zap.String("traceId", traceId),
			zap.String("addr", addr),
			zap.Duration("cost", time.Since(start)),
			zap.String("service", path.Dir(info.FullMethod)[1:]),
			zap.String("method", path.Base(info.FullMethod)),
			zap.String("srv", jsonFormat(srv)),
		)
	}()

	return handler(srv, stream)
}

func jsonFormat(data interface{}) string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(data)
	if err == nil {
		return string(buf)
	}

	if a, ok := data.(StringAble); ok {
		return a.String()
	}

	return ""
}

// StringAble string
type StringAble interface {
	String() string
}
