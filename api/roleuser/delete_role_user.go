//nolint:nolintlint,dupl
package roleuser

import (
	"context"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	grpc "github.com/NpoolPlatform/appuser-manager/pkg/client"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appusergw"
	"github.com/NpoolPlatform/message/npool/appusergw/approleuser"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteRoleUser(ctx context.Context, in *approleuser.DeleteRoleUserRequest) (*approleuser.DeleteRoleUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteRoleUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Error("ID is invalid")
		return &approleuser.DeleteRoleUserResponse{}, status.Error(npool.ErrParams, appusergw.ErrMsgIDInvalid)
	}

	span.AddEvent("call grpc DeleteAppRoleUserV2")
	resp, err := grpc.DeleteAppRoleUserV2(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("fail delete app role user: %v", err)
		return &approleuser.DeleteRoleUserResponse{}, status.Error(npool.ErrService, npool.ErrMsgServiceErr)
	}

	return &approleuser.DeleteRoleUserResponse{
		Info: resp,
	}, nil
}
