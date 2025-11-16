package logger

import "go.uber.org/zap"

type HTTPServerLogger struct {
	logger *zap.Logger
}

func (l *HTTPServerLogger) Write(p []byte) (int, error) {
	l.logger.Error(string(p))

	return len(p), nil
}

func NewHTTPServerLogger(l *zap.Logger) *HTTPServerLogger {
	return &HTTPServerLogger{l}
}
