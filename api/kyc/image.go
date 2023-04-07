//nolint:dupl
package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"

	kyc1 "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) UploadKycImage(ctx context.Context, in *npool.UploadKycImageRequest) (resp *npool.UploadKycImageResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UploadKycImage")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

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

	commontracer.TraceInvoker(span, "kyc", "kyc", "UploadKycImage")

	key, err := kyc1.UploadKycImage(ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetImageType(),
		in.GetImageBase64(),
	)
	if err != nil {
		logger.Sugar().Errorw("UploadKycImage", "error", err)
		return &npool.UploadKycImageResponse{}, status.Error(codes.Internal, "fail create kyc")
	}

	return &npool.UploadKycImageResponse{
		Info: key,
	}, nil
}

func (s *Server) GetKycImage(ctx context.Context, in *npool.GetKycImageRequest) (resp *npool.GetKycImageResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetKycImage")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

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

	span = commontracer.TraceInvoker(span, "kyc", "kyc", "GetKycImage")

	imgBase64, err := kyc1.GetKycImage(ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetImageType(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetKycImage", "error", err)
		return &npool.GetKycImageResponse{}, status.Error(codes.Internal, "fail get kyc image")
	}

	return &npool.GetKycImageResponse{
		Info: imgBase64,
	}, nil
}

func (s *Server) GetUserKycImage(ctx context.Context, in *npool.GetUserKycImageRequest) (resp *npool.GetUserKycImageResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetUserKycImage")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetUserKycImage", "AppID", in.GetAppID())
		return &npool.GetUserKycImageResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		logger.Sugar().Errorw("GetUserKycImage", "TargetUserID", in.GetTargetUserID())
		return &npool.GetUserKycImageResponse{}, status.Error(codes.InvalidArgument, "TargetUserID is invalid")
	}

	switch in.GetImageType() {
	case kycmgrpb.KycImageType_FrontImg:
	case kycmgrpb.KycImageType_BackImg:
	case kycmgrpb.KycImageType_SelfieImg:
	default:
		logger.Sugar().Errorw("GetUserKycImage", "ImageType", in.GetImageType())
		return &npool.GetUserKycImageResponse{}, status.Error(codes.InvalidArgument, "ImageType is invalid")
	}

	span = commontracer.TraceInvoker(span, "kyc", "kyc", "GetUserKycImage")

	imgBase64, err := kyc1.GetKycImage(ctx,
		in.GetAppID(),
		in.GetTargetUserID(),
		in.GetImageType(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetUserKycImage", "error", err)
		return &npool.GetUserKycImageResponse{}, status.Error(codes.Internal, "fail get kyc image")
	}

	return &npool.GetUserKycImageResponse{
		Info: imgBase64,
	}, nil
}

func (s *Server) GetAppUserKycImage(
	ctx context.Context, in *npool.GetAppUserKycImageRequest,
) (
	resp *npool.GetAppUserKycImageResponse, err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppUserKycImage")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppUserKycImage", "TargetAppID", in.GetTargetAppID())
		return &npool.GetAppUserKycImageResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}
	if _, err := uuid.Parse(in.GetTargetUserID()); err != nil {
		logger.Sugar().Errorw("GetAppUserKycImage", "TargetUserID", in.GetTargetUserID())
		return &npool.GetAppUserKycImageResponse{}, status.Error(codes.InvalidArgument, "TargetUserID is invalid")
	}

	switch in.GetImageType() {
	case kycmgrpb.KycImageType_FrontImg:
	case kycmgrpb.KycImageType_BackImg:
	case kycmgrpb.KycImageType_SelfieImg:
	default:
		logger.Sugar().Errorw("GetAppUserKycImage", "ImageType", in.GetImageType())
		return &npool.GetAppUserKycImageResponse{}, status.Error(codes.InvalidArgument, "ImageType is invalid")
	}

	span = commontracer.TraceInvoker(span, "kyc", "kyc", "GetAppUserKycImage")

	imgBase64, err := kyc1.GetKycImage(ctx,
		in.GetTargetAppID(),
		in.GetTargetUserID(),
		in.GetImageType(),
	)
	if err != nil {
		logger.Sugar().Errorw("GetAppUserKycImage", "error", err)
		return &npool.GetAppUserKycImageResponse{}, status.Error(codes.Internal, "fail get kyc image")
	}

	return &npool.GetAppUserKycImageResponse{
		Info: imgBase64,
	}, nil
}