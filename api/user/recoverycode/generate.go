package recoverycode

import (
	"context"

	recoverycode1 "github.com/NpoolPlatform/appuser-gateway/pkg/user/recoverycode"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user/recoverycode"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GenerateRecoveryCodes(ctx context.Context, in *npool.GenerateRecoveryCodesRequest) (*npool.GenerateRecoveryCodesResponse, error) {
	handler, err := recoverycode1.NewHandler(
		ctx,
		recoverycode1.WithAppID(&in.AppID, true),
		recoverycode1.WithUserID(&in.UserID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GenerateRecoveryCodes",
			"In", in,
			"Error", err,
		)
		return &npool.GenerateRecoveryCodesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, err := handler.GenerateRecoveryCodes(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GenerateRecoveryCodes",
			"In", in,
			"Error", err,
		)
		return &npool.GenerateRecoveryCodesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GenerateRecoveryCodesResponse{
		Infos: infos,
	}, nil
}
