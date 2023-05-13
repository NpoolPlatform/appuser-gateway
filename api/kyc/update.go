package kyc

import (
	"context"

	kyc1 "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateKyc(ctx context.Context, in *npool.UpdateKycRequest) (resp *npool.UpdateKycResponse, err error) {
	frontImg := basetypes.KycImageType_FrontImg
	backImg := basetypes.KycImageType_BackImg
	selfieImg := basetypes.KycImageType_SelfieImg

	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithID(&in.KycID),
		kyc1.WithAppID(in.GetAppID()),
		kyc1.WithUserID(in.GetUserID()),
		kyc1.WithDocumentType(in.DocumentType),
		kyc1.WithIDNumber(in.IDNumber),
		kyc1.WithImage(&frontImg, in.FrontImg),
		kyc1.WithImage(&backImg, in.BackImg),
		kyc1.WithImage(&selfieImg, in.SelfieImg),
		kyc1.WithEntityType(in.EntityType),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateKyc",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateKycResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateKyc(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateKyc",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateKycResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateKycResponse{
		Info: info,
	}, nil
}
