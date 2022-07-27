package banapp

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/banapp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	banapp.UnimplementedAppUserGatewayBanAppServer
}

func Register(server grpc.ServiceRegistrar) {
	banapp.RegisterAppUserGatewayBanAppServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return banapp.RegisterAppUserGatewayBanAppHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
