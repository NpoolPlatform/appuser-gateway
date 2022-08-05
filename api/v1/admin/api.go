package admin

import (
	"context"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	admin.UnimplementedAdminGwServer
}

func Register(server grpc.ServiceRegistrar) {
	admin.RegisterAdminGwServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return admin.RegisterAdminGwHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
