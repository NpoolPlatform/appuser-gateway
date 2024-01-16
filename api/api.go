package api

import (
	"context"

	"github.com/NpoolPlatform/appuser-gateway/api/admin"
	"github.com/NpoolPlatform/appuser-gateway/api/app"
	"github.com/NpoolPlatform/appuser-gateway/api/authing/auth"
	authhistory "github.com/NpoolPlatform/appuser-gateway/api/authing/history"
	"github.com/NpoolPlatform/appuser-gateway/api/ga"
	"github.com/NpoolPlatform/appuser-gateway/api/kyc"
	"github.com/NpoolPlatform/appuser-gateway/api/role"
	roleuser "github.com/NpoolPlatform/appuser-gateway/api/role/user"
	"github.com/NpoolPlatform/appuser-gateway/api/subscriber"
	appsubscribe "github.com/NpoolPlatform/appuser-gateway/api/subscriber/app/subscribe"
	"github.com/NpoolPlatform/appuser-gateway/api/user"
	"github.com/NpoolPlatform/appuser-gateway/api/user/recoverycode"

	"github.com/NpoolPlatform/appuser-gateway/api/oauth"
	"github.com/NpoolPlatform/appuser-gateway/api/oauth/appoauththirdparty"
	"github.com/NpoolPlatform/appuser-gateway/api/oauth/oauththirdparty"
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
	subscriber.Register(server)
	appsubscribe.Register(server)
	role.Register(server)
	roleuser.Register(server)
	user.Register(server)
	recoverycode.Register(server)
	ga.Register(server)
	auth.Register(server)
	authhistory.Register(server)
	kyc.Register(server)
	oauth.Register(server)
	oauththirdparty.Register(server)
	appoauththirdparty.Register(server)
}

//nolint:gocyclo
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
	if err := recoverycode.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := subscriber.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appsubscribe.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := role.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := roleuser.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := user.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := ga.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := auth.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := authhistory.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := kyc.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := oauth.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := oauththirdparty.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appoauththirdparty.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}

	return nil
}
