package extra

import (
	"context"
	"github.com/NpoolPlatform/message/npool/appusergw/appuserextra"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appuserextra.UnimplementedAppUserGatewayExtraServer
}

func Register(server grpc.ServiceRegistrar) {
	appuserextra.RegisterAppUserGatewayExtraServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return appuserextra.RegisterAppUserGatewayExtraHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
