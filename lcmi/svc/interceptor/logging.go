// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package interceptor

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
)

type LoggerInterceptor struct {
	logger *slog.Logger
}

func NewLoggerInterceptor(log *slog.Logger) connect.Interceptor {
	return &LoggerInterceptor{logger: log}
}

func (l *LoggerInterceptor) WrapUnary(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		l.logger.Info(req.Spec().Procedure, "request", req.Any())
		reqCtx := logr.NewContextWithSlogLogger(ctx, l.logger)
		response, err := unaryFunc(reqCtx, req)
		if err != nil {
			l.logger.Error(err.Error())
		}
		return response, err
	}
}

func (l *LoggerInterceptor) WrapStreamingClient(clientFunc connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return &streamingClientInterceptor{
			StreamingClientConn: clientFunc(ctx, spec),
			logger:              l.logger,
		}
	}
}

func (l *LoggerInterceptor) WrapStreamingHandler(handlerFunc connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return handlerFunc(ctx, &streamingHandlerInterceptor{
			StreamingHandlerConn: conn,
			logger:               l.logger,
		})
	}
}

type streamingClientInterceptor struct {
	connect.StreamingClientConn
	logger *slog.Logger
}

func (s *streamingClientInterceptor) Spec() connect.Spec {
	// TODO implement me
	panic("implement me")
}

func (s *streamingClientInterceptor) Peer() connect.Peer {
	// TODO implement me
	panic("implement me")
}

// func (s *streamingClientInterceptor) Send(a any) error {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingClientInterceptor) RequestHeader() http.Header {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingClientInterceptor) CloseRequest() error {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingClientInterceptor) Receive(a any) error {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingClientInterceptor) ResponseHeader() http.Header {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingClientInterceptor) ResponseTrailer() http.Header {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingClientInterceptor) CloseResponse() error {
// 	// TODO implement me
// 	panic("implement me")
// }

type streamingHandlerInterceptor struct {
	connect.StreamingHandlerConn
	logger *slog.Logger
}

// func (s *streamingHandlerInterceptor) Spec() connect.Spec {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingHandlerInterceptor) Peer() connect.Peer {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingHandlerInterceptor) Receive(a any) error {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingHandlerInterceptor) RequestHeader() http.Header {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingHandlerInterceptor) Send(a any) error {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (s *streamingHandlerInterceptor) ResponseHeader() http.Header {
// 	// TODO implement me
// 	panic("implement me")
// }

func (s *streamingHandlerInterceptor) ResponseTrailer() http.Header {
	// TODO implement me
	panic("implement me")
}
