package ban

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/ban"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	ban.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	ban.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return ban.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
