package subscriber

import (
	"context"
	"net/mail"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	mgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/subscriber"
	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	mgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/subscriber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) DeleteSubscriber(ctx context.Context, in *npool.DeleteSubscriberRequest) (*npool.DeleteSubscriberResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := mail.ParseAddress(in.GetEmailAddress()); err != nil {
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	info, err := mgrcli.GetSubscriberOnly(ctx, &mgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
		EmailAddress: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetEmailAddress(),
		},
	})
	if err != nil {
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if info == nil {
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.InvalidArgument, "Subscriber is invalid")
	}

	info1, err := mwcli.DeleteSubscriber(ctx, info.ID)
	if err != nil {
		return &npool.DeleteSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteSubscriberResponse{
		Info: info1,
	}, nil
}
