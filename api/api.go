package api

import (
	"context"

	"github.com/NpoolPlatform/appuser-gateway/api/v1/admin"
	"github.com/NpoolPlatform/appuser-gateway/api/v1/app"
	"github.com/NpoolPlatform/appuser-gateway/api/v1/authing"
	"github.com/NpoolPlatform/appuser-gateway/api/v1/ban"
	"github.com/NpoolPlatform/appuser-gateway/api/v1/role"
	"github.com/NpoolPlatform/appuser-gateway/api/v1/user"

	appusergw "github.com/NpoolPlatform/message/npool/appuser/gw/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appusergw.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	appusergw.RegisterGatewayServer(server, &Server{})
	admin.Register(server)
	app.Register(server)
	ban.Register(server)
	role.Register(server)
	user.Register(server)
	authing.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := appusergw.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := admin.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := app.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := ban.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := role.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := user.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := authing.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}

	return nil
}
