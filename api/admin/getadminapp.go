//nolint:nolintlint,dupl
package admin

import (
	"context"

	admin1 "github.com/NpoolPlatform/appuser-gateway/pkg/admin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAdminApps(ctx context.Context, in *npool.GetAdminAppsRequest) (*npool.GetAdminAppsResponse, error) {
	handler, err := admin1.NewHandler(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAdminApps",
			"In", in,
			"Error", err,
		)
		return &npool.GetAdminAppsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, err := handler.GetGenesisApps(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAdminApps",
			"In", in,
			"Error", err,
		)
		return &npool.GetAdminAppsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAdminAppsResponse{
		Infos: infos,
	}, nil
}
