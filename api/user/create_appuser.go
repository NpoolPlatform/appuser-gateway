package user

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/appuser"
	appusercrud "github.com/NpoolPlatform/message/npool/appusermgrv2/appuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, in *appuser.CreateUserRequest) (*appuser.CreateUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validate(in.GetInfo())
	if err != nil {
		return &appuser.CreateUserResponse{}, err
	}

	if in.GetInfo().GetID() != "" {
		span.AddEvent("call grpc ExistAppUserV2")
		exist, err := grpc.ExistAppUserV2(ctx, in.GetInfo().GetID())
		if err != nil {
			logger.Sugar().Errorw("fail check user: %v", err)
			return &appuser.CreateUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
		}
		if exist {
			return &appuser.CreateUserResponse{}, status.Error(npool.ErrAlreadyExists, appusergw.ErrMsgUserIDAlreadyExists)
		}
	}

	importFromApp := in.GetInfo().GetImportFromApp()
	if in.GetInfo().GetImportFromApp() == "" {
		importFromApp = uuid.UUID{}.String()
	}
	phoneNo := in.GetInfo().GetPhoneNo()
	id := in.GetInfo().GetID()
	appID := in.GetInfo().GetAppID()
	emailAddress := in.GetInfo().GetEmailAddress()

	span.AddEvent("call grpc CreateAppUserV2")
	resp, err := grpc.CreateAppUserV2(ctx, &appusercrud.AppUserReq{
		PhoneNo:       &phoneNo,
		ImportFromApp: &importFromApp,
		ID:            &id,
		AppID:         &appID,
		EmailAddress:  &emailAddress,
	})
	if err != nil {
		logger.Sugar().Errorw("fail create user: %v", err)
		return &appuser.CreateUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &appuser.CreateUserResponse{
		Info: resp,
	}, nil
}
