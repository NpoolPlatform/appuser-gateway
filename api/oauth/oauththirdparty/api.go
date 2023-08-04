package oauththirdparty

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth/oauththirdparty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	oauththirdparty.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	oauththirdparty.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return oauththirdparty.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
