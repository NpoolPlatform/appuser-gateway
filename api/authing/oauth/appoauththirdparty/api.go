package appoauththirdparty

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/oauth/appoauththirdparty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appoauththirdparty.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	appoauththirdparty.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return appoauththirdparty.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
