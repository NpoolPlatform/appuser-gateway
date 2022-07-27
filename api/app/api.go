package app

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/app"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	app.UnimplementedAppUserGatewayAppServer
}

func Register(server grpc.ServiceRegistrar) {
	app.RegisterAppUserGatewayAppServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return app.RegisterAppUserGatewayAppHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
