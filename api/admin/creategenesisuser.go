//nolint:nolintlint,dupl
package admin

import (
	"context"

	admin1 "github.com/NpoolPlatform/appuser-gateway/pkg/admin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/admin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateGenesisUser(
	ctx context.Context,
	in *admin.CreateGenesisUserRequest,
) (
	*admin.CreateGenesisUserResponse,
	error,
) {
	handler, err := admin1.NewHandler(
		ctx,
		admin1.WithAppID(in.GetTargetAppID()),
		admin1.WithEmailAddress(&in.EmailAddress),
		admin1.WithPasswordHash(&in.PasswordHash),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateGenesisUser",
			"In", in,
			"Error", err,
		)
		return &admin.CreateGenesisUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateGenesisUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateGenesisUser",
			"In", in,
			"Error", err,
		)
		return &admin.CreateGenesisUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &admin.CreateGenesisUserResponse{
		Info: info,
	}, nil
}
