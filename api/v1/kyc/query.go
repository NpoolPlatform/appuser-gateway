//nolint:dupl
package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"

	mkyc "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"
	mgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/kyc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetKyc(ctx context.Context, in *kyc.GetKycRequest) (resp *kyc.GetKycResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetKyc", "AppID", in.GetAppID())
		return &kyc.GetKycResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("GetKyc", "UserID", in.GetUserID())
		return &kyc.GetKycResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKycs")

	infos, _, err := mkyc.GetKycs(ctx, &mwpb.Conds{
		Conds: &mgrpb.Conds{
			AppID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: in.GetAppID(),
			},
			UserID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: in.GetUserID(),
			},
		},
	}, 0, 1)
	if err != nil {
		logger.Sugar().Errorw("GetKyc", "error", err)
		return &kyc.GetKycResponse{}, status.Error(codes.Internal, "fail get kyc")
	}

	if len(infos) == 0 {
		logger.Sugar().Errorw("GetKyc", "error", "not found")
		return &kyc.GetKycResponse{}, status.Error(codes.NotFound, "not found")
	}

	return &kyc.GetKycResponse{
		Info: infos[0],
	}, nil
}

func (s *Server) GetKycs(ctx context.Context, in *kyc.GetKycsRequest) (resp *kyc.GetKycsResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("GetKyc", "AppID", in.GetAppID())
		return &kyc.GetKycsResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKycs")

	infos, total, err := mkyc.GetKycs(ctx, &mwpb.Conds{
		Conds: &mgrpb.Conds{
			AppID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: in.GetAppID(),
			},
		},
	}, in.GetOffset(), in.GetLimit())

	if err != nil {
		logger.Sugar().Errorw("GetKyc", "error", err)
		return &kyc.GetKycsResponse{}, status.Error(codes.Internal, "fail get kyc")
	}

	return &kyc.GetKycsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppKycs(ctx context.Context, in *kyc.GetAppKycsRequest) (resp *kyc.GetAppKycsResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateKyc")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetKyc", "TargetAppID", in.GetTargetAppID())
		return &kyc.GetAppKycsResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	span = commontracer.TraceInvoker(span, "kyc", "middleware", "GetKycs")

	infos, total, err := mkyc.GetKycs(ctx, &mwpb.Conds{
		Conds: &mgrpb.Conds{
			AppID: &npool.StringVal{
				Op:    cruder.EQ,
				Value: in.GetTargetAppID(),
			},
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetKyc", "error", err)
		return &kyc.GetAppKycsResponse{}, status.Error(codes.Internal, "fail get kyc")
	}

	return &kyc.GetAppKycsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
