package etcd

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type Option func(*P)

// WithVersion sets image version.
func WithVersion(version string) Option {
	return func(o *P) {
		o.Version = version
	}
}

// WithAutoSyncInterval is the interval to update endpoints with its latest members.
// 0 disables auto-sync. By default auto-sync is disabled.
func WithAutoSyncInterval(autoSyncInterval time.Duration) Option {
	return func(o *P) {
		o.AutoSyncInterval = autoSyncInterval
	}
}

// WithDialTimeout is the timeout for failing to establish a connection.
func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *P) {
		o.DialTimeout = dialTimeout
	}
}

// WithDialKeepAliveTime is the time after which client pings the server to see if
// transport is alive.
func WithDialKeepAliveTime(dialKeepAliveTime time.Duration) Option {
	return func(o *P) {
		o.DialKeepAliveTime = dialKeepAliveTime
	}
}

// WithDialKeepAliveTimeout is the time that the client waits for a response for the
// keep-alive probe. If the response is not received in this time, the connection is closed.
func WithDialKeepAliveTimeout(dialKeepAliveTimeout time.Duration) Option {
	return func(o *P) {
		o.DialKeepAliveTimeout = dialKeepAliveTimeout
	}
}

// WithDialOptions is a list of dial options for the grpc client (e.g., for interceptors).
// For example, pass "grpc.WithBlock()" to block until the underlying connection is up.
// Without this, Dial returns immediately and connecting the server happens in background.
func WithDialOptions(dialOptions []grpc.DialOption) Option {
	return func(o *P) {
		o.DialOptions = dialOptions
	}
}

// WithLogConfig configures client-side logger.
// If nil, use the default logger.
// TODO: configure gRPC logger
func WithLogConfig(logConfig *zap.Config) Option {
	return func(o *P) {
		o.LogConfig = logConfig
	}
}
