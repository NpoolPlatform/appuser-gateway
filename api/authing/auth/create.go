//nolint:dupl
package auth

import (
	"context"

	auth1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing/auth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAppAuth(ctx context.Context, in *npool.CreateAppAuthRequest) (resp *npool.CreateAppAuthResponse, err error) {
	handler, err := auth1.NewHandler(
		ctx,
		auth1.WithAppID(&in.TargetAppID, true),
		auth1.WithRoleID(in.RoleID, false),
		auth1.WithUserID(in.TargetUserID, false),
		auth1.WithResource(&in.Resource, true),
		auth1.WithMethod(&in.Method, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppAuth",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppAuthResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateAuth(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateAppAuth",
			"In", in,
			"Error", err,
		)
		return &npool.CreateAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppAuthResponse{
		Info: info,
	}, nil
}
