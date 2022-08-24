package authing

import (
	"context"

	mauthing "github.com/NpoolPlatform/appuser-gateway/pkg/authing"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/authing/auth"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/authing"
	mgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/authing/auth"
)

func (s *Server) CreateAppAuth(ctx context.Context, in *npool.CreateAppAuthRequest) (resp *npool.CreateAppAuthResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppAuth")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in)
	if err != nil {
		logger.Sugar().Errorw("CreateAppAuth", "err", err)
		return nil, err
	}

	span = tracer.Trace(span, &mgrpb.AuthReq{
		AppID:    &in.TargetAppID,
		RoleID:   in.RoleID,
		UserID:   in.TargetUserID,
		Resource: &in.Resource,
		Method:   &in.Method,
	})
	span = commontracer.TraceInvoker(span, "auth", "auth", "CreateAuth")

	info, err := mauthing.CreateAuth(ctx, in.TargetAppID, in.TargetUserID, in.RoleID, in.Resource, in.Method)
	if err != nil {
		logger.Sugar().Errorw("CreateAppAuth", "err", err)
		return &npool.CreateAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAppAuthResponse{
		Info: info,
	}, nil
}
