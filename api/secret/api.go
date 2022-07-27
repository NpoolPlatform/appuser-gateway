package secret

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/appusersecret"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appusersecret.UnimplementedAppUserGatewaySecretServer
}

func Register(server grpc.ServiceRegistrar) {
	appusersecret.RegisterAppUserGatewaySecretServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return appusersecret.RegisterAppUserGatewaySecretHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
