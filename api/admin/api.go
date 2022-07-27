package admin

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/admin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	admin.UnimplementedAppUserGatewayAdminServer
}

func Register(server grpc.ServiceRegistrar) {
	admin.RegisterAppUserGatewayAdminServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return admin.RegisterAppUserGatewayAdminHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
