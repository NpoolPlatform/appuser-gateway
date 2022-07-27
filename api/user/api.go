package user

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/appuser"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appuser.UnimplementedAppUserGatewayAppUserServer
}

func Register(server grpc.ServiceRegistrar) {
	appuser.RegisterAppUserGatewayAppUserServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return appuser.RegisterAppUserGatewayAppUserHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
