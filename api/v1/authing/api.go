package authing

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	authing.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	authing.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return authing.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
