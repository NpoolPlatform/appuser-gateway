package main

import (
	"fmt"

	"github.com/NpoolPlatform/appuser-gateway/pkg/migrator"

	"github.com/NpoolPlatform/appuser-gateway/api"
	"github.com/NpoolPlatform/go-service-framework/pkg/oss"
	ossconst "github.com/NpoolPlatform/go-service-framework/pkg/oss/const"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	apimgrcli "github.com/NpoolPlatform/api-manager/pkg/client"

	cli "github.com/urfave/cli/v2"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const BukectKey = "kyc_bucket"

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"s"},
	Usage:   "Run the daemon",
	Action: func(c *cli.Context) error {
		if err := migrator.Migrate(c.Context); err != nil {
			return err
		}

		if err := oss.Init(ossconst.SecretStoreKey, BukectKey); err != nil {
			return fmt.Errorf("fail to init s3: %v", err)
		}

		go func() {
			if err := grpc2.RunGRPC(rpcRegister); err != nil {
				logger.Sugar().Errorf("fail to run grpc server: %v", err)
			}
		}()
		return grpc2.RunGRPCGateWay(rpcGatewayRegister)
	},
}

func rpcRegister(server grpc.ServiceRegistrar) error {
	api.Register(server)

	apimgrcli.RegisterGRPC(server)

	return nil
}

func rpcGatewayRegister(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	err := api.RegisterGateway(mux, endpoint, opts)
	if err != nil {
		return err
	}

	_ = apimgrcli.Register(mux)
	return nil
}
