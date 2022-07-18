package api

import (
	"context"
	"github.com/NpoolPlatform/message/npool/appusergateway/app"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type AppServer struct {
	app.UnimplementedAppUserGatewayAppServer
}

func Register(server grpc.ServiceRegistrar) {
	app.RegisterAppUserGatewayAppServer(server, &AppServer{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	err := app.RegisterAppUserGatewayAppHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
	if err != nil {
		return err
	}
	return nil
}
