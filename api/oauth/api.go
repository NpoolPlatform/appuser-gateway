package oauth

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	oauth.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	oauth.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return oauth.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
