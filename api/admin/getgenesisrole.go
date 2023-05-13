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

func (s *Server) GetGenesisRoles(ctx context.Context, in *npool.GetGenesisRolesRequest) (*npool.GetGenesisRolesResponse, error) {
	handler, err := admin1.NewHandler(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetGenesisRoles",
			"In", in,
			"Error", err,
		)
		return &npool.GetGenesisRolesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, err := handler.GetGenesisRoles(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetGenesisRoles",
			"In", in,
			"Error", err,
		)
		return &npool.GetGenesisRolesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetGenesisRolesResponse{
		Infos: infos,
		Total: uint32(len(infos)),
	}, nil
}
