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

func (s *Server) AuthorizeGenesis(ctx context.Context, in *npool.AuthorizeGenesisRequest) (*npool.AuthorizeGenesisResponse, error) {
	handler, err := admin1.NewHandler(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AuthorizeGenesis",
			"In", in,
			"Error", err,
		)
		return &npool.AuthorizeGenesisResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	infos, total, err := handler.AuthorizeGenesis(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"AuthorizeGenesis",
			"In", in,
			"Error", err,
		)
		return &npool.AuthorizeGenesisResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.AuthorizeGenesisResponse{
		Infos: infos,
		Total: total,
	}, nil
}
