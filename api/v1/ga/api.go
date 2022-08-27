package ga

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/ga"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	ga.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	ga.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return ga.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
