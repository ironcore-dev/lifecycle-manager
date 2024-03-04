// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"log/slog"

	protov "github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

func UnaryServerValidatorInterceptor(val *protov.Validator) grpc.UnaryServerInterceptor {
	return protovalidate.UnaryServerInterceptor(val)
}

func UnaryServerLoggerInterceptor(log *slog.Logger, opts ...logging.Option) grpc.UnaryServerInterceptor {
	l := interceptLogger(log)
	return logging.UnaryServerInterceptor(l, opts...)
}

func StreamServerValidatorInterceptor(val *protov.Validator) grpc.StreamServerInterceptor {
	return protovalidate.StreamServerInterceptor(val)
}

func StreamServerLoggingInterceptor(log *slog.Logger, opts ...logging.Option) grpc.StreamServerInterceptor {
	l := interceptLogger(log)
	return logging.StreamServerInterceptor(l, opts...)
}

func interceptLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, level logging.Level, msg string, fields ...any) {
		switch level {
		case logging.LevelDebug:
			l.Debug(msg, fields...)
		case logging.LevelInfo:
			l.Info(msg, fields...)
		case logging.LevelWarn:
			l.Warn(msg, fields...)
		case logging.LevelError:
			l.Error(msg, fields...)
		}
	})
}
