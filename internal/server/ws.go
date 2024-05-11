package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/acceptor"
	"github.com/topfreegames/pitaya/v2/component"
	"github.com/topfreegames/pitaya/v2/config"
	"github.com/topfreegames/pitaya/v2/constants"
	"github.com/topfreegames/pitaya/v2/groups"
	"go.uber.org/zap"
	"nunu-layout-chat/internal/handler"
	"nunu-layout-chat/pkg/log"
	"nunu-layout-chat/pkg/server/ws"
	"strings"
	"time"
)

type WebSocketServer struct {
	logger *log.Logger
	conf   *viper.Viper
}

func NewWebSocketServer(
	logger *log.Logger,
	conf *viper.Viper,
	app pitaya.Pitaya,
	roomHandler *handler.RoomHandler,
) *ws.Server {
	pLogger := log.NewPitayaLog(conf)
	pitaya.SetLogger(pLogger)
	s := ws.NewServer(
		logger,
		ws.WithPitayaApp(app),
	)
	// rewrite component and handler name
	app.Register(roomHandler,
		component.WithName("room"),
		component.WithNameFunc(strings.ToLower),
	)

	// more handler
	//app.Register(voiceHandler,
	//	component.WithName("voice"),
	//	component.WithNameFunc(strings.ToLower),
	//)
	return s
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
			return ctx, nil, errors.New("unauthorized")
		}
		if in != nil {
		}
		return ctx, in, nil
	})
	builder.HandlerHooks.AfterHandler.PushBack(func(ctx context.Context, out interface{}, err error) (interface{}, error) {
		logger.WithContext(ctx).Info("Response", zap.Any("data", out))
		return out, err
	})

	builder.AddAcceptor(acceptor.NewWSAcceptor(fmt.Sprintf("%s:%d", conf.GetString("ws.host"), conf.GetInt("ws.port"))))
	//builder.AddAcceptor(acceptor.NewTCPAcceptor(""))
	builder.Groups = groups.NewMemoryGroupService(builder.Config.Groups.Memory)
	return builder.Build()
}
