package recoverycode

import (
	"context"

	recoverycode1 "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user/recoverycode"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	recoverycode1.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	recoverycode1.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return recoverycode1.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
