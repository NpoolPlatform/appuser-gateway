package user

import (
	"context"

	user1 "github.com/NpoolPlatform/appuser-gateway/pkg/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, in *npool.UpdateUserRequest) (*npool.UpdateUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithUserID(&in.UserID),
		user1.WithNewAccount(in.NewAccount, in.NewAccountType),
		user1.WithPasswordHash(in.PasswordHash),
		user1.WithOldPasswordHash(in.OldPasswordHash),
		user1.WithVerificationCode(in.VerificationCode),
		user1.WithNewVerificationCode(in.NewVerificationCode),
		user1.WithIDNumber(in.IDNumber),
		// user1.WithUsername(&in.Username),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateUser",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateUser(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateUser",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateUserResponse{
		Info: info,
	}, nil
}

func (s *Server) ResetUser(ctx context.Context, in *npool.ResetUserRequest) (*npool.ResetUserResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithUserID(in.UserID),
		user1.WithAccount(in.GetAccount(), in.GetAccountType()),
		user1.WithVerificationCode(&in.VerificationCode),
		user1.WithPasswordHash(in.PasswordHash),
		user1.WithRecoveryCode(in.RecoveryCode),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"ResetUser",
			"In", in,
			"Error", err,
		)
		return &npool.ResetUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := handler.ResetUser(ctx); err != nil {
		logger.Sugar().Errorw(
			"ResetUser",
			"In", in,
			"Error", err,
		)
		return &npool.ResetUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ResetUserResponse{}, nil
}

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

func (s *Server) UpdateAppUserKol(ctx context.Context, in *npool.UpdateAppUserKolRequest) (*npool.UpdateAppUserKolResponse, error) {
	handler, err := user1.NewHandler(
		ctx,
		user1.WithAppID(in.GetAppID()),
		user1.WithTargetUserID(&in.TargetUserID),
		user1.WithKol(&in.Kol),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUserKol",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserKolResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateUserKol(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateAppUserKol",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppUserKolResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppUserKolResponse{
		Info: info,
	}, nil
}
