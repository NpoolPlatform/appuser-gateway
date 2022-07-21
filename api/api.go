package api

import (
	"context"
	"github.com/NpoolPlatform/appuser-gateway/api/appuser"
	"github.com/NpoolPlatform/appuser-gateway/api/appusersecret"

	"github.com/NpoolPlatform/appuser-gateway/api/app"
	"github.com/NpoolPlatform/appuser-gateway/api/banapp"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appusergw.UnimplementedAppUserGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	appusergw.RegisterAppUserGatewayServer(server, &Server{})
	app.Register(server)
	banapp.Register(server)
	appuser.Register(server)
	appusersecret.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := appusergw.RegisterAppUserGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := app.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := banapp.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appuser.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appusersecret.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
