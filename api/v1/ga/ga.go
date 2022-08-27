package ga

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	ga1 "github.com/NpoolPlatform/appuser-gateway/pkg/ga"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/ga"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/google/uuid"
)

func (s *Server) SetupGoogleAuth(
	ctx context.Context, in *npool.SetupGoogleAuthRequest,
) (
	*npool.SetupGoogleAuthResponse, error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "SetupGoogleAuth")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("SetupGoogleAuth", "AppID", in.GetAppID)
		return nil, fmt.Errorf("AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("SetupGoogleAuth", "UserID", in.GetUserID)
		return nil, fmt.Errorf("UserID is invalid")
	}

	span = commontracer.TraceInvoker(span, "ga", "ga", "SetupGoogleAuth")

	info, err := ga1.SetupGoogleAuth(ctx,
		in.GetAppID(),
		in.GetUserID(),
	)
	if err != nil {
		logger.Sugar().Errorw("SetupGoogleAuth", "error", err)
		return &npool.SetupGoogleAuthResponse{}, status.Error(codes.Internal, "fail create ga")
	}

	return &npool.SetupGoogleAuthResponse{
		Info: info,
	}, nil
}

func (s *Server) VerifyGoogleAuth(
	ctx context.Context, in *npool.VerifyGoogleAuthRequest,
) (
	*npool.VerifyGoogleAuthResponse, error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "VerifyGoogleAuth")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("SetupGoogleAuth", "AppID", in.GetAppID())
		return nil, fmt.Errorf("AppID is invalid")
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Errorw("SetupGoogleAuth", "UserID", in.GetUserID())
		return nil, fmt.Errorf("UserID is invalid")
	}
	if in.GetCode() == "" {
		logger.Sugar().Errorw("SetupGoogleAuth", "GACode", in.GetCode())
		return nil, fmt.Errorf("GACode is invalid")
	}

	span = commontracer.TraceInvoker(span, "ga", "ga", "VerifyGoogleAuth")

	info, err := ga1.VerifyGoogleAuth(ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetCode(),
	)
	if err != nil {
		logger.Sugar().Errorw("VerifyGoogleAuth", "error", err)
		return &npool.VerifyGoogleAuthResponse{}, status.Error(codes.Internal, "fail create ga")
	}

	return &npool.VerifyGoogleAuthResponse{
		Info: info,
	}, nil
}
