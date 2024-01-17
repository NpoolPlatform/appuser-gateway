package app

import (
	"context"

	app1 "github.com/NpoolPlatform/appuser-gateway/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateApp(ctx context.Context, in *npool.UpdateAppRequest) (*npool.UpdateAppResponse, error) {
	handler, err := app1.NewHandler(
		ctx,
		app1.WithID(&in.ID, true),
		app1.WithEntID(&in.EntID, true),
		app1.WithNewEntID(in.NewEntID, false),
		app1.WithName(in.Name, false),
		app1.WithLogo(in.Logo, false),
		app1.WithDescription(in.Description, false),
		app1.WithSignupMethods(in.GetSignupMethods(), false),
		app1.WithExtSigninMethods(in.GetExtSigninMethods(), false),
		app1.WithRecaptchaMethod(in.RecaptchaMethod, false),
		app1.WithKycEnable(in.KycEnable, false),
		app1.WithSigninVerifyEnable(in.SigninVerifyEnable, false),
		app1.WithInvitationCodeMust(in.InvitationCodeMust, false),
		app1.WithCreateInvitationCodeWhen(in.CreateInvitationCodeWhen, false),
		app1.WithMaxTypedCouponsPerOrder(in.MaxTypedCouponsPerOrder, false),
		app1.WithMaintaining(in.Maintaining, false),
		app1.WithCouponWithdrawEnable(in.CouponWithdrawEnable, false),
		app1.WithCommitButtonTargets(in.GetCommitButtonTargets(), false),
		app1.WithResetUserMethod(in.ResetUserMethod, false),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateApp",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateApp(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateApp",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateAppResponse{
		Info: info,
	}, nil
}
