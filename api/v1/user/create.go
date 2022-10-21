package user

import (
	"context"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/appuser-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/appuser-gateway/pkg/tracer"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	npool "github.com/NpoolPlatform/message/npool/appuser/gw/v1/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	appuserextracli "github.com/NpoolPlatform/appuser-manager/pkg/client/appuserextra"
	appuserextrapb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuserextra"

	npoolpb "github.com/NpoolPlatform/message/npool"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func (s *Server) CreateUser(ctx context.Context, in *npool.CreateUserRequest) (*npool.CreateUserResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Infow("CreateUser", "AppID", in.GetAppID())
		return &npool.CreateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetEmailAddress() == "" && in.GetPhoneNO() == "" {
		logger.Sugar().Infow("CreateUser", "AppID", in.GetAppID())
		return &npool.CreateUserResponse{}, status.Error(codes.InvalidArgument, "invalid account")
	}
	if in.GetPasswordHash() == "" {
		logger.Sugar().Infow("CreateUser", "PasswordHash", in.GetPasswordHash())
		return &npool.CreateUserResponse{}, status.Error(codes.InvalidArgument, "invalid password")
	}
	if in.ImportedFromAppID != nil {
		if _, err := uuid.Parse(in.GetImportedFromAppID()); err != nil {
			logger.Sugar().Infow("CreateUser", "ImportedFromAppID", in.GetImportedFromAppID())
			return &npool.CreateUserResponse{}, status.Error(codes.InvalidArgument, "invalid password")
		}
	}

	if in.IDNumber != nil {
		if in.GetIDNumber() == "" {
			logger.Sugar().Infow("UpdateUser", "IDNumber", in.GetIDNumber())
			return &npool.CreateUserResponse{}, status.Error(codes.InvalidArgument, "IDNumber is invalid")
		}
		exist, err := appuserextracli.ExistAppUserExtraConds(ctx, &appuserextrapb.Conds{
			IDNumber: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: in.GetIDNumber(),
			},
		})
		if err != nil {
			logger.Sugar().Infow("CreateUser", "exist", exist, "err", err)
			return &npool.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
		}
		if exist {
			logger.Sugar().Infow("CreateUser", "IDNumber", in.GetIDNumber())
			return &npool.CreateUserResponse{}, status.Error(codes.InvalidArgument, "IDNumber is already exists")
		}
	}

	span = commontracer.TraceInvoker(span, "user", "middleware", "CreateUser")

	userID := uuid.NewString()
	info, err := usermwcli.CreateUser(ctx, &usermwpb.UserReq{
		ID:                &userID,
		AppID:             &in.AppID,
		EmailAddress:      in.EmailAddress,
		PhoneNO:           in.PhoneNO,
		ImportedFromAppID: in.ImportedFromAppID,
		Username:          in.Username,
		AddressFields:     in.AddressFields,
		Gender:            in.Gender,
		PostalCode:        in.PostalCode,
		Age:               in.Age,
		Birthday:          in.Birthday,
		Avatar:            in.Avatar,
		Organization:      in.Organization,
		FirstName:         in.FirstName,
		LastName:          in.LastName,
		IDNumber:          in.IDNumber,
		PasswordHash:      in.PasswordHash,
	})
	if err != nil {
		logger.Sugar().Errorw("CreateUser", "err", err)
		return &npool.CreateUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateUserResponse{
		Info: info,
	}, nil
}
