//nolint:dupl
package authing

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing"

	authing1 "github.com/NpoolPlatform/appuser-gateway/pkg/authing"
	mw "github.com/NpoolPlatform/appuser-middleware/pkg/client/authing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) Authenticate(ctx context.Context, in *npool.AuthenticateRequest) (*npool.AuthenticateResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("Authenticate", "AppID", in.GetAppID())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if in.UserID != nil && in.GetUserID() != "" {
		if _, err := uuid.Parse(in.GetUserID()); err != nil {
			logger.Sugar().Errorw("Authenticate", "UserID", in.GetUserID())
			return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
		}
	}
	if in.Token != nil && in.GetToken() == "" {
		logger.Sugar().Errorw("Authenticate", "Token", in.GetToken())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "Token is invalid")
	}
	if in.GetResource() == "" {
		logger.Sugar().Errorw("Authenticate", "Resource", in.GetResource())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "Resource is invalid")
	}
	if in.GetMethod() == "" {
		logger.Sugar().Errorw("Authenticate", "Method", in.GetMethod())
		return &npool.AuthenticateResponse{}, status.Error(codes.InvalidArgument, "Method is invalid")
	}

	allowed, err := authing1.Authenticate(ctx, in.GetAppID(), in.UserID, in.Token, in.GetResource(), in.GetMethod())
	if err != nil {
		logger.Sugar().Errorw("Authenticate", "error", err)
		return &npool.AuthenticateResponse{}, status.Error(codes.Internal, "fail authenticate")
	}

	return &npool.AuthenticateResponse{
		Info: allowed,
	}, nil
}

func (s *Server) GetAppAuths(ctx context.Context, in *npool.GetAppAuthsRequest) (resp *npool.GetAppAuthsResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppAuths")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppAuths", "TargetAppID", in.GetTargetAppID())
		return &npool.GetAppAuthsResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	infos, total, err := mw.GetAuths(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppAuths", "error", err)
		return &npool.GetAppAuthsResponse{}, status.Error(codes.Internal, "fail get app auths")
	}

	return &npool.GetAppAuthsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppHistories(ctx context.Context, in *npool.GetAppHistoriesRequest) (resp *npool.GetAppHistoriesResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppHistories")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorw("GetAppHistories", "TargetAppID", in.GetTargetAppID())
		return &npool.GetAppHistoriesResponse{}, status.Error(codes.InvalidArgument, "TargetAppID is invalid")
	}

	infos, total, err := mw.GetHistories(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppHistories", "error", err)
		return &npool.GetAppHistoriesResponse{}, status.Error(codes.Internal, "fail get app auth histories")
	}

	return &npool.GetAppHistoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
