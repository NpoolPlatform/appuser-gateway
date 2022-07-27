package api

import (
	"context"

	"github.com/NpoolPlatform/appuser-gateway/api/admin"
	"github.com/NpoolPlatform/appuser-gateway/api/banappuser"
	"github.com/NpoolPlatform/appuser-gateway/api/extra"
	"github.com/NpoolPlatform/appuser-gateway/api/role"
	"github.com/NpoolPlatform/appuser-gateway/api/roleuser"
	"github.com/NpoolPlatform/appuser-gateway/api/secret"
	"github.com/NpoolPlatform/appuser-gateway/api/thirdparty"
	"github.com/NpoolPlatform/appuser-gateway/api/user"

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
	banappuser.Register(server)
	extra.Register(server)
	role.Register(server)
	roleuser.Register(server)
	secret.Register(server)
	user.Register(server)
	admin.Register(server)
	thirdparty.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := appusergw.RegisterAppUserGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := admin.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := app.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := banapp.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := banappuser.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := extra.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := role.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := roleuser.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := secret.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := user.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := thirdparty.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}

	return nil
}
