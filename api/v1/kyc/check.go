package kyc

import (
	"context"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/gw/v1/kyc"
	kycmgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/kyc"
	"github.com/google/uuid"

	mgrcli "github.com/NpoolPlatform/appuser-manager/pkg/client/kyc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint:gocyclo
func validateKycCreate(ctx context.Context, info *kyc.CreateKycRequest) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Errorw("validateKycCreate", "AppID", info.GetAppID())
		return status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Errorw("validateKycCreate", "UserID", info.GetUserID())
		return status.Error(codes.InvalidArgument, "UserID is invalid")
	}
	if info.IDNumber != nil && info.GetIDNumber() == "" {
		logger.Sugar().Errorw("validateKycCreate", "IDNumber", info.GetIDNumber())
		return status.Error(codes.InvalidArgument, "IDNumber is invalid")
	}
	if info.GetFrontImg() == "" {
		logger.Sugar().Errorw("validateKycCreate", "FrontImg", info.GetFrontImg())
		return status.Error(codes.InvalidArgument, "FrontImg is invalid")
	}
	if info.BackImg != nil && info.GetBackImg() == "" {
		logger.Sugar().Errorw("validateKycCreate", "BackImg", info.GetBackImg())
		return status.Error(codes.InvalidArgument, "BackImg is invalid")
	}
	if info.GetSelfieImg() == "" {
		logger.Sugar().Errorw("validateKycCreate", "SelfieImg", info.GetSelfieImg())
		return status.Error(codes.InvalidArgument, "SelfieImg is invalid")
	}

	switch info.GetDocumentType() {
	case kycmgrpb.KycDocumentType_IDCard:
	case kycmgrpb.KycDocumentType_DriverLicense:
	case kycmgrpb.KycDocumentType_Passport:
	default:
		logger.Sugar().Errorw("validateKycCreate", "DocumentType", info.GetDocumentType())
		return status.Error(codes.InvalidArgument, "DocumentType is invalid")
	}

	switch info.GetEntityType() {
	case kycmgrpb.KycEntityType_Individual:
	case kycmgrpb.KycEntityType_Organization:
	default:
		logger.Sugar().Errorw("validateKycCreate", "EntityType", info.GetEntityType())
		return status.Error(codes.InvalidArgument, "EntityType is invalid")
	}

	exist, err := mgrcli.ExistKycConds(ctx, &kycmgrpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetAppID(),
		},
		UserID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: info.GetUserID(),
		},
	})
	if err != nil {
		return err
	}

	if exist {
		logger.Sugar().Errorw("validateKycCreate", "")
		return status.Error(codes.AlreadyExists, "already exists")
	}

	return checkIDNumberDuplicate(ctx, info.GetDocumentType().String(), info.GetIDNumber())
}

func validateKycUpdate(ctx context.Context, info *kyc.UpdateKycRequest) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Errorw("validateKycUpdate", "AppID", info.GetAppID())
		return status.Error(codes.InvalidArgument, "AppID is invalid")
	}
	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		logger.Sugar().Errorw("validateKycUpdate", "UserID", info.GetUserID())
		return status.Error(codes.InvalidArgument, "UserID is invalid")
	}
	if info.IDNumber != nil && info.GetIDNumber() == "" {
		logger.Sugar().Errorw("validateKycUpdate", "IDNumber", info.GetIDNumber())
		return status.Error(codes.InvalidArgument, "IDNumber is invalid")
	}
	if info.FrontImg != nil && info.GetFrontImg() == "" {
		logger.Sugar().Errorw("validateKycUpdate", "FrontImg", info.GetFrontImg())
		return status.Error(codes.InvalidArgument, "FrontImg is invalid")
	}
	if info.BackImg != nil && info.GetBackImg() == "" {
		logger.Sugar().Errorw("validateKycUpdate", "BackImg", info.GetBackImg())
		return status.Error(codes.InvalidArgument, "BackImg is invalid")
	}
	if info.SelfieImg != nil && info.GetSelfieImg() == "" {
		logger.Sugar().Errorw("validateKycUpdate", "SelfieImg", info.GetSelfieImg())
		return status.Error(codes.InvalidArgument, "SelfieImg is invalid")
	}

	if info.IDNumber != nil && info.GetIDNumber() != "" {
		kycInfo, err := mgrcli.GetKyc(ctx, info.KycID)
		if err != nil {
			logger.Sugar().Errorw("validateKycUpdate", "err", err)
			return status.Error(codes.Internal, err.Error())
		}

		if info.GetIDNumber() == kycInfo.IDNumber {
			return nil
		}

		return checkIDNumberDuplicate(ctx, kycInfo.DocumentType.String(), info.GetIDNumber())
	}
	return nil
}

func checkIDNumberDuplicate(ctx context.Context, documentType, idNumber string) error {
	existIDNumber, err := mgrcli.ExistKycConds(ctx, &kycmgrpb.Conds{
		DocumentType: &npool.StringVal{
			Op:    cruder.EQ,
			Value: documentType,
		},
		IDNumber: &npool.StringVal{
			Op:    cruder.EQ,
			Value: idNumber,
		},
	})
	if err != nil {
		logger.Sugar().Errorw("validateKycUpdate", "err", err)
		return status.Error(codes.Internal, err.Error())
	}

	if existIDNumber {
		logger.Sugar().Errorw("validateKycUpdate", "idNumber", idNumber)
		return status.Error(codes.AlreadyExists, "IDNumber already exists")
	}
	return nil
}
