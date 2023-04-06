package role

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/role"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	role.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	role.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return role.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
