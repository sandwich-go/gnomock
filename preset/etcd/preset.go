package etcd

import (
	"context"
	"fmt"
	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/internal/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/pkg/logutil"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

const (
	defaultPort    = 2379
	defaultVersion = "3.4.24"
)

var setLoggerOnce sync.Once

func init() {
	registry.Register("etcd", func() gnomock.Preset { return &P{} })
}

type P struct {
	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	AutoSyncInterval time.Duration `json:"auto-sync-interval"`
	// DialTimeout is the timeout for failing to establish a connection.
	DialTimeout time.Duration `json:"dial-timeout"`

	// DialKeepAliveTime is the time after which client pings the server to see if
	// transport is alive.
	DialKeepAliveTime time.Duration `json:"dial-keep-alive-time"`

	// DialKeepAliveTimeout is the time that the client waits for a response for the
	// keep-alive probe. If the response is not received in this time, the connection is closed.
	DialKeepAliveTimeout time.Duration `json:"dial-keep-alive-timeout"`
	// Username is a user name for authentication.
	Username string `json:"username"`

	// Password is a password for authentication.
	Password string `json:"password"`
	// DialOptions is a list of dial options for the grpc client (e.g., for interceptors).
	// For example, pass "grpc.WithBlock()" to block until the underlying connection is up.
	// Without this, Dial returns immediately and connecting the server happens in background.
	DialOptions []grpc.DialOption
	// LogConfig configures client-side logger.
	// If nil, use the default logger.
	// TODO: configure gRPC logger
	LogConfig *zap.Config

	Version string `json:"version"`
}

func Preset(opts ...Option) gnomock.Preset {
	p := &P{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *P) setDefaults() {
	if p.DialOptions == nil {
		p.DialOptions = []grpc.DialOption{grpc.WithBlock()}
	}
	if p.LogConfig == nil {
		zapLoggerConfig := logutil.DefaultZapLoggerConfig
		if zapLoggerConfig.InitialFields == nil {
			zapLoggerConfig.InitialFields = make(map[string]interface{})
		}
	}
	if p.Version == "" {
		p.Version = defaultVersion
	}
}

// Image returns an image that should be pulled to create this container.
func (p *P) Image() string {
	return fmt.Sprintf("quay.io/coreos/etcd:v%s", p.Version)
}

// Ports returns ports that should be used to access this container.
func (p *P) Ports() gnomock.NamedPorts {
	return gnomock.DefaultTCP(defaultPort)
}

// Options returns a list of options to configure this container.
func (p *P) Options() []gnomock.Option {
	setLoggerOnce.Do(func() {
		// err is always nil for non-nil logger
		_ = mysqldriver.SetLogger(log.New(ioutil.Discard, "", -1))
	})

	p.setDefaults()

	opts := []gnomock.Option{
		gnomock.WithEnv("ETCD_ADVERTISE_CLIENT_URLS=http://192.168.1.21:2379"),
		gnomock.WithEnv("ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379"),
		gnomock.WithHealthCheck(p.healthcheck),
		gnomock.WithInit(p.initf()),
	}
	return opts
}

func (p *P) initf() gnomock.InitFunc {
	return func(ctx context.Context, c *gnomock.Container) error {
		addr := c.Address(gnomock.DefaultPort)
		db, err := p.connect(ctx, addr)
		if err != nil {
			return err
		}
		defer func() { _ = db.Close() }()
		return err
	}
}

func (p *P) healthcheck(ctx context.Context, c *gnomock.Container) error {
	addr := c.Address(gnomock.DefaultPort)
	db, err := p.connect(ctx, addr)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()
	return err
}

func (p *P) connect(ctx context.Context, endpoints ...string) (*clientv3.Client, error) {
	v3cc := clientv3.Config{
		Endpoints:            endpoints,
		AutoSyncInterval:     p.AutoSyncInterval,
		DialTimeout:          p.DialTimeout,
		DialKeepAliveTime:    p.DialKeepAliveTime,
		DialKeepAliveTimeout: p.DialKeepAliveTimeout,
		DialOptions:          p.DialOptions,
		LogConfig:            p.LogConfig,
	}
	cc, err := clientv3.New(v3cc)
	if err != nil {
		return nil, err
	}
	_, err = cc.Status(ctx, endpoints[0])
	if err != nil {
		return nil, err
	}
	return cc, err
}
