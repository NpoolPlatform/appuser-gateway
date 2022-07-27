package banappuser

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/banappuser"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	banappuser.UnimplementedAppUserGatewayBanAppUserServer
}

func Register(server grpc.ServiceRegistrar) {
	banappuser.RegisterAppUserGatewayBanAppUserServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return banappuser.RegisterAppUserGatewayBanAppUserHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
