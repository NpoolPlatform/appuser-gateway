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

func (s *Server) UpdateKyc(ctx context.Context, in *npool.UpdateKycRequest) (resp *npool.UpdateKycResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	tracer.Trace(span, &mgrpb.KycReq{
		AppID:     &in.AppID,
		UserID:    &in.UserID,
		IDNumber:  in.IDNumber,
		FrontImg:  in.FrontImg,
		BackImg:   in.BackImg,
		SelfieImg: in.SelfieImg,
	})

	err = validateKycUpdate(ctx, in)
	if err != nil {
		return &npool.UpdateKycResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "kyc", "kyc", "UpdateKyc")

	info, err := kyc1.UpdateKyc(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("UpdateKyc", "error", err)
		return &npool.UpdateKycResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateKycResponse{
		Info: info,
	}, nil
}
