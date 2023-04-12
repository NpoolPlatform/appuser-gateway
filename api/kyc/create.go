package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	kyc1 "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/kyc"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"
	mgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateKyc(ctx context.Context, in *npool.CreateKycRequest) (resp *npool.CreateKycResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	tracer.Trace(span, &mgrpb.KycReq{
		AppID:        &in.AppID,
		UserID:       &in.UserID,
		DocumentType: &in.DocumentType,
		IDNumber:     in.IDNumber,
		FrontImg:     &in.FrontImg,
		BackImg:      in.BackImg,
		SelfieImg:    &in.SelfieImg,
		EntityType:   &in.EntityType,
	})

	err = validateKycCreate(ctx, in)
	if err != nil {
		return &npool.CreateKycResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "kyc", "kyc", "CreateKyc")

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
