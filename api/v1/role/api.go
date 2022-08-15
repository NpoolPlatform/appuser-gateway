package role

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/approle"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	approle.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	approle.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return approle.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
