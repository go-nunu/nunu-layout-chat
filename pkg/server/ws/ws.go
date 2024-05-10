package ws

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/acceptor"
	"github.com/topfreegames/pitaya/v2/config"
	"github.com/topfreegames/pitaya/v2/constants"
	"github.com/topfreegames/pitaya/v2/groups"
	"go.uber.org/zap"
	"net/http"
	"nunu-layout-chat/pkg/log"
	"time"
)

type Server struct {
	*gin.Engine
	httpSrv *http.Server
	host    string
	port    int
	logger  *log.Logger
	app     pitaya.Pitaya
}
type Option func(s *Server)

func NewServer(engine *gin.Engine, logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		Engine: engine,
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}
func WithPitayaApp(app pitaya.Pitaya) Option {
	return func(s *Server) {
		s.app = app
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.app.Start()
	return nil
}
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Sugar().Info("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.app.Shutdown()

	s.logger.Sugar().Info("Server exiting")
	return nil
}

func NewPitaya(logger *log.Logger, conf *viper.Viper) pitaya.Pitaya {
	pitayaConf := config.NewDefaultPitayaConfig()
	pitayaConf.Buffer.Handler.LocalProcess = 15
	pitayaConf.Heartbeat.Interval = 15 * time.Second
	pitayaConf.Buffer.Agent.Messages = 32
	pitayaConf.Handler.Messages.Compression = false
	builder := pitaya.NewDefaultBuilder(true, "chat", pitaya.Standalone, map[string]string{}, *pitayaConf)

	// ws middleware
	var exceptInRoute = []string{
		"room.create",
		"room.join",
	}
	builder.HandlerHooks.BeforeHandler.PushBack(func(ctx context.Context, in interface{}) (c context.Context, out interface{}, err error) {
		route := pitaya.GetFromPropagateCtx(ctx, constants.RouteKey)
		requestID := pitaya.GetFromPropagateCtx(ctx, constants.RequestIDKey)
		session := pitaya.GetSessionFromCtx(ctx)
		ctx = logger.WithValue(ctx, zap.String("trace", requestID.(string)))
		ctx = logger.WithValue(ctx, zap.String("route", route.(string)))
		ctx = logger.WithValue(ctx, zap.String("userId", session.UID()))
		logger.WithContext(ctx).Info("Request")
		for _, v := range exceptInRoute {
			if route == v {
				return ctx, in, nil
			}
		}
		if session.UID() == "" {
			return ctx, nil, errors.New("Unauthorized")
		}
		if in != nil {
		}
		return ctx, in, nil
	})
	builder.HandlerHooks.AfterHandler.PushBack(func(ctx context.Context, out interface{}, err error) (interface{}, error) {
		logger.WithContext(ctx).Info("Response", zap.Any("data", out))
		return out, err
	})

	builder.AddAcceptor(acceptor.NewWSAcceptor(":" + conf.GetString("ws.port")))
	builder.Groups = groups.NewMemoryGroupService(builder.Config.Groups.Memory)
	return builder.Build()
}
