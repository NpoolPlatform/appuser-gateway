package appoauththirdparty

import (
	"context"

	oauth1 "github.com/NpoolPlatform/appuser-gateway/pkg/oauth/appoauththirdparty"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/oauth/appoauththirdparty"
	oauthmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/oauth/appoauththirdparty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOAuthThirdParties(ctx context.Context, in *npool.GetOAuthThirdPartiesRequest) (resp *npool.GetOAuthThirdPartiesResponse, err error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithAppID(in.GetAppID()),
		oauth1.WithOffset(in.GetOffset()),
		oauth1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthThirdParties",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthThirdPartiesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetOAuthThirdParties(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetOAuthThirdParties",
			"In", in,
			"Error", err,
		)
		return &npool.GetOAuthThirdPartiesResponse{}, status.Error(codes.Internal, err.Error())
	}

	lists := []*oauthmwpb.OAuthThirdParty{}
	for _, info := range infos {
		thirdPartyInfo := &oauthmwpb.OAuthThirdParty{
			ClientName:    info.ClientName,
			ClientTag:     info.ClientTag,
			ClientLogoURL: info.ClientLogoURL,
		}
		lists = append(lists, thirdPartyInfo)
	}

	return &npool.GetOAuthThirdPartiesResponse{
		Infos: lists,
		Total: total,
	}, nil
}

func (s *Server) GetAppOAuthThirdParties(ctx context.Context, in *npool.GetAppOAuthThirdPartiesRequest) (resp *npool.GetAppOAuthThirdPartiesResponse, err error) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithAppID(in.GetTargetAppID()),
		oauth1.WithOffset(in.GetOffset()),
		oauth1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppOAuthThirdParties",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppOAuthThirdPartiesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetOAuthThirdParties(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppOAuthThirdParties",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppOAuthThirdPartiesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppOAuthThirdPartiesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
