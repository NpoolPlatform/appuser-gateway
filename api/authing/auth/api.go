package auth

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	auth.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	auth.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return auth.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
