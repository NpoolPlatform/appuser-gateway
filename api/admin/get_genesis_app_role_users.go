package admin

import (
	"context"
	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	mw "github.com/NpoolPlatform/appuser-gateway/pkg/middleware/admin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/admin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAppGenesisAppRoleUsers(ctx context.Context, in *admin.GetAppGenesisAppRoleUsersRequest) (*admin.GetAppGenesisAppRoleUsersResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateExtra")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Error("TargetAppID is invalid")
		return &admin.GetAppGenesisAppRoleUsersResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgAppIDInvalid)
	}
	if _, err := uuid.Parse(in.GetUserID()); err != nil {
		logger.Sugar().Error("UserID is invalid")
		return &admin.GetAppGenesisAppRoleUsersResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgUserIDInvalid)
	}

	span.AddEvent("call middleware GetAppGenesisAppRoleUsers")
	resp, err := mw.GetAppGenesisAppRoleUsers(ctx, in.GetTargetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorw("fail get genesis app role user: %v", err)
		return &admin.GetAppGenesisAppRoleUsersResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &admin.GetAppGenesisAppRoleUsersResponse{
		Infos: resp,
	}, nil
}
