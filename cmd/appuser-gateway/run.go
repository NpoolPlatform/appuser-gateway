package main

import (
	"context"

	"github.com/NpoolPlatform/appuser-gateway/pkg/migrator"

	"github.com/NpoolPlatform/appuser-gateway/api"
	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
	ossconst "github.com/NpoolPlatform/go-service-framework/pkg/oss/const"

	apicli "github.com/NpoolPlatform/basal-middleware/pkg/client/api"

	cli "github.com/urfave/cli/v2"

	"github.com/NpoolPlatform/go-service-framework/pkg/action"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const BukectKey = "kyc_bucket"

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"s"},
	Usage:   "Run the daemon",
	Action: func(c *cli.Context) error {
		return action.Run(
			c.Context,
			run,
			rpcRegister,
			rpcGatewayRegister,
			watch,
		)
	},
}

func run(ctx context.Context) error {
	if err := migrator.Migrate(ctx); err != nil {
		return err
	}

	if err := oss.Init(ossconst.SecretStoreKey, BukectKey); err != nil {
		return err
	}
	return nil
}

func watch(ctx context.Context, cancel context.CancelFunc) error {
	return nil
}

func rpcRegister(server grpc.ServiceRegistrar) error {
	api.Register(server)

	apicli.RegisterGRPC(server)

	return nil
}

func rpcGatewayRegister(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	err := api.RegisterGateway(mux, endpoint, opts)
	if err != nil {
		return err
	}

	_ = apicli.Register(mux)
	return nil
}
