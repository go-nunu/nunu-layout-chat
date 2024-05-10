package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"
	"nunu-layout-chat/internal/handler"
	"nunu-layout-chat/pkg/log"
	"nunu-layout-chat/pkg/server/ws"
	"strings"
)

type WebSocketServer struct {
	logger *log.Logger
	conf   *viper.Viper
}

func NewWebSocketServer(
	logger *log.Logger,
	roomHandler *handler.RoomHandler,
	conf *viper.Viper,
	app pitaya.Pitaya,
) *ws.Server {
	pLogger := log.NewPitayaLog(conf)
	pitaya.SetLogger(pLogger)
	s := ws.NewServer(
		gin.Default(),
		logger,
		ws.WithPitayaApp(app),
		ws.WithServerHost(conf.GetString("ws.host")),
		ws.WithServerPort(conf.GetInt("ws.port")),
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
