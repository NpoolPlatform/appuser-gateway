//nolint:dupl
package kyc

import (
	"context"

	kyc1 "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UploadKycImage(ctx context.Context, in *npool.UploadKycImageRequest) (resp *npool.UploadKycImageResponse, err error) {
	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithAppID(&in.AppID, true),
		kyc1.WithUserID(&in.UserID, true),
		kyc1.WithImage(&in.ImageType, &in.ImageBase64, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.UploadKycImageResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	key, err := handler.UploadKycImage(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UploadKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.UploadKycImageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UploadKycImageResponse{
		Info: key,
	}, nil
}

func (s *Server) GetKycImage(ctx context.Context, in *npool.GetKycImageRequest) (resp *npool.GetKycImageResponse, err error) {
	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithAppID(&in.AppID, true),
		kyc1.WithUserID(&in.UserID, true),
		kyc1.WithImage(&in.ImageType, nil, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.GetKycImageResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	imgBase64, err := handler.GetKycImage(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.GetKycImageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetKycImageResponse{
		Info: imgBase64,
	}, nil
}

func (s *Server) GetUserKycImage(ctx context.Context, in *npool.GetUserKycImageRequest) (resp *npool.GetUserKycImageResponse, err error) {
	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithAppID(&in.AppID, true),
		kyc1.WithUserID(&in.TargetUserID, true),
		kyc1.WithImage(&in.ImageType, nil, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetUserKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.GetUserKycImageResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	imgBase64, err := handler.GetKycImage(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetUserKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.GetUserKycImageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetUserKycImageResponse{
		Info: imgBase64,
	}, nil
}

func (s *Server) GetAppUserKycImage(ctx context.Context, in *npool.GetAppUserKycImageRequest) (resp *npool.GetAppUserKycImageResponse, err error) {
	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithAppID(&in.TargetAppID, true),
		kyc1.WithUserID(&in.TargetUserID, true),
		kyc1.WithImage(&in.ImageType, nil, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserKycImageResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	imgBase64, err := handler.GetKycImage(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppUserKycImage",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppUserKycImageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppUserKycImageResponse{
		Info: imgBase64,
	}, nil
}
