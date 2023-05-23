package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUserKol(ctx context.Context, in *npool.UpdateUserKolRequest) (*npool.UpdateUserKolResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithUserID(&in.UserID),
		user1.WithTargetUserID(&in.TargetUserID),
		user1.WithCheckInvitation(true),
		user1.WithKol(&in.Kol),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateUserKol",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateUserKol(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateUserKol",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateUserKolResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateUserKolResponse{
		Info: info,
	}, nil
}
