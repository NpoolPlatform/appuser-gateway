package app

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	app.UnimplementedAppGwServer
}

func Register(server grpc.ServiceRegistrar) {
	app.RegisterAppGwServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return app.RegisterAppGwHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
