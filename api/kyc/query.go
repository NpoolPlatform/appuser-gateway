//nolint:dupl
package kyc

import (
	"context"

	kyc1 "github.com/NpoolPlatform/appuser-gateway/pkg/kyc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetKyc(ctx context.Context, in *npool.GetKycRequest) (resp *npool.GetKycResponse, err error) {
	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithAppID(in.GetAppID()),
		kyc1.WithUserID(in.GetUserID()),
		kyc1.WithOffset(0),
		kyc1.WithLimit(1),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetKyc",
			"In", in,
			"Error", err,
		)
		return &npool.GetKycResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, _, err := handler.GetKycs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetKyc",
			"In", in,
			"Error", err,
		)
		return &npool.GetKycResponse{}, status.Error(codes.Internal, err.Error())
	}
	if len(infos) == 0 {
		return &npool.GetKycResponse{}, status.Error(codes.Internal, "not found")
	}

	return &npool.GetKycResponse{
		Info: infos[0],
	}, nil
}

func (s *Server) GetKycs(ctx context.Context, in *npool.GetKycsRequest) (resp *npool.GetKycsResponse, err error) {
	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithAppID(in.GetAppID()),
		kyc1.WithOffset(in.GetOffset()),
		kyc1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetKycs",
			"In", in,
			"Error", err,
		)
		return &npool.GetKycsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetKycs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetKycs",
			"In", in,
			"Error", err,
		)
		return &npool.GetKycsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetKycsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppKycs(ctx context.Context, in *npool.GetAppKycsRequest) (resp *npool.GetAppKycsResponse, err error) {
	handler, err := kyc1.NewHandler(
		ctx,
		kyc1.WithAppID(in.GetTargetAppID()),
		kyc1.WithOffset(in.GetOffset()),
		kyc1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppKycs",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppKycsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetKycs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppKycs",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppKycsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppKycsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
