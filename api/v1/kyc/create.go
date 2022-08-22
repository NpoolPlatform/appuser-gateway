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

func (s *Server) CreateKyc(ctx context.Context, in *npool.CreateKycRequest) (*npool.CreateKycResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("CreateKyc", "AppID", in.GetAppID())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("CreateKyc", "UserID", in.GetUserID())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}
	if in.IDNumber != nil && in.GetIDNumber() == "" {
		logger.Sugar().Errorw("CreateKyc", "IDNumber", in.GetIDNumber())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "IDNumber is invalid")
	}
	if in.GetFrontImg() == "" {
		logger.Sugar().Errorw("CreateKyc", "FrontImg", in.GetFrontImg())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "FrontImg is invalid")
	}
	if in.BackImg != nil && in.GetBackImg() == "" {
		logger.Sugar().Errorw("CreateKyc", "BackImg", in.GetBackImg())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "BackImg is invalid")
	}
	if in.GetSelfieImg() == "" {
		logger.Sugar().Errorw("CreateKyc", "SelfieImg", in.GetSelfieImg())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "SelfieImg is invalid")
	}

	switch in.GetDocumentType() {
	case kycmgrpb.KycDocumentType_IDCard:
	case kycmgrpb.KycDocumentType_DriverLicense:
	case kycmgrpb.KycDocumentType_Passport:
	default:
		logger.Sugar().Errorw("CreateKyc", "DocumentType", in.GetDocumentType())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "DocumentType is invalid")
	}

	switch in.GetEntityType() {
	case kycmgrpb.KycEntityType_Individual:
	case kycmgrpb.KycEntityType_Organization:
	default:
		logger.Sugar().Errorw("CreateKyc", "EntityType", in.GetEntityType())
		return &npool.CreateKycResponse{}, status.Error(codes.InvalidArgument, "EntityType is invalid")
	}

	info, err := kyc1.CreateKyc(ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetFrontImg(),
		in.GetSelfieImg(),
		in.IDNumber,
		in.BackImg,
		in.GetDocumentType(),
		in.GetEntityType(),
	)
	if err != nil {
		logger.Sugar().Errorw("CreateKyc", "error", err)
		return &npool.CreateKycResponse{}, status.Error(codes.Internal, "fail create kyc")
	}

	return &npool.CreateKycResponse{
		Info: info,
	}, nil
}
