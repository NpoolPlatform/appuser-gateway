//nolint:nolintlint,dupl
package ban

import (
	"context"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	tracerbanuser "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banappuser"

	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"
	tracer "github.com/NpoolPlatform/appuser-manager/pkg/tracer/banapp"
	"google.golang.org/grpc/codes"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	banappmgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banapp"
	banappusermgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/banappuser"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/ban"
	banappcrud "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/banapp"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateBanApp(ctx context.Context, in *ban.UpdateBanAppRequest) (*ban.UpdateBanAppResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateBanApp")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorw("UpdateBanApp", "ID", in.GetInfo().GetID(), "err", err)
		return &ban.UpdateBanAppResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	if _, err := uuid.Parse(in.GetInfo().GetAppID()); err != nil {
		logger.Sugar().Errorw("UpdateBanApp", "AppID", in.GetInfo().GetAppID(), "err", err)
		return &ban.UpdateBanAppResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	if in.GetInfo().GetMessage() == "" {
		logger.Sugar().Errorw("UpdateBanApp", "Message", in.GetInfo().GetMessage(), "err", "Message is empty")
		return &ban.UpdateBanAppResponse{}, status.Error(codes.InvalidArgument, "Message is empty")
	}

	span = commontracer.TraceInvoker(span, "banapp", "manager", "UpdateBanApp")

	_, err = banappmgrcli.UpdateBanApp(ctx, &banappcrud.BanAppReq{
		ID:      in.GetInfo().ID,
		Message: in.GetInfo().Message,
	})
	if err != nil {
		logger.Sugar().Errorw("UpdateBanApp", "err", err)
		return &ban.UpdateBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := appmwcli.GetApp(ctx, in.GetInfo().GetAppID())
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &ban.UpdateBanAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &ban.UpdateBanAppResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateBanAppUser(ctx context.Context,
	in *ban.UpdateBanUserRequest) (*ban.UpdateBanUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateBanAppUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracerbanuser.Trace(span, in.GetInfo())

	err = validateBanUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("UpdateBanAppUser", "err", err)
		return nil, err
	}

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorw("UpdateBanAppUser", "ID", in.GetInfo().GetID(), "err", err)
		return &ban.UpdateBanUserResponse{}, status.Error(codes.InvalidArgument, "ID is invalid")
	}

	span = commontracer.TraceInvoker(span, "banappuser", "manager", "UpdateBanAppUser")

	_, err = banappusermgrcli.UpdateBanAppUser(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorw("UpdateBanAppUser", "err", err)
		return &ban.UpdateBanUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := usermwcli.GetUser(ctx, in.GetInfo().GetAppID(), in.GetInfo().GetUserID())
	if err != nil {
		logger.Sugar().Errorw("CreateBanApp", "err", err)
		return &ban.UpdateBanUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &ban.UpdateBanUserResponse{
		Info: info,
	}, nil
}
