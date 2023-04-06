package app

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	app.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	app.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return app.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
