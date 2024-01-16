package recoverycode

import (
	"context"

	recoverycode1 "github.com/NpoolPlatform/appuser-gateway/pkg/user/recoverycode"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user/recoverycode"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRecoveryCodes(ctx context.Context, in *npool.GetRecoveryCodesRequest) (*npool.GetRecoveryCodesResponse, error) {
	handler, err := recoverycode1.NewHandler(
		ctx,
		recoverycode1.WithAppID(&in.AppID, true),
		recoverycode1.WithUserID(&in.UserID, true),
		recoverycode1.WithOffset(in.GetOffset()),
		recoverycode1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetRecoveryCodes",
			"In", in,
			"Error", err,
		)
		return &npool.GetRecoveryCodesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetRecoveryCodes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetRecoveryCodes",
			"In", in,
			"Error", err,
		)
		return &npool.GetRecoveryCodesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetRecoveryCodesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
