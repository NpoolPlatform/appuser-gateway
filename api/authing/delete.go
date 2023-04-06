package authing

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/authing/auth"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing"
	mwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/authing"
)

func (s *Server) DeleteAppAuth(ctx context.Context, in *npool.DeleteAppAuthRequest) (resp *npool.DeleteAppAuthResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteAppAuth")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	span.SetAttributes(attribute.String("ID", in.GetID()))

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("DeleteAppAuth", "ID", in.GetID(), "error", err)
		return &npool.DeleteAppAuthResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	dInfo, err := mgrcli.DeleteAuth(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("DeleteAppAuth", "error", err)
		return &npool.DeleteAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteAppAuthResponse{
		Info: &mwpb.Auth{
			AppID:     dInfo.AppID,
			RoleID:    dInfo.RoleID,
			UserID:    dInfo.UserID,
			Resource:  dInfo.Resource,
			Method:    dInfo.Method,
			CreatedAt: dInfo.CreatedAt,
		},
	}, nil
}
