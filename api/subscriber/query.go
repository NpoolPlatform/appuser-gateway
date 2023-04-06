package subscriber

import (
	"context"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	mgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/subscriber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) GetSubscriberes(ctx context.Context, in *npool.GetSubscriberesRequest) (*npool.GetSubscriberesResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetSubscriberesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.GetSubscriberesResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.GetSubscriberesResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	infos, total, err := mwcli.GetSubscriberes(ctx, &mgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		return &npool.GetSubscriberesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetSubscriberesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
