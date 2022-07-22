package roleuser

import (
	"context"
	"github.com/NpoolPlatform/message/npool/appusergw/approleuser"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	approleuser.UnimplementedAppUserGatewayAppRoleUserServer
}

func Register(server grpc.ServiceRegistrar) {
	approleuser.RegisterAppUserGatewayAppRoleUserServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return approleuser.RegisterAppUserGatewayAppRoleUserHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
