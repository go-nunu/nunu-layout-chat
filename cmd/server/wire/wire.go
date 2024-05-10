//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"nunu-layout-chat/internal/handler"
	"nunu-layout-chat/internal/repository"
	"nunu-layout-chat/internal/server"
	"nunu-layout-chat/internal/service"
	"nunu-layout-chat/pkg/app"
	"nunu-layout-chat/pkg/jwt"
	"nunu-layout-chat/pkg/log"
	"nunu-layout-chat/pkg/server/http"
	"nunu-layout-chat/pkg/server/ws"
	"nunu-layout-chat/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewRoomService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewRoomHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
	server.NewWebSocketServer,
)

// build App
func newApp(
	httpServer *http.Server,
	job *server.Job,
	wsServer *ws.Server,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job, wsServer),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {

	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		ws.NewPitaya,
		newApp,
	))
}
