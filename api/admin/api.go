package admin

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
)

type Server struct {
	admin.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	admin.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return admin.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
