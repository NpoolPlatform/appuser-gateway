package subscriber

import (
	"context"
	"net/mail"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/subscriber"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	mwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/subscriber"
	mgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/subscriber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateSubscriber(ctx context.Context, in *npool.CreateSubscriberRequest) (*npool.CreateSubscriberResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.CreateSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := mail.ParseAddress(in.GetEmailAddress()); err != nil {
		return &npool.CreateSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		return &npool.CreateSubscriberResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		return &npool.CreateSubscriberResponse{}, status.Error(codes.InvalidArgument, "AppID is invalid")
	}

	info, err := mwcli.CreateSubscriber(ctx, &mgrpb.SubscriberReq{
		AppID:        &in.AppID,
		EmailAddress: &in.EmailAddress,
	})
	if err != nil {
		return &npool.CreateSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateSubscriberResponse{
		Info: info,
	}, nil
}
