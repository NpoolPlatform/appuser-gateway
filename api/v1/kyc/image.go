package kyc

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"

	kyc1 "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UploadKycImage(ctx context.Context, in *npool.UploadKycImageRequest) (*npool.UploadKycImageResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("UploadKycImage", "AppID", in.GetAppID())
		return &npool.UploadKycImageResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("UploadKycImage", "UserID", in.GetUserID())
		return &npool.UploadKycImageResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}
	if in.GetImageBase64() == "" {
		logger.Sugar().Errorw("UploadKycImage", "ImageBase64", in.GetImageBase64())
		return &npool.UploadKycImageResponse{}, status.Error(codes.InvalidArgument, "ImageBase64 is invalid")
	}

	switch in.GetImageType() {
	case kycmgrpb.KycImageType_FrontImg:
	case kycmgrpb.KycImageType_BackImg:
	case kycmgrpb.KycImageType_SelfieImg:
	default:
		logger.Sugar().Errorw("UploadKycImage", "ImageType", in.GetImageType())
		return &npool.UploadKycImageResponse{}, status.Error(codes.InvalidArgument, "ImageType is invalid")
	}

	err := kyc1.UploadKycImage(ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetImageType(),
		in.GetImageBase64(),
	)
	if err != nil {
		logger.Sugar().Errorw("UploadKycImage", "error", err)
		return &npool.UploadKycImageResponse{}, status.Error(codes.Internal, "fail create kyc")
	}

	return &npool.UploadKycImageResponse{}, nil
}

func (s *Server) GetKycImage(ctx context.Context, in *npool.GetKycImageRequest) (*npool.GetKycImageResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetKycImage", "AppID", in.GetAppID())
		return &npool.GetKycImageResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetKycImage", "UserID", in.GetUserID())
		return &npool.GetKycImageResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	switch in.GetImageType() {
	case kycmgrpb.KycImageType_FrontImg:
	case kycmgrpb.KycImageType_BackImg:
	case kycmgrpb.KycImageType_SelfieImg:
	default:
		logger.Sugar().Errorw("GetKycImage", "ImageType", in.GetImageType())
		return &npool.GetKycImageResponse{}, status.Error(codes.InvalidArgument, "ImageType is invalid")
	}

	imgBase64, err := kyc1.GetKycImage(ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetImageType(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetKycImage", "error", err)
		return &npool.GetKycImageResponse{}, status.Error(codes.Internal, "fail create kyc")
	}

	return &npool.GetKycImageResponse{
		Info: imgBase64,
	}, nil
}
