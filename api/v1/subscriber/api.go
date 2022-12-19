package subscriber

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	subscriber.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	subscriber.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return subscriber.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
