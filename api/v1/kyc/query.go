package kyc

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"

	kyc1 "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetKyc(ctx context.Context, in *npool.GetKycRequest) (*npool.GetKycResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetKyc", "AppID", in.GetAppID())
		return &npool.GetKycResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetKyc", "UserID", in.GetUserID())
		return &npool.GetKycResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	info, err := kyc1.GetKyc(ctx, in.GetAppID(), in.UserID)
	if err != nil {
		logger.Sugar().Errorw("GetKyc", "error", err)
		return &npool.GetKycResponse{}, status.Error(codes.Internal, "fail get kyc")
	}

	return &npool.GetKycResponse{
		Info: info,
	}, nil
}

func (s *Server) GetKycs(ctx context.Context, in *npool.GetKycsRequest) (*npool.GetKycsResponse, error) {
	return nil, nil
}

func (s *Server) GetAppKycs(ctx context.Context, in *npool.GetAppKycsRequest) (*npool.GetAppKycsResponse, error) {
	return nil, nil
}
