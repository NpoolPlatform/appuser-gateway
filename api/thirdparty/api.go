package thirdparty

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appusergw/thirdparty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	thirdparty.UnimplementedAppUserGatewayThirdPartyServer
}

func Register(server grpc.ServiceRegistrar) {
	thirdparty.RegisterAppUserGatewayThirdPartyServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return thirdparty.RegisterAppUserGatewayThirdPartyHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
