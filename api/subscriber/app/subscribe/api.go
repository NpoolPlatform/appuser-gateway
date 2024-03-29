package appsubscribe

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber/app/subscribe"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	subscribe.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	subscribe.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return subscribe.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
