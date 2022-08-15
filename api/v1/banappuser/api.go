package banappuser

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/banappuser"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	banappuser.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	banappuser.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return banappuser.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
