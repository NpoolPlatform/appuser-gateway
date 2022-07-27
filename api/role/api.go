package role

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/approle"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	approle.UnimplementedAppUserGatewayAppRoleServer
}

func Register(server grpc.ServiceRegistrar) {
	approle.RegisterAppUserGatewayAppRoleServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return approle.RegisterAppUserGatewayAppRoleHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
